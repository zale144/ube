package actions

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/imdario/mergo"
	"go.uber.org/zap"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/zale144/ube/model"
)

type (
	// Enrich is a wrapper for enriching the business event be with a provided EnrichFn function
	Enrich struct {
		enrichers []EnricherOption
		Base
	}
	// EnricherMapping maps a string (event name) to a EnrichFn function
	EnricherMapping map[string]EnrichFn
	EnrichFn        func(ctx context.Context, be model.Medium) (model.Medium, int, error)
	GetEntityFn     func(context.Context, string, interface{}) error
	GetSubEntityFn  func(ctx context.Context, table, key string, entity interface{}) error
)

type EnricherOption interface {
	enrich(ctx context.Context, be model.Medium) (model.Medium, int, error)
}

func (e EnrichFn) enrich(ctx context.Context, be model.Medium) (model.Medium, int, error) {
	return e(ctx, be)
}

func (Enrich) Name() string {
	return "Enricher"
}

func (e Enrich) DepCallNames() []string {
	return []string{"GetEntity", "EntityExists"}
}

// Enricher constructs a new
func Enricher(enrichers ...EnricherOption) *Enrich {
	return &Enrich{
		Base: Base{
			batchSize:      1000,
			failureMandate: model.StopAndRaiseError,
		},
		enrichers: enrichers,
	}
}

func (em EnricherMapping) enrich(ctx context.Context, be model.Medium) (model.Medium, int, error) {
	ev := be.GetEventName()
	enrich, ok := em[ev]
	if !ok {
		be.SetError(fmt.Errorf("failed to find enricher for event: '%s'", ev))
		return be, 0, be.GetError()
	}
	return enrich(ctx, be)
}

// Process implements the action interface in UBE, executes the underlying embedded device
func (e Enrich) Process(ctx context.Context, bes ...model.Medium) {
	var counter int

	for i := range bes {
		be := bes[i]

		for _, enr := range e.enrichers {
			var (
				count int
				err   error
			)
			be, count, err = enr.enrich(ctx, be)
			if err != nil {
				be.SetError(fmt.Errorf("enrich business event '%s' fail: %w", be.GetID(), err))
			} else {
				counter += count
			}
			bes[i] = be
		}
	}

	zap.L().Info("entities enriched", zap.Int("size", counter))
}

/*

create: !exist; get nested entities; enrich nested entities

update: get original; enrich only patch attributes; (maybe get and enrich nested attributes)

*/

// EnrichEvent combines a variadic list of parameter enrich functions and returns a single one
func EnrichEvent(actions ...EnrichFn) EnrichFn {
	return func(ctx context.Context, be model.Medium) (model.Medium, int, error) {
		var (
			err   error
			count int
		)

		for _, act := range actions {
			be, count, err = act(ctx, be)
			if err != nil {
				return be, count, fmt.Errorf("enrich event fail: %w", err)
			}
		}

		return be, count, nil
	}
}

// WithSubEntity fetches and sets the nested be
func WithSubEntity(subEntityName string, subRepo IRepository) EnrichFn {
	return func(ctx context.Context, beIfc model.Medium) (model.Medium, int, error) {
		var (
			err   error
			count int
		)

		subEntities := make([]model.Entity, len(beIfc.GetEntities()))
		for i, ent := range beIfc.GetEntities() {
			subEntities[i], err = subEntity(ctx, ent, subEntityName, subRepo)
			if err != nil {
				return beIfc, 0, err
			}
			count++
		}

		if len(subEntities) > 0 {
			beIfc.SetEntities(subEntities)
		}

		be, ok := beIfc.(interface{ UpdateMetadata(func() time.Time) })
		if ok {
			be.UpdateMetadata(model.Now)
		}

		return beIfc, count, err
	}
}

func subEntity(ctx context.Context, ent model.Entity, subEntityName string, subRepo IRepository) (model.Entity, error) {
	val, err := getEntityValue(ent)
	if err != nil {
		return ent, err
	}

	caser := cases.Title(language.Und, cases.NoLower)
	subEntityFieldName := caser.String(subEntityName)
	subEntityField := val.FieldByName(subEntityFieldName)
	if !subEntityField.IsValid() {
		return ent, fmt.Errorf("no such field '%s' in type '%s'", subEntityFieldName, val.Type())
	}

	subi := subEntityField.Interface()
	subEnt := reflect.TypeOf(subi)
	if subEnt.Kind() != reflect.Ptr {
		return ent, fmt.Errorf("sub-entity '%s' in type '%s' must be a pointer", subEntityFieldName, val.Type())
	}

	if subEntityField.IsNil() {
		return ent, fmt.Errorf("sub-entity '%s' in type '%s' is nil", subEntityFieldName, val.Type())
	}

	sub, ok := subi.(model.Entity)
	if !ok {
		return ent, fmt.Errorf("sub-entity '%s' in type '%s' is not a model.entity, it is '%s'",
			subEntityFieldName, val.Type(), subEntityField.Type().String())
	}

	subKey := sub.GetKey()

	if err = subRepo.GetEntity(ctx, subKey, sub); err != nil /*|| sub.GetKey() != subKey */ {
		return ent, fmt.Errorf("failed to get sub-entity '%s'", subEntityName)
	}

	f := subEntityField
	if f.IsValid() && f.CanSet() {
		sev := reflect.ValueOf(sub)
		f.Set(sev)
		zap.L().Info("enrichment WithSubEntity: set sub-Data",
			zap.String("field", subEntityFieldName),
			zap.String("Data", val.Type().String()),
			zap.Any("value", sev.Interface()))
	}

	return ent, nil
}

