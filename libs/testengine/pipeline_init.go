//go:build generatetests

package testengine

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/go-yaml/yaml"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/zale144/ube/actions"
	"github.com/zale144/ube/model"
	pl "github.com/zale144/ube/pipeline"
)

/*

* get seed data: can be from file, or generated randomly
* instantiate Test struct
* setup dependencies and pass Test struct to them
* create new pipeline from dependencies
* use data to run pipeline
* each mock dependency should add its calls with inputs, outputs and methods to Test
* when done, use the filled Test struct to generate a yaml file

 */

type Case struct {
	Entity            model.Entity
	Feed              interface{} // optional
	EventMapping      map[string]string
	EventNames        []string
	NumMessages       int
	RecordsPerMessage int
	WithNegative      bool
}

func EventHandlerInit(t *testing.T, testCase Case) []*Test {
	model.Now = func() time.Time { return time.Date(2021, 11, 22, 3, 4, 5, 0, time.UTC) }
	model.UUIDStr = func() string { return "5de1ea04-61c9-4cf8-bdf8-320479e62d31" }

	// success test
	cases := initCases(t, testCase, nil)

	if testCase.WithNegative {
		for _, dp := range cases[0].Dependencies {
			// no fail tests for the republisher
			if dp.Name == "Republisher" {
				continue
			}

			d := *dp

			// each call gets a negative test
			for i := range d.Calls {
				d.failCall = d.Calls[i].Method
				failCases := initCases(t, testCase, &d)
				cases = append(cases, failCases...)
			}
		}
	}

	yml, err := yaml.Marshal(cases)
	require.NoError(t, err)
	require.NoError(t, os.MkdirAll(testDir, os.ModePerm))
	require.NoError(t, os.WriteFile(fmt.Sprintf("%s/generated.yaml", testDir), yml, 0644))

	return cases
}

func initCases(t *testing.T, testCase Case, failDep *dependency) []*Test {
	var tests []*Test

	if len(testCase.EventMapping) > 0 {
		for sourceURI, eventName := range testCase.EventMapping {
			tests = append(tests, createTestCase(t, testCase, sourceURI, eventName, failDep))
		}
	} else if len(testCase.EventNames) > 0 {
		for _, eventName := range testCase.EventNames {
			tests = append(tests, createTestCase(t, testCase, "", eventName, failDep))
		}
	} else {
		tests = append(tests, createTestCase(t, testCase, "", "", failDep))
	}

	require.NotEmptyf(t, tests, "not tests were generated")
	return tests
}

func createTestCase(t *testing.T, tc Case, sourceURI, eventName string, failDep *dependency) *Test {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	event := eventName

	inputsRaw := createInputs(t, tc, sourceURI, &event)

	name := fmt.Sprintf("%s: %d messages", event, tc.NumMessages)
	if tc.RecordsPerMessage > 0 {
		name += fmt.Sprintf(" (%d records each)", tc.RecordsPerMessage)
	}

	if failDep != nil {
		name = fmt.Sprintf("Fail %s - %s : %s", failDep.Name, failDep.failCall, name)
	}

	test := &Test{
		Name:         name,
		SourceURI:    sourceURI,
		Inputs:       inputsRaw,
		Dependencies: []*dependency{},
	}

	return runSimulation(t, tc, ctrl, test, failDep)
}

func runSimulation(t *testing.T, tc Case, ctrl *gomock.Controller, test *Test, failDep *dependency) *Test {
	depAck := &dependency{}
	ack := MockIAckerInit(ctrl, depAck)

	wg := &sync.WaitGroup{}
	inChan := make(chan []model.Input)

	pipeline, willRepublish := getPipelineInit(t, test, tc.Entity, ctrl, ack, depAck, failDep, inChan, wg)

	if willRepublish {
		test.Name = "Retry " + test.Name
	}

	test.Dependencies = append(test.Dependencies, depAck)

	// wait for the successful message
	wg.Add(1)
	if failDep != nil && willRepublish {
		// wait for both the re-publish and pre-acknowledge
		wg.Add(2 * tc.NumMessages)
	}

	handleTest(t, pipeline, ack, wg, test, inChan)

	return test
}

