package handler

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/zale144/ube/model"
	pl "github.com/zale144/ube/pipeline"
)

func TestNewRequestHandler(t *testing.T) {
	type args struct {
		pipeline        iPipeline
		outputTransform APIOUTTransform
	}
	tests := []struct {
		name string
		args args
		want RequestHandler
	}{
		{
			name: "success",
			args: args{
				pipeline: pl.NewPipeline(&Model{}, pl.Action(&fakeAction{})),
			},
			want: RequestHandler{
				pipeline: pl.NewPipeline(&Model{}, pl.Action(&fakeAction{})),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewRequestHandler(tt.args.pipeline, tt.args.outputTransform), "NewRequestHandler(%v, %v)", tt.args.pipeline, tt.args.outputTransform)
		})
	}
}

func stubOutputTransform(err error) func(*pl.EventProcessingResult) (*model.Response, error) {
	return func(res *pl.EventProcessingResult) (*model.Response, error) {
		code := http.StatusOK
		if strings.ToLower(res.Status) != "ok" || err != nil {
			code = http.StatusInternalServerError
		}
		return &model.Response{
			StatusCode: code,
			Headers:    nil,
			Body:       "{}",
		}, err
	}
}

func TestRequestHandler_Handle(t *testing.T) {
	type fields struct {
		outputTransform APIOUTTransform
	}
	type args struct {
		ctx context.Context
		req *model.Request
	}
	tests := []struct {
		name            string
		fields          fields
		mockExpectation func(mp *MockPipeline)
		args            args
		want            *model.Response
		wantErr         assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			fields: fields{
				outputTransform: stubOutputTransform(nil),
			},
			args: args{
				ctx: nil,
				req: &model.Request{
					Path:       "/",
					HTTPMethod: http.MethodPost,
					RequestID:  "123",
					Body:       "{}",
					Reference:  model.Reference{},
				},
			},
			mockExpectation: func(mp *MockPipeline) {
				mp.EXPECT().InvokePipeline(gomock.Any(), gomock.Any()).
					Return(pl.EventProcessingResult{
						Status: "ok",
					}, nil)
			},
			want: &model.Response{
				StatusCode: http.StatusOK,
				Body:       "{}",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
		{
			name: "failed pipeline",
			fields: fields{
				outputTransform: stubOutputTransform(nil),
			},
			args: args{
				ctx: nil,
				req: &model.Request{
					Path:       "/",
					HTTPMethod: http.MethodPost,
					RequestID:  "123",
					Body:       "{}",
					Reference:  model.Reference{},
				},
			},
			mockExpectation: func(mp *MockPipeline) {
				mp.EXPECT().InvokePipeline(gomock.Any(), gomock.Any()).
					Return(pl.EventProcessingResult{
						Status: "failure",
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
			},
			want: &model.Response{
				StatusCode: http.StatusInternalServerError,
				Headers:    map[string]string{"Content-type": "application/json"},
				Body:       `{"message": "process record fail: failed pipeline"}`,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
		{
			name: "got errors",
			fields: fields{
				outputTransform: stubOutputTransform(nil),
			},
			args: args{
				ctx: nil,
				req: &model.Request{
					Path:       "/",
					HTTPMethod: http.MethodPost,
					RequestID:  "123",
					Body:       "{}",
					Reference:  model.Reference{},
				},
			},
			mockExpectation: func(mp *MockPipeline) {
				mp.EXPECT().InvokePipeline(gomock.Any(), gomock.Any()).
					Return(pl.EventProcessingResult{
						Status: "failure",
						Errors: []string{"failed an action"},
					}, nil)
			},
			want: &model.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       "{}",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err != nil
			},
		},
		{
			name: "failed to transform output",
			fields: fields{
				outputTransform: stubOutputTransform(fmt.Errorf("error during conversion")),
			},
			args: args{
				ctx: nil,
				req: &model.Request{
					Path:       "/",
					HTTPMethod: http.MethodPost,
					RequestID:  "123",
					Body:       "{}",
					Reference:  model.Reference{},
				},
			},
			mockExpectation: func(mp *MockPipeline) {
				mp.EXPECT().InvokePipeline(gomock.Any(), gomock.Any()).
					Return(pl.EventProcessingResult{
						Status: "ok",
					}, nil)
			},
			want: &model.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       "{}",
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

			pipe := NewMockPipeline(ctrl)

			tt.mockExpectation(pipe)

			p := NewRequestHandler(pipe, tt.fields.outputTransform)
			got, err := p.Handle(tt.args.ctx, tt.args.req)
			if !tt.wantErr(t, err, fmt.Sprintf("Handle(%v, %v)", tt.args.ctx, tt.args.req)) {
				t.Errorf("unexpected error: %s", err)
				return
			}
			assert.Equalf(t, tt.want, got, "Handle(%v, %v)", tt.args.ctx, tt.args.req)
		})
	}
}