// Override is an override option for patching the original Data

type Override struct {
	FieldName string
	Value     interface{}
}

// WithPatchOriginal fetches the original be and enriches it with the patch one from the event
func WithPatchOriginal(repo IRepository, overrides ...Override) EnrichFn {
	return func(ctx context.Context, beIfc model.Medium) (model.Medium, int, error) {
		var (
			err   error
			count int
		)

		subEntities := make([]model.Entity, len(beIfc.GetEntities()))
		for i, ent := range beIfc.GetEntities() {
			subEntities[i], err = patchOriginal(ctx, ent, repo, overrides...)
			if err != nil {
				return beIfc, 0, err
			}
			count++
		}

		if len(subEntities) > 0 {
			beIfc.SetEntities(subEntities)
		}

		be, ok := beIfc.(interface{ UpdateMetadata(func() time.Time) })
		if ok {
			be.UpdateMetadata(model.Now)
		}

		return beIfc, count, err
	}
}

var merge = mergo.Merge

func patchOriginal(ctx context.Context, ent model.Entity, repo IRepository, overrides ...Override) (model.Entity, error) {
	originalVal, err := getEntityValue(ent)
	if err != nil {
		return ent, fmt.Errorf("get business event value fail: %w", err)
	}

	entKey := ent.GetKey()
	key := entKey
	original := reflect.New(originalVal.Type()).Interface().(model.Entity)

	if err = repo.GetEntity(ctx, entKey, original); err != nil /*|| entKey != original.GetKey() */ {
		return ent, fmt.Errorf("get entity with key '%s' fail: %w", model.StringifyKey(key), err)
	}

	patch := ent
	patchV := reflect.ValueOf(patch)
	patchVal := patchV.Elem()

	for _, ovrd := range overrides {
		fld := patchVal.FieldByName(ovrd.FieldName)
		if fld.IsValid() && fld.CanSet() {
			ovrdValType := reflect.TypeOf(ovrd.Value)
			if fld.Type() != ovrdValType {
				return ent, fmt.Errorf("field '%s' is not of type '%s', it is '%s'",
					fld.String(), ovrdValType.String(), fld.Type())
			}

			fld.Set(reflect.ValueOf(ovrd.Value))
		}
	}

	if err = merge(original, patch, mergo.WithOverride); err != nil {
		return ent, fmt.Errorf("enrich original business event fail: %w", err)
	}

	zap.L().Info("enrichment WithPatchOriginal: merged original with patch",
		zap.String("Data", originalVal.Type().String()),
		zap.Any("original", original),
		zap.Any("patch", patch))

	patchV.Elem().Set(reflect.ValueOf(original).Elem())

	return original.(model.Entity), nil
}

// WithDedupe ensures there is no existing record with the same be key
func WithDedupe(repo IRepository) EnrichFn {
	return func(ctx context.Context, be model.Medium) (model.Medium, int, error) {
		var count int
		for _, ent := range be.GetEntities() {
			if err := dedupe(ctx, ent.GetKey(), repo); err != nil {
				return be, 0, err
			}
			count++
		}

		return be, count, nil
	}
}

func dedupe(ctx context.Context, key model.Key, repo IRepository) error {
	exists, err := repo.EntityExists(ctx, key)
	if err != nil {
		return fmt.Errorf("get entity fail: %w", err)
	}

	if exists {
		return fmt.Errorf("attempt to create a duplicate for business event ID: '%s'", model.StringifyKey(key))
	}

	zap.L().Info("enrichment WithDedupe: successfully de-duplicated",
		zap.String("Data key", model.StringifyKey(key)))

	return nil
}

func getEntityValue(ent model.Entity) (reflect.Value, error) {
	val := reflect.ValueOf(ent)

	if val.Kind() != reflect.Ptr {
		return reflect.Value{}, errors.New("entity must be a pointer to a struct")
	} else if val.IsNil() {
		return reflect.Value{}, errors.New("entity must be a non-nil struct")
	}

	if val.Kind() == reflect.Interface {
		elm := val.Elem()
		if elm.Kind() == reflect.Ptr && !elm.IsNil() && elm.Elem().Kind() == reflect.Ptr {
			val = elm
		}
	}

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return reflect.Value{}, errors.New("business event is not a struct")
	}

	return val, nil
}