func createInputs(t *testing.T, tc Case, sourceURI string, event *string) (inputsRaw []string) {
	for i := 0; i < tc.NumMessages; i++ {
		var (
			record   interface{}
			seedType interface{}
		)

		seedType = tc.Entity
		if tc.Feed != nil {
			seedType = tc.Feed
		}

		if tc.RecordsPerMessage > 1 {
			entities := make([]interface{}, tc.RecordsPerMessage)
			for j := 0; j < tc.RecordsPerMessage; j++ {
				seed := reflect.New(reflect.TypeOf(seedType).Elem()).Interface()
				require.NoError(t, gofakeit.Struct(seed))
				entities[j] = seed
			}
			record = entities
		} else {
			// TODO: check if update or similar case
			// get number of fields
			// skip faking a portion of them
			require.NoError(t, gofakeit.Struct(seedType))
			record = seedType
		}

		eventName := *event

		if sourceURI == "" && eventName != "" {
			record = map[string]interface{}{
				eventName: record,
			}
		}

		if eventName == "" {
			*event = "Test"
		}

		jsn, err := json.Marshal(record)
		require.NoError(t, err)

		inputsRaw = append(inputsRaw, string(jsn))
	}

	return inputsRaw
}

func getPipelineInit(t *testing.T,
	test *Test,
	ent model.Entity,
	ctrl *gomock.Controller,
	ack actions.IAcker,
	depAck *dependency,
	failDep *dependency,
	inChan chan []model.Input,
	wg *sync.WaitGroup,
) (*pl.Pipeline, bool) {
	var deps []reflect.Value

	method := getConstructorValues(t, ent, "Pipeline")
	paramsNum := method.Type().NumIn()
	pipeDepNames := make([]string, paramsNum)

	for i := 0; i < paramsNum; i++ {
		mi := method.Type().In(i)
		name := mi.Name()[1:] // remove the I
		pipeDepNames[i] = name

		failCall := ""
		if failDep != nil && failDep.Name == name {
			failCall = failDep.failCall
		}

		switch name {
		case "Uploader":
			dep := &dependency{
				failCall: failCall,
			}
			deps = append(deps, reflect.ValueOf(MockUploaderInit(ctrl, dep)))
			test.Dependencies = append(test.Dependencies, dep)
		case "Downloader":
			dep := &dependency{
				failCall: failCall,
			}
			deps = append(deps, reflect.ValueOf(MockDownloaderInit(ctrl, dep, ent)))
			test.Dependencies = append(test.Dependencies, dep)
		case "Repository":
			dep := &dependency{
				failCall: failCall,
			}
			deps = append(deps, reflect.ValueOf(MockRepositoryInit(t, ctrl, dep)))
			test.Dependencies = append(test.Dependencies, dep)
		case "Publisher":
			dep := &dependency{
				failCall: failCall,
			}
			deps = append(deps, reflect.ValueOf(MockPublisherInit(ctrl, dep)))
			test.Dependencies = append(test.Dependencies, dep)
		case "Republisher":
			dep := &dependency{}
			deps = append(deps, reflect.ValueOf(MockRepublisherInit(ctrl, dep, inChan, wg)))
			test.Dependencies = append(test.Dependencies, dep)
		case "Acker":
			deps = append(deps, reflect.ValueOf(ack))
			test.Dependencies = append(test.Dependencies, depAck)
		}
	}

	if pli := method.Call(deps); len(pli) > 0 {
		if pipeline, ok := pli[0].Interface().(*pl.Pipeline); ok {
			var willRepublish bool

			if failDep != nil && failDep.failCall != "" {
			outer:
				for _, act := range pipeline.GetActions() {
					if act.FailureMandate() == model.StopAndRetry {
						for _, depCallName := range act.DepCallNames() {
							if depCallName == failDep.failCall {
								willRepublish = true
								break outer
							}
						}
					}
				}
			}

			return pipeline, willRepublish
		}
	}

	return nil, false
}

