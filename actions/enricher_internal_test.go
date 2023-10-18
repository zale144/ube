package actions

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/imdario/mergo"
	"github.com/stretchr/testify/assert"

	"github.com/zale144/ube/model"
)

func TestMain(m *testing.M) {
	model.Now = func() time.Time {
		return time.Date(2021, 11, 22, 3, 4, 5, 0, time.UTC)
	}

	m.Run()
}

func TestEnrich_Process(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		enrichers []EnricherOption
	}
	type args struct {
		ctx context.Context
		bes []model.Medium
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "success: one entity; sub-entity",
			fields: fields{
				enrichers: []EnricherOption{
					WithSubEntity("store", &storeRepo{}),
				},
			},
			args: args{
				ctx: context.Background(),
				bes: []model.Medium{
					&model.BusinessEvent{
						Event: &model.Event{},
						Entities: []model.Entity{&product{
							Store: &store{},
						}},
					},
				},
			},
		}, {
			name: "failure: one entity; sub-entity - nil sub-entity",
			fields: fields{
				enrichers: []EnricherOption{
					WithSubEntity("store", &storeRepo{}),
				},
			},
			args: args{
				ctx: context.Background(),
				bes: []model.Medium{
					&model.BusinessEvent{
						Event:    &model.Event{},
						Entities: []model.Entity{&product{}},
					},
				},
			},
		}, {
			name: "failure: one entity; sub-entity - entity not pointer",
			fields: fields{
				enrichers: []EnricherOption{
					WithSubEntity("store", &storeRepo{}),
				},
			},
			args: args{
				ctx: context.Background(),
				bes: []model.Medium{
					&model.BusinessEvent{
						Event:    &model.Event{},
						Entities: []model.Entity{&product{}},
					},
				},
			},
		}, {
			name: "success: list entities; sub-entity",
			fields: fields{
				enrichers: []EnricherOption{
					WithSubEntity("store", &storeRepo{}),
				},
			},
			args: args{
				ctx: context.Background(),
				bes: []model.Medium{
					&model.BusinessEvent{
						Event: &model.Event{},
						Entities: []model.Entity{
							&product{
								Store: &store{},
							},
						},
					},
				},
			},
		}, {
			name: "failure: list entities; sub-entity - nil sub-entity",
			fields: fields{
				enrichers: []EnricherOption{
					WithSubEntity("store", &storeRepo{}),
				},
			},
			args: args{
				ctx: context.Background(),
				bes: []model.Medium{
					&model.BusinessEvent{
						Event: &model.Event{},
						Entities: []model.Entity{
							&product{},
						},
					},
				},
			},
		}, {
			name: "failure: one entity; sub-entity - repo fail",
			fields: fields{
				enrichers: []EnricherOption{
					WithSubEntity("store", &storeRepo{
						err: fmt.Errorf("failed to get store"),
					}),
				},
			},
			args: args{
				ctx: context.Background(),
				bes: []model.Medium{
					&model.BusinessEvent{
						Event: &model.Event{},
						Entities: []model.Entity{&product{
							Store: &store{},
						}},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Enricher(tt.fields.enrichers...)
			e.Process(tt.args.ctx, tt.args.bes...)
		})
	}
}

type product struct {
	productKey
	AnotherOne int
	Store      *store
}

type productKey struct {
	SomeField string
}

func (p productKey) PK() string {
	return p.SomeField
}

func (p product) GetKey() model.Key {
	return p.productKey
}

type store struct {
	ID
	Name    string
	Address string
}

type ID int

func (i ID) PK() string {
	return fmt.Sprint(i)
}

func (s store) GetKey() model.Key {
	return s.ID
}

type prodRepo struct {
	p      product
	exists bool
	err    error
}

func (r prodRepo) SaveEntities(context.Context, ...model.Entity) error {
	return nil
}

func (r prodRepo) GetEntity(_ context.Context, _ model.Key, i interface{}) error {
	if r.err != nil {
		return r.err
	}
	*i.(*product) = r.p
	return nil
}

func (r prodRepo) EntityExists(context.Context, model.Key) (bool, error) {
	return r.exists, nil
}

type storeRepo struct {
	s   store
	err error
}

func (r storeRepo) EntityExists(context.Context, model.Key) (bool, error) {
	return false, nil
}

func (r storeRepo) SaveEntities(context.Context, ...model.Entity) error {
	return nil
}

func (r storeRepo) GetEntity(_ context.Context, _ model.Key, i interface{}) error {
	if r.err != nil {
		return r.err
	}
	*i.(*store) = r.s
	return nil
}

