package model

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessage_MarshalJSON(t *testing.T) {
	type fields struct {
		ID        string
		Reference string
		SourceURI string
		Body      json.RawMessage
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				ID:        "123",
				Reference: "ref1",
				SourceURI: "src1",
				Body:      []byte(`{}`),
			},
			want: []byte(
				`{
	"id": "123",
	"reference": "ref1",
	"body": {},
	"source_uri": "src1"
}
`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Message{
				ID:        tt.fields.ID,
				Reference: tt.fields.Reference,
				SourceURI: tt.fields.SourceURI,
				Body:      tt.fields.Body,
			}
			got, err := m.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.JSONEq(t, string(got), string(tt.want))
		})
	}
}

func TestNewInputEvent(t *testing.T) {
	type args struct {
		messages []Input
	}
	tests := []struct {
		name string
		args args
		want *InputEvent
	}{
		{
			name: "success",
			args: args{
				messages: []Input{
					&Message{
						ID:        "123",
						Reference: "ref1",
						SourceURI: "src1",
						Body:      []byte(`{}`),
					}},
			},
			want: &InputEvent{
				inputs: []Input{
					&Message{
						ID:        "123",
						Reference: "ref1",
						SourceURI: "src1",
						Body:      []byte(`{}`),
					}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewInputEvent(tt.args.messages), "NewInputEvent(%v)", tt.args.messages)
			assert.Equal(t, tt.want.inputs, tt.args.messages)
			for i, want := range tt.want.inputs {
				assert.Equal(t, want.GetID(), tt.args.messages[i].GetID())
				assert.Equal(t, want.GetReference(), tt.args.messages[i].GetReference())
				assert.Equal(t, want.GetSourceURI(), tt.args.messages[i].GetSourceURI())
				assert.Equal(t, string(want.GetBody()), tt.args.messages[i].GetBody())
			}
		})
	}
}

func TestInputEvent_Messages(t *testing.T) {
	type fields struct {
		messages []Input
	}
	tests := []struct {
		name   string
		fields fields
		want   []Input
	}{
		{
			name: "success",
			fields: fields{
				messages: []Input{
					&Message{
						ID:        "123",
						Reference: "ref1",
						SourceURI: "src1",
						Body:      []byte(`{}`),
					}},
			},
			want: []Input{&Message{
				ID:        "123",
				Reference: "ref1",
				SourceURI: "src1",
				Body:      []byte(`{}`),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ie := &InputEvent{
				inputs: tt.fields.messages,
			}
			assert.Equalf(t, tt.want, ie.Inputs(), "Inputs()")
		})
	}
}