func MockIAckerInit(ctrl *gomock.Controller, dep *dependency) *actions.MockIAcker {
	pub := actions.NewMockIAcker(ctrl)
	dep.Name = "Acker"

	pub.EXPECT().AckMessages(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, msgs ...*model.Message) error {
			inputs := make([]interface{}, len(msgs))
			for i := range msgs {
				jsn, err := json.Marshal(msgs[i])
				if err != nil {
					return err
				}
				inputs[i] = string(jsn)
			}

			c := call{
				Method:       "AckMessages",
				ExpectInputs: inputs,
			}

			var err error
			if dep.failCall == c.Method {
				err = fmt.Errorf("failed to acknowledge")
			}

			if err != nil {
				c.ExpectError = err.Error()
			}

			dep.Calls = append(dep.Calls, c)
			return err
		}).AnyTimes()

	return pub
}

func MockUploaderInit(ctrl *gomock.Controller, dep *dependency) *actions.MockIUploader {
	upl := actions.NewMockIUploader(ctrl)
	dep.Name = "Uploader"

	upl.EXPECT().UploadFile(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, key string, r io.Reader) error {
			file, err := io.ReadAll(r)
			if err != nil {
				return err
			}
			c := call{
				Method: "UploadFile",
				ExpectInputs: []interface{}{
					key, string(file),
				},
			}

			if dep.failCall == c.Method {
				err = fmt.Errorf("failed to upload")
			}

			if err != nil {
				c.ExpectError = err.Error()
			}

			dep.Calls = append(dep.Calls, c)
			return err
		}).AnyTimes()
	return upl
}

func MockDownloaderInit(ctrl *gomock.Controller, dep *dependency, entity model.Entity) *actions.MockIDownloader {
	pub := actions.NewMockIDownloader(ctrl)
	dep.Name = "Downloader"

	pub.EXPECT().DownloadFile(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, key string, r io.Writer) error {
			var outputs []model.Entity
			for i := 0; i < 10; i++ {
				if err := gofakeit.Struct(entity); err != nil {
					return err
				}
				outputs = append(outputs, entity)
			}

			jsn, err := json.Marshal(outputs)
			if err != nil {
				return err
			}

			c := call{
				Method:        "DownloadFile",
				ExpectInputs:  []interface{}{key},
				ExpectOutputs: []interface{}{string(jsn)},
			}

			if dep.failCall == c.Method {
				err = fmt.Errorf("failed to download file")
				dep.failCall = ""
			}

			if err != nil {
				c.ExpectError = err.Error()
			}

			dep.Calls = append(dep.Calls, c)
			return err
		}).AnyTimes()

	pub.EXPECT().DownloadFileFromBucket(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, bucket, key string, r io.Writer) error {
			var outputs []model.Entity
			for i := 0; i < 10; i++ {
				if err := gofakeit.Struct(entity); err != nil {
					return err
				}
				outputs = append(outputs, entity)
			}

			jsn, err := json.Marshal(outputs)
			if err != nil {
				return err
			}

			c := call{
				Method:        "DownloadFileFromBucket",
				ExpectInputs:  []interface{}{bucket, key},
				ExpectOutputs: []interface{}{string(jsn)},
			}

			if dep.failCall == c.Method {
				err = fmt.Errorf("failed to download file")
				dep.failCall = ""
			}

			if err != nil {
				c.ExpectError = err.Error()
			}

			dep.Calls = append(dep.Calls, c)
			return err
		}).AnyTimes()
	return pub
}

func MockPublisherInit(ctrl *gomock.Controller, dep *dependency) *actions.MockIPublisher {
	pub := actions.NewMockIPublisher(ctrl)
	dep.Name = "Publisher"

	pub.EXPECT().PublishEvents(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, msgs ...*model.Message) error {
			inputs := make([]interface{}, len(msgs))
			for i := range msgs {
				inputs[i] = msgs[i].GetBody()
			}

			c := call{
				Method:       "PublishEvents",
				ExpectInputs: inputs,
			}

			var err error
			if dep.failCall == c.Method {
				err = fmt.Errorf("failed to publish")
				dep.failCall = ""
			}

			if err != nil {
				c.ExpectError = err.Error()
			}

			dep.Calls = append(dep.Calls, c)
			return err
		}).AnyTimes()

	return pub
}

