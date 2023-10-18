package testengine

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/go-yaml/yaml"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/zale144/ube/actions"
	"github.com/zale144/ube/handler"
	"github.com/zale144/ube/model"
	pl "github.com/zale144/ube/pipeline"
)

// TODO: testing the retry mechanism ?

const (
	depAcker       = "Acker"
	depRepublisher = "Republisher"
)

// add 'when_to_process' timestamp to business event input, and pass and republish if it's not >= now

type Test struct {
	Name         string        `json:"test_name" yaml:"test_name"`
	SourceURI    string        `json:"source_uri" yaml:"source_uri,omitempty"`
	Inputs       []string      `json:"inputs" yaml:"inputs"`
	Dependencies []*dependency `json:"dependencies" yaml:"dependencies"`
}

const testDir = "tests"

type dependency struct {
	Name  string `json:"name" yaml:"name"`
	Calls []call `json:"calls" yaml:"calls,omitempty"`
}

type call struct {
	Method        string        `json:"method" yaml:"method"`
	ExpectInputs  []interface{} `json:"expect_inputs" yaml:"expect_inputs,omitempty"`
	ExpectOutputs []interface{} `json:"expect_outputs" yaml:"expect_outputs,omitempty"`
	ExpectError   string        `json:"expect_error,omitempty" yaml:"expect_error,omitempty"`
}

// EventHandler handles an event by iterating over its messages
func EventHandler(t *testing.T, ent model.Entity, idx ...int) {
	testFiles, err := os.ReadDir(testDir) // reads the test sorted
	require.NoError(t, err)

	for _, testFile := range testFiles {
		fmt.Printf("\n-----< Test file: %s > -----\n\n", testFile.Name())
		testData, err := os.ReadFile(filepath.Join(testDir, testFile.Name()))
		require.NoError(t, err)

		eventHandler(t, testData, ent, idx...)
	}
}

func eventHandler(t *testing.T, testData []byte, ent model.Entity, idx ...int) {
	var tests []*Test

	require.NoError(t, yaml.Unmarshal(testData, &tests), "parsing test data failed")

	logger := zap.NewExample()
	defer func() { _ = logger.Sync() }()

	undo := zap.ReplaceGlobals(logger)
	defer undo()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	model.Now = func() time.Time { return time.Date(2021, 11, 22, 3, 4, 5, 0, time.UTC) }
	model.UUIDStr = func() string { return "5de1ea04-61c9-4cf8-bdf8-320479e62d31" }

	for i := range tests {
		i := i
		test := tests[i]
		t.Run(test.Name, func(t *testing.T) {
			eventHandlerInnerTest(t, test, ent, ctrl, i, idx)
		})
	}
}

func eventHandlerInnerTest(t *testing.T, test *Test, ent model.Entity, ctrl *gomock.Controller, i int, idx []int) bool {
	if len(test.Inputs) == 0 {
		t.Fatal("cannot run test without inputs")
	}

	if len(idx) > 0 && idx[0] != i {
		return false
	}

	inChan := make(chan []model.Input)
	wg := &sync.WaitGroup{}
	ack := actions.NewMockIAcker(ctrl)
	rep := actions.NewMockIRepublisher(ctrl)

	for _, d := range test.Dependencies {
		if d.Name == depAcker {
			ack = MockAcker(t, ctrl, d)
		}

		if d.Name == depRepublisher {
			rep = MockRepublisher(t, ctrl, d, inChan, wg)
		}
	}

	pipeline := getPipeline(t, test, ent, ctrl, ack, rep)

	wg.Add(1)

	handleTest(t, pipeline, ack, wg, test, inChan)
	return true
}

func handleTest(
	t *testing.T,
	pipeline *pl.Pipeline,
	ack *actions.MockIAcker,
	wg *sync.WaitGroup,
	test *Test,
	inChan chan []model.Input,
) {
	h := handler.NewEventHandler(pipeline, ack)

	go func() {
		defer wg.Done()

		var ins []model.Input
		for j, input := range test.Inputs {
			ins = append(ins, makeInput(j, test.SourceURI, []byte(input)))
		}
		inChan <- ins
	}()

	go func() {
		wg.Wait()
		close(inChan)
	}()

	for ins := range inChan {
		if err := h.Handle(context.Background(), model.NewInputEvent(ins)); err != nil {
			t.Logf("error from the pipeline: %s", err)
		}

		result := h.GetResult()
		zap.L().Info("pipeline result", zap.Any("result", result))
	}
}

func makeInput(i int, sourceURI string, body []byte) *model.Message {
	return &model.Message{
		ID:        fmt.Sprintf("msg_%d", i+1),
		Reference: fmt.Sprintf("ref_%d", i+1),
		SourceURI: sourceURI,
		Body:      body,
	}
}

