package actions

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"go.uber.org/zap"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/zale144/ube/libs/converter"

	"github.com/zale144/ube/model"
)

type (
	InputTransform struct {
		depCallNames []string
		Transforms   []TransformFn
		Base
	}
	TransformFn func(ctx context.Context, be model.InputActionMedium) (model.InputActionMedium, int, error)
)

// InputTransformer constructs a new InputTransform
func InputTransformer(options ...TransformOption) *InputTransform {
	it := &InputTransform{
		Base: Base{
			critical:       true,
			failureMandate: model.StopAndRaiseError,
			batchSize:      100,
		},
	}

	for _, opt := range options {
		opt(it)
	}

	return it
}

func (e *InputTransform) AddInputTransformer(options ...TransformOption) {
	for _, opt := range options {
		opt(e)
	}
}

// Process implements the action interface in UBE, executes the underlying embedded device
func (e InputTransform) Process(ctx context.Context, bes ...model.Medium) {
	if len(bes) == 0 {
		return
	}

	be0 := bes[0]

	if len(be0.GetEntities()) == 0 {
		return
	}

	entType := reflect.TypeOf(be0.GetEntities()[0])
	if entType.Kind() == reflect.Ptr {
		entType = entType.Elem()
	}

	var counter int

	for i := range bes {
		be, ok := bes[i].(model.InputActionMedium)
		if !ok {
			be.SetError(fmt.Errorf("unexpected input type %T", bes[0]))
			return
		}

		// TODO: how to figure out which is already business event??
		if be.GetEventName() != "" {
			continue
		}

		for _, transform := range e.Transforms {
			var (
				nbe   model.InputActionMedium
				count int
				err   error
			)
			nbe, count, err = transform(ctx, be)
			if err != nil {
				nbe.SetError(fmt.Errorf("transform event '%s' fail: %w", be.GetID(), err))
			} else {
				counter += count
			}
			be = nbe
		}

		if err := be.InitEntity(entType); err != nil {
			be.SetError(fmt.Errorf("failed to transform business events: %w", err))
			return
		}

		be.SetBody(nil)

		bes[i] = be
	}

	zap.L().Info("entities transformed", zap.Int("size", counter))
}

func (e InputTransform) DepCallNames() []string {
	return e.depCallNames
}

func (InputTransform) Name() string {
	return "InputTransformer"
}

func CreateEvent(category, source string) TransformOption {
	return func(transform *InputTransform) {
		transform.Transforms = append(transform.Transforms, func(
			ctx context.Context,
			be model.InputActionMedium,
		) (model.InputActionMedium, int, error) {
			be.SetEventCategory(category)
			caser := cases.Title(language.Und, cases.NoLower)
			be.SetEventName("Create" + caser.String(category))
			//be.BaseWarehouse = source // TODO: ?
			be.SetSource(source)

			return be, 1, nil
		})
	}
}

type TransformOption func(*InputTransform)

func EventFromQueueSource(category, source string, eventNameMap map[string]string) TransformOption {
	return func(transform *InputTransform) {
		transform.Transforms = append(transform.Transforms, func(
			ctx context.Context,
			be model.InputActionMedium,
		) (model.InputActionMedium, int, error) {
			be.SetEventCategory(category)
			//be.BaseWarehouse = source // TODO: ?
			be.SetSource(source)

			eventName, ok := eventNameMap[be.GetEventSource()]
			if !ok {
				err := fmt.Errorf("event name mapping not defined for provided source URI: '%s'", be.GetEventSource())
				return be, 0, err
			}

			be.SetEventName(eventName)

			return be, 1, nil
		})
	}
}