func overrideEntityKey(entity model.Entity, key model.Key) error {
	entityVal := reflect.ValueOf(entity)
	if entityVal.Kind() == reflect.Ptr {
		entityVal = entityVal.Elem()
	}

	keyType := reflect.TypeOf(key)

	for j := 0; j < entityVal.NumField(); j++ {
		fld := entityVal.Field(j)

		if fld.Type().AssignableTo(keyType) {
			fld.Set(reflect.ValueOf(key))
		}
	}
	return nil
}

func MockRepositoryInit(t *testing.T, ctrl *gomock.Controller, dep *dependency) *actions.MockIRepository {
	rep := actions.NewMockIRepository(ctrl)
	dep.Name = "Repository"

	rep.EXPECT().EntityExists(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, key model.Key) (bool, error) {
			c := call{
				Method:       "EntityExists",
				ExpectInputs: []interface{}{model.StringifyKey(key)},
			}

			var (
				err    error
				exists bool
			)

			if dep.failCall == c.Method {
				dep.failCall = ""
				exists = true
				c.ExpectOutputs = []interface{}{true}
			} else {
				c.ExpectOutputs = []interface{}{false}
			}

			dep.Calls = append(dep.Calls, c)
			// err = fmt.Errorf("failed to check if entity exists")
			// c.ExpectError = err.Error()

			return exists, err
		}).AnyTimes()

	rep.EXPECT().GetEntity(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, key model.Key, i interface{}) error {
			// fill in the struct with new fake data
			require.NoError(t, gofakeit.Struct(i))
			// override the key portion of it, so it matches the key from the previously generated one
			require.NoError(t, overrideEntityKey(i.(model.Entity), key))

			jsn, err := json.Marshal(i)
			require.NoError(t, err)

			c := call{
				Method:        "GetEntity",
				ExpectInputs:  []interface{}{model.StringifyKey(key)},
				ExpectOutputs: []interface{}{string(jsn)},
			}

			if dep.failCall == c.Method {
				err = fmt.Errorf("failed to get entity")
				dep.failCall = ""
			}

			if err != nil {
				c.ExpectError = err.Error()
			}

			dep.Calls = append(dep.Calls, c)
			return err
		}).AnyTimes()

	rep.EXPECT().SaveEntities(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, entities ...model.Entity) error {
			inputs := make([]interface{}, len(entities))
			for i := range entities {
				jsn, err := json.Marshal(entities[i])
				require.NoError(t, err)
				inputs[i] = string(jsn)
			}

			c := call{
				Method:       "SaveEntities",
				ExpectInputs: inputs,
			}

			var err error

			if dep.failCall == c.Method {
				err = fmt.Errorf("failed to save")
				dep.failCall = ""
			}

			if err != nil {
				c.ExpectError = err.Error()
			}

			dep.Calls = append(dep.Calls, c)
			return err
		}).AnyTimes()

	return rep
}

func MockRepublisherInit(ctrl *gomock.Controller, dep *dependency, ch chan []model.Input, wg *sync.WaitGroup) *actions.MockIRepublisher {
	rep := actions.NewMockIRepublisher(ctrl)
	dep.Name = "Republisher"

	rep.EXPECT().AckMessages(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, msgs ...model.Input) error {
			defer wg.Done()

			inputs := make([]interface{}, len(msgs))
			for i := range msgs {
				jsn, err := json.Marshal(msgs[i])
				if err != nil {
					return err
				}
				inputs[i] = string(jsn)
			}

			c := call{
				Method:       "AckMessages",
				ExpectInputs: inputs,
			}

			dep.Calls = append(dep.Calls, c)
			return nil
		}).AnyTimes()

	rep.EXPECT().PublishEvents(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, msgs ...model.Input) error {
			inputs := make([]interface{}, len(msgs))
			for i := range msgs {
				inputs[i] = msgs[i].GetBody()
			}

			c := call{
				Method:       "PublishEvents",
				ExpectInputs: inputs,
			}

			go func(m []model.Input) {
				ch <- m
				wg.Done()
			}(msgs)

			dep.Calls = append(dep.Calls, c)

			return nil
		}).AnyTimes()
	return rep
}