func getPipeline(
	t *testing.T,
	test *Test,
	ent model.Entity,
	ctrl *gomock.Controller,
	ack actions.IAcker,
	rep actions.IRepublisher,
) *pl.Pipeline {
	var deps []reflect.Value

	method := getConstructorValues(t, ent, "Pipeline")
	testDepNum := len(test.Dependencies)
	paramsNum := method.Type().NumIn()
	require.Equalf(
		t,
		testDepNum,
		paramsNum+1, /*(account for acker)*/
		"number of pipeline dependencies ('%d') is not the same as test dependencies ('%d')",
		paramsNum,
		testDepNum,
	)

	testDepNames := make([]string, testDepNum)
	for i := range test.Dependencies {
		testDepNames[i] = test.Dependencies[i].Name
	}

	pipeDepNames := make([]string, paramsNum)

	for i := 0; i < paramsNum; i++ {
		mi := method.Type().In(i)
		name := mi.Name()[1:] // remove the I
		pipeDepNames[i] = name

		switch name {
		case "Uploader":
			deps = append(deps, reflect.ValueOf(MockUploader(t, ctrl, test.Dependencies[i])))
		case "Downloader":
			deps = append(deps, reflect.ValueOf(MockDownloader(t, ctrl, test.Dependencies[i])))
		case "Repository":
			deps = append(deps, reflect.ValueOf(MockRepository(t, ctrl, test.Dependencies[i])))
		case "Publisher":
			deps = append(deps, reflect.ValueOf(MockPublisher(t, ctrl, test.Dependencies[i])))
		case depRepublisher:
			deps = append(deps, reflect.ValueOf(rep))
		case depAcker:
			deps = append(deps, reflect.ValueOf(ack))
		default:
			assert.Fail(t, fmt.Sprintf("Can't handle name '%s' in pipeline", name))
		}
	}

	pipeDepNames = append(pipeDepNames, depAcker) // account for acker

	require.EqualValuesf(t, pipeDepNames, testDepNames, "the pipeline dependencies don't match the test dependencies")

	if pli := method.Call(deps); len(pli) > 0 {
		if pipeline, ok := pli[0].Interface().(*pl.Pipeline); ok {
			return pipeline
		}
	}

	return nil
}

func getConstructorValues(t *testing.T, ent model.Entity, constructorName string) reflect.Value {
	entVal := reflect.ValueOf(ent)
	method := entVal.MethodByName(constructorName)
	require.Truef(t, method.IsValid(), "entity '%s' does not have a pipeline method", entVal.Type().String())

	return method
}

func MockAcker(t *testing.T, ctrl *gomock.Controller, dep *dependency) *actions.MockIAcker {
	ack := actions.NewMockIAcker(ctrl)

	for i := range dep.Calls {
		depCall := dep.Calls[i]

		ack.EXPECT().AckMessages(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, msgs ...*model.Message) error {
				for j, msg := range msgs {
					jsn, err := msg.MarshalJSON()
					require.NoError(t, err)
					zap.L().Info("acknowledging message", zap.Any("msg", json.RawMessage(jsn)))
					if len(depCall.ExpectInputs) > 0 {
						assert.JSONEqf(t, depCall.ExpectInputs[j].(string), string(jsn), "AckMessages() got: %s", string(jsn))
					}
				}

				var err error
				if depCall.ExpectError != "" {
					err = errors.New(depCall.ExpectError)
				}

				return err
			}).Times(1)
	}

	return ack
}

func MockRepublisher(
	t *testing.T,
	ctrl *gomock.Controller,
	dep *dependency,
	ch chan []model.Input,
	wg *sync.WaitGroup,
) *actions.MockIRepublisher {
	numCalls := len(dep.Calls)
	rep := actions.NewMockIRepublisher(ctrl)

	wg.Add(numCalls)

	for i := range dep.Calls {
		depCall := dep.Calls[i]
		switch depCall.Method {
		case "AckMessages":
			rep = mockAckMessages(t, rep, depCall, wg)
		case "PublishEvents":
			rep = mockPublishEvents(t, rep, depCall, wg, ch)
		}
	}

	return rep
}

