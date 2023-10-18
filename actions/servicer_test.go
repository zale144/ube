package actions

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zale144/ube/model"
)

func TestService_Process(t *testing.T) {
	type fields struct {
		service IService
	}
	type args struct {
		ctx context.Context
		bes []model.Medium
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []model.Medium
	}{
		{
			name: "success",
			fields: fields{
				service: stubService{},
			},
			args: args{
				ctx: context.Background(),
				bes: []model.Medium{
					&model.BusinessEvent{
						ID: "asdf",
					},
				},
			},
			want: []model.Medium{
				&model.BusinessEvent{
					ID: "asdf",
				},
			},
		},
		{
			name: "failed service",
			fields: fields{
				service: stubService{
					err: fmt.Errorf("failed service"),
				},
			},
			args: args{
				ctx: context.Background(),
				bes: []model.Medium{
					&model.BusinessEvent{
						ID: "asdf",
					},
				},
			},
			want: []model.Medium{
				&model.BusinessEvent{
					ID:    "asdf",
					Error: fmt.Errorf("execute business service fail: %w", fmt.Errorf("failed service")),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Service{
				service: tt.fields.service,
			}
			e.Process(tt.args.ctx, tt.args.bes...)

			assert.Equal(t, tt.want, tt.args.bes)
		})
	}
}

type stubService struct {
	err error
}

func (s stubService) Execute(context.Context, ...model.Medium) error {
	return s.err
}