func TestEnrichWithSubEntity(t *testing.T) {
	sr := storeRepo{
		s: store{
			ID:      12,
			Name:    "Adidas",
			Address: "Some new address",
		},
	}

	type args struct {
		subEntityName string
		ctx           context.Context
		be            model.Medium
		repo          IRepository
	}
	tests := []struct {
		name        string
		args        args
		expectedBE  model.Medium
		expectedErr error
	}{
		{
			name: "success: entity",
			args: args{
				subEntityName: "store",
				ctx:           context.Background(),
				repo:          sr,
				be: &model.BusinessEvent{
					Event: &model.Event{},
					Entities: []model.Entity{&product{
						productKey: productKey{SomeField: "bla bla"},
						AnotherOne: 234,
						Store: &store{
							ID: 12,
						},
					}},
				},
			},
			expectedBE: &model.BusinessEvent{
				Event: &model.Event{
					Metadata: &model.Metadata{
						LastUpdated: "2021-11-22T03:04:05Z",
					},
				},
				Entities: []model.Entity{&product{
					productKey: productKey{SomeField: "bla bla"},
					AnotherOne: 234,
					Store: &store{
						ID:      12,
						Name:    "Adidas",
						Address: "Some new address",
					},
				}},
			},
			expectedErr: nil,
		}, {
			name: "success: entities",
			args: args{
				subEntityName: "store",
				ctx:           context.Background(),
				repo:          sr,
				be: &model.BusinessEvent{
					Event: &model.Event{},
					Entities: []model.Entity{
						&product{
							productKey: productKey{SomeField: "bla bla"},
							AnotherOne: 234,
							Store: &store{
								ID: 12,
							},
						},
					},
				},
			},
			expectedBE: &model.BusinessEvent{
				Event: &model.Event{
					Metadata: &model.Metadata{
						LastUpdated: "2021-11-22T03:04:05Z",
					},
				},
				Entities: []model.Entity{
					&product{
						productKey: productKey{SomeField: "bla bla"},
						AnotherOne: 234,
						Store: &store{
							ID:      12,
							Name:    "Adidas",
							Address: "Some new address",
						},
					},
				},
			},
			expectedErr: nil,
		}, {
			name: "failure: getEntityValue",
			args: args{
				subEntityName: "store",
				ctx:           context.Background(),
				repo:          sr,
				be: &model.BusinessEvent{
					Event:    &model.Event{},
					Entities: []model.Entity{product{}},
				},
			},
			expectedBE: &model.BusinessEvent{
				Event:    &model.Event{},
				Entities: []model.Entity{product{}},
			},
			expectedErr: fmt.Errorf("entity must be a pointer to a struct"),
		}, {
			name: "failure: entities - getEntityValue",
			args: args{
				subEntityName: "store",
				ctx:           context.Background(),
				repo:          sr,
				be: &model.BusinessEvent{
					Event:    &model.Event{},
					Entities: []model.Entity{product{}},
				},
			},
			expectedBE: &model.BusinessEvent{
				Event:    &model.Event{},
				Entities: []model.Entity{product{}},
			},
			expectedErr: fmt.Errorf("entity must be a pointer to a struct"),
		}, {
			name: "failure: no sub entity",
			args: args{
				subEntityName: "store",
				ctx:           context.Background(),
				repo:          sr,
				be: &model.BusinessEvent{
					Event:    &model.Event{},
					Entities: []model.Entity{&entityNoSub{}},
				},
			},
			expectedBE: &model.BusinessEvent{
				Event:    &model.Event{},
				Entities: []model.Entity{&entityNoSub{}},
			},
			expectedErr: fmt.Errorf("no such field 'Store' in type 'actions.entityNoSub'"),
		}, {
			name: "failure: sub-entity not pointer",
			args: args{
				subEntityName: "store",
				ctx:           context.Background(),
				repo:          sr,
				be: &model.BusinessEvent{
					Event:    &model.Event{},
					Entities: []model.Entity{&entitySubNotPtr{}},
				},
			},
			expectedBE: &model.BusinessEvent{
				Event:    &model.Event{},
				Entities: []model.Entity{&entitySubNotPtr{}},
			},
			expectedErr: fmt.Errorf("sub-entity 'Store' in type 'actions.entitySubNotPtr' must be a pointer"),
		}, {
			name: "failure: sub-entity is nil",
			args: args{
				subEntityName: "store",
				ctx:           context.Background(),
				repo:          sr,
				be: &model.BusinessEvent{
					Event:    &model.Event{},
					Entities: []model.Entity{&product{}},
				},
			},
			expectedBE: &model.BusinessEvent{
				Event:    &model.Event{},
				Entities: []model.Entity{&product{}},
			},
			expectedErr: fmt.Errorf("sub-entity 'Store' in type 'actions.product' is nil"),
		}, {
			name: "failure: sub-entity is not entity",
			args: args{
				subEntityName: "store",
				ctx:           context.Background(),
				repo:          sr,
				be: &model.BusinessEvent{
					Event: &model.Event{},
					Entities: []model.Entity{&entitySubNotModel{
						Store: &struct{}{},
					}},
				},
			},
			expectedBE: &model.BusinessEvent{
				Event: &model.Event{},
				Entities: []model.Entity{&entitySubNotModel{
					Store: &struct{}{},
				}},
			},
			expectedErr: fmt.Errorf("sub-entity 'Store' in type 'actions.entitySubNotModel' is not a model.entity, it is '*struct {}'"),
		}, {
			name: "failure: get sub-entity",
			args: args{
				subEntityName: "store",
				ctx:           context.Background(),
				repo: storeRepo{
					err: fmt.Errorf("failed to get sub-entity"),
				},
				be: &model.BusinessEvent{
					Event: &model.Event{},
					Entities: []model.Entity{&product{
						productKey: productKey{SomeField: "bla bla"},
						AnotherOne: 234,
						Store: &store{
							ID: 12,
						},
					}},
				},
			},
			expectedBE: &model.BusinessEvent{
				Event: &model.Event{},
				Entities: []model.Entity{
					&product{
						productKey: productKey{SomeField: "bla bla"},
						AnotherOne: 234,
						Store: &store{
							ID: 12,
						},
					},
				},
			},
			expectedErr: fmt.Errorf("failed to get sub-entity 'store'"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			be, _, err := WithSubEntity(tt.args.subEntityName, tt.args.repo)(tt.args.ctx, tt.args.be)

			asserts := assert.New(t)
			if tt.expectedErr != nil {
				asserts.EqualError(err, tt.expectedErr.Error())
			} else {
				asserts.NoError(err)
			}

			asserts.EqualValues(tt.expectedBE, be)
		})
	}
}

