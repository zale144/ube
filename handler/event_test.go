package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/zale144/ube/actions"
	"github.com/zale144/ube/model"
	pl "github.com/zale144/ube/pipeline"
)

type Model struct {
	Field string `json:"field"`
}

func (mdl Model) GetKey() model.Key {
	return mdl
}

func (mdl Model) PK() string {
	return mdl.Field
}

func TestMain(m *testing.M) {
	m.Run()
}

type fakeAction struct {
	actions.Base
}

func (a fakeAction) Name() string                             { return "fake" }
func (a fakeAction) Process(context.Context, ...model.Medium) {}
func (a fakeAction) DepCallNames() []string                   { return nil }

func TestEventHandler_process_record_fail(t *testing.T) {
	p := pl.NewPipeline(&Model{}, pl.Action(fakeAction{}))

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ack := actions.NewMockIAcker(ctrl)
	ack.EXPECT().AckMessages(gomock.Any(), gomock.Any()).AnyTimes()

	h := NewEventHandler(p, ack)

	msgs := []model.Input{&model.Message{ID: "1234", Reference: "ref1", Body: json.RawMessage("{}")}}
	in := model.NewInputEvent(msgs)

	var ctx context.Context
	err := h.Handle(ctx, in)

	assert.NoError(t, err)
}

func TestNewEventHandler(t *testing.T) {
	type args struct {
		pipeline iPipeline
		acker    actions.IAcker
	}
	tests := []struct {
		name string
		args args
		want EventHandler
	}{
		{
			name: "success",
			args: args{
				pipeline: pl.NewPipeline(&Model{}, pl.Action(&fakeAction{})),
				acker:    &actions.MockIAcker{},
			},
			want: EventHandler{
				pipeline: pl.NewPipeline(&Model{}, pl.Action(&fakeAction{})),
				acker:    &actions.MockIAcker{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewEventHandler(tt.args.pipeline, tt.args.acker), "NewEventHandler(%v, %v)", tt.args.pipeline, tt.args.acker)
		})
	}
}