func mockAckMessages(t *testing.T, rep *actions.MockIRepublisher, depCall call, wg *sync.WaitGroup) *actions.MockIRepublisher {
	rep.EXPECT().AckMessages(gomock.Any(), gomock.Any()).Do(func(ctx context.Context, msgs ...*model.Message) error {
		defer wg.Done()

		for j, msg := range msgs {
			jsn, err := msg.MarshalJSON()
			require.NoError(t, err)
			zap.L().Info("pre-acknowledging message", zap.Any("msg", json.RawMessage(jsn)))
			if len(depCall.ExpectInputs) > 0 {
				assert.JSONEqf(t, depCall.ExpectInputs[j].(string), string(jsn), "AckMessages() got: %s", string(jsn))
			}
		}

		return nil
	}).Times(1)

	return rep
}

func mockPublishEvents(
	t *testing.T,
	rep *actions.MockIRepublisher,
	depCall call,
	wg *sync.WaitGroup,
	ch chan []model.Input,
) *actions.MockIRepublisher {
	rep.EXPECT().PublishEvents(gomock.Any(), gomock.Any()).Do(
		func(_ context.Context, msgs ...*model.Message) error {
			for j, msg := range msgs {
				zap.L().Info("re-publishing message", zap.Any("msg", json.RawMessage(msg.GetBody())))
				if len(depCall.ExpectInputs) > 0 {
					assert.JSONEqf(t, depCall.ExpectInputs[j].(string), msg.GetBody(), "PublishEvents() got: %s", msg.GetBody())
				}

				go func(m *model.Message) {
					ch <- []model.Input{m}
					wg.Done()
				}(msg)
			}
			return nil
		}).Times(1)

	return rep
}

func MockUploader(t *testing.T, ctrl *gomock.Controller, dep *dependency) *actions.MockIUploader {
	upl := actions.NewMockIUploader(ctrl)

	for i := range dep.Calls {
		depCall := dep.Calls[i]

		upl.EXPECT().UploadFile(gomock.Any(), gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, key string, r io.Reader) error {
				if len(depCall.ExpectInputs) > 0 {
					jsn, err := io.ReadAll(r)
					if err != nil {
						return err
					}

					assert.Equal(t, depCall.ExpectInputs[0], key)
					assert.JSONEqf(t, depCall.ExpectInputs[1].(string), string(jsn), "UploadFile() got: %s", string(jsn))
				}

				var err error
				if depCall.ExpectError != "" {
					err = errors.New(depCall.ExpectError)
				}

				return err
			}).Times(1)
	}

	return upl
}