type entityNoSub struct{}

func (e entityNoSub) GetKey() model.Key {
	return nil
}

type entitySubNotPtr struct {
	Store struct{}
}

func (e entitySubNotPtr) GetKey() model.Key {
	return nil
}

type entitySubNotModel struct {
	Store *struct{}
}

func (e entitySubNotModel) GetKey() model.Key {
	return nil
}

func TestWithPatchOriginal(t *testing.T) {
	pr := prodRepo{
		p: product{
			productKey: productKey{SomeField: "bla"},
			AnotherOne: 2223,
			Store: &store{
				ID:      12,
				Name:    "Adidas",
				Address: "Some address",
			},
		},
	}

	type args struct {
		beforeFn  func()
		repo      IRepository
		overrides []Override
		ctx       context.Context
		be        model.Medium
	}
	tests := []struct {
		name        string
		args        args
		expectedBE  model.Medium
		expectedErr error
	}{
		{
			name: "success",
			args: args{
				repo: pr,
				overrides: []Override{
					{
						FieldName: "SomeField",
						Value:     "bla bla bla",
					},
				},
				ctx: context.Background(),
				be: &model.BusinessEvent{
					Event: &model.Event{},
					Entities: []model.Entity{
						&product{
							productKey: productKey{SomeField: "bla"},
							AnotherOne: 2223,
							Store: &store{
								ID:      12,
								Name:    "Adidas",
								Address: "Some address",
							},
						},
					},
				},
			},
			expectedBE: &model.BusinessEvent{
				Event: &model.Event{
					Metadata: &model.Metadata{
						LastUpdated: "2021-11-22T03:04:05Z",
					},
				},
				Entities: []model.Entity{
					&product{
						productKey: productKey{SomeField: "bla bla bla"},
						AnotherOne: 2223,
						Store: &store{
							ID:      12,
							Name:    "Adidas",
							Address: "Some address",
						},
					},
				},
			},
			expectedErr: nil,
		}, {
			name: "success: entities",
			args: args{
				repo: pr,
				overrides: []Override{
					{
						FieldName: "SomeField",
						Value:     "bla bla bla",
					},
				},
				ctx: context.Background(),
				be: &model.BusinessEvent{
					Event: &model.Event{},
					Entities: []model.Entity{
						&product{
							productKey: productKey{SomeField: "bla"},
							AnotherOne: 2223,
							Store: &store{
								ID:      12,
								Name:    "Adidas",
								Address: "Some address",
							},
						},
					},
				},
			},
			expectedBE: &model.BusinessEvent{
				Event: &model.Event{
					Metadata: &model.Metadata{
						LastUpdated: "2021-11-22T03:04:05Z",
					},
				},
				Entities: []model.Entity{
					&product{
						productKey: productKey{SomeField: "bla bla bla"},
						AnotherOne: 2223,
						Store: &store{
							ID:      12,
							Name:    "Adidas",
							Address: "Some address",
						},
					},
				},
			},
			expectedErr: nil,
		}, {
			name: "failure: getEntityValue",
			args: args{
				repo: pr,
				overrides: []Override{
					{
						FieldName: "SomeField",
						Value:     "bla bla bla",
					},
				},
				ctx: context.Background(),
				be: &model.BusinessEvent{
					Event:    &model.Event{},
					Entities: []model.Entity{product{}},
				},
			},
			expectedBE: &model.BusinessEvent{
				Event:    &model.Event{},
				Entities: []model.Entity{product{}},
			},
			expectedErr: fmt.Errorf("get business event value fail: entity must be a pointer to a struct"),
		}, {
			name: "failure: entities - getEntityValue",
			args: args{
				repo: pr,
				overrides: []Override{
					{
						FieldName: "SomeField",
						Value:     "bla bla bla",
					},
				},
				ctx: context.Background(),
				be: &model.BusinessEvent{
					Event: &model.Event{},
					Entities: []model.Entity{
						product{},
					},
				},
			},
			expectedBE: &model.BusinessEvent{
				Event:    &model.Event{},
				Entities: []model.Entity{product{}},
			},
			expectedErr: fmt.Errorf("get business event value fail: entity must be a pointer to a struct"),
		}, {
			name: "failure: get entity",
			args: args{
				repo: prodRepo{
					err: fmt.Errorf("failed to get product"),
				},
				overrides: []Override{
					{
						FieldName: "SomeField",
						Value:     "bla bla bla",
					},
				},
				ctx: context.Background(),
				be: &model.BusinessEvent{
					Event: &model.Event{},
					Entities: []model.Entity{
						&product{},
					},
				},
			},
			expectedBE: &model.BusinessEvent{
				Event:    &model.Event{},
				Entities: []model.Entity{&product{}},
			},
			expectedErr: fmt.Errorf("get entity with key '' fail: failed to get product"),
		}, {
			name: "failure: filed type mismatch",
			args: args{
				repo: pr,
				overrides: []Override{
					{
						FieldName: "SomeField",
						Value:     1,
					},
				},
				ctx: context.Background(),
				be: &model.BusinessEvent{
					Event: &model.Event{},
					Entities: []model.Entity{
						&product{
							productKey: productKey{SomeField: "bla"},
						},
					},
				},
			},
			expectedBE: &model.BusinessEvent{
				Event: &model.Event{},
				Entities: []model.Entity{
					&product{
						productKey: productKey{SomeField: "bla"},
					},
				},
			},
			expectedErr: fmt.Errorf("field 'bla' is not of type 'int', it is 'string'"),
		}, {
			name: "failure: merge",
			args: args{
				beforeFn: func() {
					merge = func(dst, src interface{}, opts ...func(*mergo.Config)) error {
						return fmt.Errorf("failed to merge")
					}
				},
				repo: pr,
				ctx:  context.Background(),
				be: &model.BusinessEvent{
					Event: &model.Event{},
					Entities: []model.Entity{
						&product{
							productKey: productKey{
								SomeField: "bla",
							},
						},
					},
				},
			},
			expectedBE: &model.BusinessEvent{
				Event: &model.Event{},
				Entities: []model.Entity{
					&product{
						productKey: productKey{
							SomeField: "bla",
						},
					},
				},
			},
			expectedErr: fmt.Errorf("enrich original business event fail: failed to merge"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.beforeFn != nil {
				tt.args.beforeFn()
			}

			be, _, err := WithPatchOriginal(tt.args.repo, tt.args.overrides...)(tt.args.ctx, tt.args.be)

			asserts := assert.New(t)
			if tt.expectedErr != nil {
				asserts.EqualError(err, tt.expectedErr.Error())
			} else {
				asserts.NoError(err)
			}

			asserts.EqualValues(tt.expectedBE, be)
		})
	}
}
