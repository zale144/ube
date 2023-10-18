package dynamodb

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/golang/mock/gomock"

	"github.com/zale144/ube/model"
)

func TestDynamoDB_GetEntity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		db        dynamoDB
		tableName string
	}

	type args struct {
		ctx context.Context
		key model.Key
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		wantNil bool
	}{
		{
			name: "success",
			fields: fields{
				db: func() dynamoDB {
					d := NewMockdynamoDB(ctrl)
					d.EXPECT().GetItemWithContext(context.Background(), gomock.Any(), gomock.Any()).Return(
						&dynamodb.GetItemOutput{
							Item: map[string]*dynamodb.AttributeValue{
								"entity": {M: map[string]*dynamodb.AttributeValue{
									"SKU": {S: aws.String("SOMESKU")},
								}},
							},
						}, nil)
					return d
				}(),
				tableName: "product",
			},
			args: args{
				ctx: context.Background(),
				key: key("1"),
			},
			wantErr: false,
			wantNil: false,
		}, {
			name: "not found",
			fields: fields{
				db: func() dynamoDB {
					d := NewMockdynamoDB(ctrl)
					d.EXPECT().GetItemWithContext(context.Background(), gomock.Any(), gomock.Any()).Return(
						&dynamodb.GetItemOutput{}, nil)
					return d
				}(),
				tableName: "product",
			},
			args: args{
				ctx: context.Background(),
				key: key("1"),
			},
			wantErr: true,
			wantNil: true,
		},
	}

	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			p := DynamoDB{
				db:        test.fields.db,
				tableName: test.fields.tableName,
			}

			if err := p.GetEntity(test.args.ctx, test.args.key, &struct{}{}); (err != nil) != test.wantErr {
				t.Errorf("GetEntity() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}

type key string

func (k key) PK() string { return string(k) }