func MockDownloader(t *testing.T, ctrl *gomock.Controller, dep *dependency) *actions.MockIDownloader {
	pub := actions.NewMockIDownloader(ctrl)

	for i := range dep.Calls {
		depCall := dep.Calls[i]

		pub.EXPECT().DownloadFileFromBucket(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			DoAndReturn(
				func(_ context.Context, bucket, key string, body io.Writer) error {
					err := mockDownloaderCheckInputs(depCall, bucket, key)
					if err != nil {
						return err
					}

					if len(depCall.ExpectOutputs) > 0 {
						jsn := []byte(depCall.ExpectOutputs[0].(string))
						_, err = body.Write(jsn)
						require.NoError(t, err)
						zap.L().Info("got file", zap.String("key", key), zap.Any("file", json.RawMessage(jsn)))
					}

					if depCall.ExpectError != "" {
						err = errors.New(depCall.ExpectError)
					}

					return err
				}).Times(1)
	}

	return pub
}

func mockDownloaderCheckInputs(depCall call, bucket, key string) error {
	if len(depCall.ExpectInputs) > 1 {
		if depCall.ExpectInputs[0].(string) != bucket {
			return fmt.Errorf("bucket '%s' not found", bucket)
		}

		if depCall.ExpectInputs[1].(string) != key {
			return fmt.Errorf("file with key '%s' not found in bucket '%s'", key, bucket)
		}
	}

	return nil
}

func MockPublisher(t *testing.T, ctrl *gomock.Controller, dep *dependency) *actions.MockIPublisher {
	pub := actions.NewMockIPublisher(ctrl)

	for i := range dep.Calls {
		depCall := dep.Calls[i]

		pub.EXPECT().PublishEvents(gomock.Any(), gomock.Any()).
			DoAndReturn(
				func(_ context.Context, msgs ...*model.Message) error {
					for i, msg := range msgs {
						zap.L().Info("publishing message", zap.Any("msg", json.RawMessage(msg.GetBody())))
						if len(depCall.ExpectInputs) == len(msgs) {
							assert.JSONEqf(t, depCall.ExpectInputs[i].(string), msg.GetBody(), "PublishEvents() got: %s", msg.GetBody())
						}
					}

					var err error
					if depCall.ExpectError != "" {
						err = errors.New(depCall.ExpectError)
					}

					return err
				}).Times(1)
	}

	return pub
}

func MockRepository(t *testing.T, ctrl *gomock.Controller, dep *dependency) *actions.MockIRepository {
	pub := actions.NewMockIRepository(ctrl)

	for i := range dep.Calls {
		depCall := dep.Calls[i]
		switch depCall.Method {
		case "GetEntity":
			pub = mockGetEntity(t, pub, depCall)
		case "EntityExists":
			pub = mockEntityExists(t, pub, depCall)
		case "SaveEntities":
			pub = mockSaveEntities(t, pub, depCall)
		}
	}

	return pub
}

func mockGetEntity(t *testing.T, pub *actions.MockIRepository, depCall call) *actions.MockIRepository {
	pub.EXPECT().GetEntity(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(
			func(_ context.Context, key model.Key, i interface{}) error {
				err := mockGetEntityCheckInputs(depCall, key)
				if err != nil {
					return err
				}

				mockGetEntityCheckOutputs(t, depCall, key, i)

				if depCall.ExpectError != "" {
					err = errors.New(depCall.ExpectError)
				}

				return err
			}).Times(1)

	return pub
}

func mockGetEntityCheckInputs(depCall call, key model.Key) error {
	if len(depCall.ExpectInputs) > 0 {
		expectStr := fmt.Sprint(depCall.ExpectInputs[0])
		parts := strings.Split(expectStr, model.KeySeparator)

		pk := parts[0]
		if pk != key.PK() {
			return fmt.Errorf("entity for PK '%s' not found", pk)
		}

		if sk, ok := key.(model.SK); ok {
			if len(parts) <= 1 {
				return fmt.Errorf("SK for entity not expected")
			}

			if parts[1] != sk.SK() {
				return fmt.Errorf("entity for SK '%s' not found", sk.SK())
			}
		}
	}

	return nil
}

func mockGetEntityCheckOutputs(t *testing.T, depCall call, key model.Key, i interface{}) {
	jsn, ok := depCall.ExpectOutputs[0].(string)
	require.Truef(t, ok, "the output should be a string")
	err := json.Unmarshal([]byte(jsn), i)
	require.NoErrorf(t, err, "the output is not a proper json")
	zap.L().Info("got entity", zap.String("key", model.StringifyKey(key)), zap.Any("entity", json.RawMessage(jsn)))
}

func mockEntityExists(t *testing.T, pub *actions.MockIRepository, depCall call) *actions.MockIRepository {
	pub.EXPECT().EntityExists(gomock.Any(), gomock.Any()).
		DoAndReturn(
			func(_ context.Context, key model.Key) (bool, error) {
				isOk, err := mockEntityExistsCheckInputs(depCall, key)
				if !isOk {
					return isOk, err
				}

				var exists, ok bool

				if len(depCall.ExpectOutputs) > 0 {
					exists, ok = depCall.ExpectOutputs[0].(bool)
					require.Truef(t, ok, "the output should be a bool")
				}

				if depCall.ExpectError != "" {
					err = errors.New(depCall.ExpectError)
				}

				return exists, err
			}).Times(1)

	return pub
}

func mockEntityExistsCheckInputs(depCall call, key model.Key) (bool, error) {
	if len(depCall.ExpectInputs) > 0 {
		expectStr := fmt.Sprint(depCall.ExpectInputs[0])
		parts := strings.Split(expectStr, model.KeySeparator)

		pk := parts[0]
		if pk != key.PK() {
			return false, nil
		}

		if sk, ok := key.(model.SK); ok {
			if len(parts) <= 1 {
				return false, fmt.Errorf("SK for entity not expected")
			}

			if parts[1] != sk.SK() {
				return false, nil
			}
		}
	}

	return true, nil
}

func mockSaveEntities(t *testing.T, pub *actions.MockIRepository, depCall call) *actions.MockIRepository {
	pub.EXPECT().SaveEntities(gomock.Any(), gomock.Any()).
		DoAndReturn(
			func(_ context.Context, entities ...model.Entity) error {
				require.Equalf(
					t,
					len(depCall.ExpectInputs),
					len(entities),
					"SaveEntities(): number of inputs should be equal to the number of entities",
				)
				for i := range entities {
					entity := entities[i]
					jsn, err := json.MarshalIndent(entity, "", "	")
					require.NoError(t, err)
					zap.L().Info("saving entity", zap.String("key", model.StringifyKey(entity.GetKey())), zap.Any("entity", json.RawMessage(jsn)))
					assert.JSONEqf(t, depCall.ExpectInputs[i].(string), string(jsn), "SaveEntities() got: %s", string(jsn))
				}

				var err error
				if depCall.ExpectError != "" {
					err = errors.New(depCall.ExpectError)
				}

				return err
			}).Times(1)

	return pub
}