func TestEventHandler_Handle(t *testing.T) {
	type args struct {
		ctx context.Context
		ev  *model.InputEvent
	}
	tests := []struct {
		name            string
		args            args
		mockExpectation func(mp *MockPipeline, ma *actions.MockIAcker)
		wantErr         assert.ErrorAssertionFunc
	}{
		{
			name: "success",

			args: args{
				ctx: context.Background(),
				ev: model.NewInputEvent([]model.Input{
					&model.Message{
						ID:        "123",
						Reference: "ref1",
						SourceURI: "src1",
						Body:      []byte{},
					},
				}),
			},
			mockExpectation: func(mp *MockPipeline, ma *actions.MockIAcker) {
				mp.EXPECT().InvokePipeline(gomock.Any(), gomock.Any()).
					Return(pl.EventProcessingResult{
						Status: "ok",
						Errors: nil,
						BusinessEvents: []model.PipelineMedium{
							&model.BusinessEvent{
								Event: &model.Event{
									ID:        "123",
									Reference: "ref1",
								},
							},
						},
					}, nil)
				ma.EXPECT().AckMessages(gomock.Any(), gomock.Any())
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
		{
			name: "failed pipeline",

			args: args{
				ctx: context.Background(),
				ev: model.NewInputEvent([]model.Input{
					&model.Message{
						ID:        "123",
						Reference: "ref1",
						SourceURI: "src1",
						Body:      []byte{},
					},
				}),
			},
			mockExpectation: func(mp *MockPipeline, ma *actions.MockIAcker) {
				mp.EXPECT().InvokePipeline(gomock.Any(), gomock.Any()).
					Return(pl.EventProcessingResult{
						Status: "failure",
						Errors: nil,
						BusinessEvents: []model.PipelineMedium{
							&model.BusinessEvent{
								Event: &model.Event{
									ID:        "123",
									Reference: "ref1",
								},
								Error: fmt.Errorf("failed pipeline"),
							},
						},
					}, fmt.Errorf("failed pipeline"))
				ma.EXPECT().AckMessages(gomock.Any(), gomock.Any())
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
		{
			name: "got errors",

			args: args{
				ctx: context.Background(),
				ev: model.NewInputEvent([]model.Input{
					&model.Message{
						ID:        "123",
						Reference: "ref1",
						SourceURI: "src1",
						Body:      []byte{},
					},
				}),
			},
			mockExpectation: func(mp *MockPipeline, ma *actions.MockIAcker) {
				mp.EXPECT().InvokePipeline(gomock.Any(), gomock.Any()).
					Return(pl.EventProcessingResult{
						Status: "failure",
						BusinessEvents: []model.PipelineMedium{
							&model.BusinessEvent{
								Event: &model.Event{
									ID:        "123",
									Reference: "ref1",
								},
							},
						},
						Errors: []string{"failed an action"},
					}, nil)
				ma.EXPECT().AckMessages(gomock.Any(), gomock.Any())
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err != nil
			},
		},
		{
			name: "nothing to acknowledge",

			args: args{
				ctx: context.Background(),
				ev: model.NewInputEvent([]model.Input{
					&model.Message{
						ID:        "123",
						Reference: "ref1",
						SourceURI: "src1",
						Body:      []byte{},
					},
				}),
			},
			mockExpectation: func(mp *MockPipeline, ma *actions.MockIAcker) {
				mp.EXPECT().InvokePipeline(gomock.Any(), gomock.Any()).
					Return(pl.EventProcessingResult{
						Status: "failure",
						BusinessEvents: []model.PipelineMedium{
							&model.BusinessEvent{
								Event: &model.Event{},
							},
						},
					}, nil)
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
		{
			name: "failed to acknowledge",

			args: args{
				ctx: context.Background(),
				ev: model.NewInputEvent([]model.Input{
					&model.Message{
						ID:        "123",
						Reference: "ref1",
						SourceURI: "src1",
						Body:      []byte{},
					},
				}),
			},
			mockExpectation: func(mp *MockPipeline, ma *actions.MockIAcker) {
				mp.EXPECT().InvokePipeline(gomock.Any(), gomock.Any()).
					Return(pl.EventProcessingResult{
						Status: "failure",
						BusinessEvents: []model.PipelineMedium{
							&model.BusinessEvent{
								Event: &model.Event{
									ID:        "123",
									Reference: "ref1",
								},
							},
						},
					}, nil)
				ma.EXPECT().AckMessages(gomock.Any(), gomock.Any()).Return(fmt.Errorf("failed to ack"))
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err != nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ack := actions.NewMockIAcker(ctrl)
			pipe := NewMockPipeline(ctrl)

			tt.mockExpectation(pipe, ack)

			p := NewEventHandler(pipe, ack)
			err := p.Handle(tt.args.ctx, tt.args.ev)
			if !tt.wantErr(t, err, fmt.Sprintf("Handle(%v, %v)", tt.args.ctx, tt.args.ev)) {
				t.Errorf("unexpected error: %s", err)
			}
		})
	}
}

func TestEventHandler_GetResult(t *testing.T) {
	type fields struct {
		result pl.EventProcessingResult
	}
	tests := []struct {
		name   string
		fields fields
		want   pl.EventProcessingResult
	}{
		{
			name: "success",
			fields: fields{
				result: pl.EventProcessingResult{
					Status:         "ok",
					Errors:         nil,
					BusinessEvents: []model.PipelineMedium{},
				},
			},
			want: pl.EventProcessingResult{
				Status:         "ok",
				Errors:         nil,
				BusinessEvents: []model.PipelineMedium{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := EventHandler{
				result: tt.fields.result,
			}
			assert.Equalf(t, tt.want, p.GetResult(), "GetResult()")
		})
	}
}