func RecordsFromFilePointer(fileKey, category, source, eventName string, downloader IDownloader) TransformOption {
	return func(transform *InputTransform) {
		transform.depCallNames = append(transform.depCallNames, "DownloadFileFromBucket")
		transform.Transforms = append(transform.Transforms, func(
			ctx context.Context,
			be model.InputActionMedium,
		) (model.InputActionMedium, int, error) {
			be.SetEventCategory(category)
			be.SetEventName(eventName)
			//be.BaseWarehouse = source // TODO: ?
			be.SetSource(source)

			m := make(map[string]interface{})
			byt := be.GetBody()
			if err := json.Unmarshal(byt, &m); err != nil {
				return be, 0, err
			}

			if len(m) > 1 {
				return be, 0, nil
			}

			body, ok := m[fileKey]
			if !ok {
				return be, 0, nil
			}

			fileInfo, ok := body.(map[string]interface{})
			if !ok {
				return be, 0, nil
			}

			bucket, ok := fileInfo["bucket"].(string)
			if !ok {
				return be, 0, nil
			}
			key, ok := fileInfo["key"].(string)
			if !ok {
				return be, 0, nil
			}

			buf := bytes.NewBuffer(nil)
			if err := downloader.DownloadFileFromBucket(context.Background(), bucket, key, buf); err != nil {
				return be, 0, err
			}
			be.SetBody(buf.Bytes())

			return be, 1, nil
		})
	}
}

func RecordsFromKey(category, source string) TransformOption {
	return func(transform *InputTransform) {
		transform.Transforms = append(transform.Transforms, func(
			ctx context.Context,
			be model.InputActionMedium,
		) (model.InputActionMedium, int, error) {
			m := make(map[string]interface{})
			byt := be.GetBody()
			if err := json.Unmarshal(byt, &m); err != nil {
				return be, 0, nil
			}

			if len(m) > 1 {
				return be, 0, nil
			}

			var (
				body   interface{}
				evName string
			)
			for evName, body = range m {
				break
			}

			be.SetEventName(evName)

			jsn, err := json.Marshal(body)
			if err != nil {
				return be, 0, err
			}

			be.SetBody(jsn)

			be.SetEventCategory(category)
			//be.BaseWarehouse = source
			be.SetSource(source)
			return be, 1, nil
		})
	}
}

type convertible interface {
	ConvertToModel() (model.Entity, error)
}

func isList(body []byte) (int, bool) {
	var list []interface{}
	if err := json.Unmarshal(body, &list); err != nil {
		return 0, false
	}

	return len(list), true
}

func FeedToModel(source interface{}, dest model.Entity) TransformOption {
	return func(transform *InputTransform) {
		transform.Transforms = append(transform.Transforms, func(
			ctx context.Context,
			be model.InputActionMedium,
		) (model.InputActionMedium, int, error) {
			srcType := reflect.TypeOf(source)
			if srcType.Kind() == reflect.Ptr {
				srcType = srcType.Elem()
			}
			destType := reflect.TypeOf(dest)
			if destType.Kind() == reflect.Ptr {
				destType = destType.Elem()
			}

			var err error

			byt := be.GetBody()

			var body interface{}

			l, ok := isList(byt)
			if ok {
				var feeds []interface{}
				for i := 0; i < l; i++ {
					feeds = append(feeds, reflect.New(srcType).Interface())
				}

				if err = json.Unmarshal(byt, &feeds); err != nil {
					return be, 0, fmt.Errorf("failed to parse feed body: %w", err)
				}

				var ents []model.Entity
				for _, feed := range feeds {
					var ent model.Entity
					ent, err = feedToModel(feed, destType)
					if err != nil {
						return be, 0, fmt.Errorf("failed to convert feed to model: %w", err)
					}

					ents = append(ents, ent)
				}
				body = ents
			} else {
				feed := reflect.New(srcType).Interface()
				if err = json.Unmarshal(byt, feed); err != nil {
					return be, 0, fmt.Errorf("failed to parse feed body: %w", err)
				}

				body, err = feedToModel(feed, destType)
				if err != nil {
					return be, 0, fmt.Errorf("failed to convert feed to model: %w", err)
				}
			}

			newBody, err := json.Marshal(body)
			if err != nil {
				return be, 0, fmt.Errorf("failed to serialise model: %w", err)
			}

			be.SetBody(newBody)

			return be, 1, nil
		})
	}
}

func feedToModel(feed interface{}, destType reflect.Type) (ent model.Entity, err error) {
	if conv, ok := feed.(convertible); ok {
		ent, err = conv.ConvertToModel()
	} else {
		entity := reflect.New(destType).Interface().(model.Entity)
		err = converter.ConvertStruct(feed, reflect.ValueOf(entity))
		ent = entity
	}
	if err != nil {
		return
	}

	return
}
