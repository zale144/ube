package model

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBusinessEvent_UnmarshalJSON(t *testing.T) {
	timestamp := time.Date(2021, 11, 22, 3, 4, 5, 0, time.UTC)

	type args struct {
		data   []byte
		entity Entity
	}
	tests := []struct {
		name    string
		args    args
		want    *BusinessEvent
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success: list of entities",
			args: args{
				data:   []byte(`{"id":"1234","doctype":"doc1","metadata":{"last_updated":"2021-11-22 03:04:05 +0000 UTC","last_update_event_occurred":"2021-11-22 03:04:05 +0000 UTC","last_update_event_id":"2222","created":"2021-11-22 03:04:05 +0000 UTC","created_event_id":"1111"},"event":{"event_name":"CreateProduct","event_category":"Product","event_source":"src1","metadata":{},"id":"msg1","reference":"ref1","event_occurred_time":"2021-11-22 03:04:05 +0000 UTC","event_received_time":"2021-11-22 03:04:05 +0000 UTC"},"raw_data_event":["e30="],"pt":"2021-11-22T03:04:05Z","base_warehouse":"Zale144","Product":[{"sku":""}]}`),
				entity: &Product{},
			},
			want: &BusinessEvent{
				ID:      "1234",
				DocType: "doc1",
				Metadata: &Metadata{
					Created:                 timestamp.String(),
					CreatedEventID:          "1111",
					LastUpdated:             timestamp.String(),
					LastUpdateEventOccurred: timestamp.String(),
					LastUpdateEventID:       "2222",
				},
				Event: &Event{
					EventHeader: EventHeader{
						EventName:     "CreateProduct",
						EventCategory: "Product",
						EventSource:   "src1",
					},
					Metadata:          &Metadata{},
					ID:                "msg1",
					EventOccurredTime: timestamp.String(),
					EventReceivedTime: timestamp.String(),
					Reference:         "ref1",
				},
				RawDataEvent:          [][]byte{[]byte("{}")},
				Pt:                    &timestamp,
				BaseWarehouse:         "Zale144",
				Entities:              []Entity{&Product{}},
				Error:                 nil,
				PreviousActionMandate: 0,
				PreviousAction:        0,
				RepublishAttempt:      nil,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
		{
			name: "success: with previous_action",
			args: args{
				data:   []byte(`{"event": {"event_category": "product"}, "product": [{}], "previous_action": 1}`),
				entity: &Product{},
			},
			want: &BusinessEvent{
				Event: &Event{
					EventHeader: EventHeader{
						EventCategory: "product",
					},
				},
				Entities:       []Entity{&Product{}},
				PreviousAction: 1,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
		{
			name: "failure: unmarshal into map",
			args: args{
				data:   []byte(``),
				entity: &Product{},
			},
			want: &BusinessEvent{
				entity: &Product{},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err != nil && err.Error() == "unexpected end of JSON input"
			},
		},
		{
			name: "failure: no event",
			args: args{
				data:   []byte(`{}`),
				entity: &Product{},
			},
			want: &BusinessEvent{
				entity: &Product{},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err != nil && err.Error() == "'event' is missing"
			},
		},
		{
			name: "failure: no event_category",
			args: args{
				data:   []byte(`{"event": {}}`),
				entity: &Product{},
			},
			want: &BusinessEvent{
				entity: &Product{},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err != nil && err.Error() == "'event_category' is missing"
			},
		},
		{
			name: "failure: entity object is missing",
			args: args{
				data:   []byte(`{"event": {"event_category": "product"}}`),
				entity: &Product{},
			},
			want: &BusinessEvent{
				entity: &Product{},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err != nil && err.Error() == "entity 'product' is missing"
			},
		},
		{
			name: "failure: unmarshal into business event alias",
			args: args{
				data:   []byte(`{"id": 123, "event": {"event_category": "product"}, "product": {}}`),
				entity: &Product{},
			},
			want: &BusinessEvent{
				entity: &Product{},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err != nil && err.Error() == "json: cannot unmarshal number into Go struct field alias.id of type string"
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			be := &BusinessEvent{
				entity: tt.args.entity,
			}

			err := be.UnmarshalJSON(tt.args.data)
			if !tt.wantErr(t, err, fmt.Sprintf("UnmarshalJSON(%v)", tt.args.data)) {
				t.Errorf("unexpected error: %s", err)
				return
			}
			assert.Equal(t, tt.want, be)
		})
	}
}

func TestBusinessEvent_MarshalJSON(t *testing.T) {
	timestamp := time.Date(2021, 11, 22, 3, 4, 5, 0, time.UTC)

	type fields struct {
		be *BusinessEvent
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			fields: fields{
				be: &BusinessEvent{
					ID:      "1234",
					DocType: "doc1",
					Metadata: &Metadata{
						Created:                 timestamp.String(),
						CreatedEventID:          "1111",
						LastUpdated:             timestamp.String(),
						LastUpdateEventOccurred: timestamp.String(),
						LastUpdateEventID:       "2222",
					},
					Event: &Event{
						EventHeader: EventHeader{
							EventName:     "CreateProduct",
							EventCategory: "Product",
							EventSource:   "src1",
						},
						Metadata:          &Metadata{},
						ID:                "msg1",
						EventOccurredTime: timestamp.String(),
						EventReceivedTime: timestamp.String(),
						Reference:         "ref1",
					},
					RawDataEvent:          [][]byte{[]byte("{}")},
					Body:                  []byte("{}"),
					Pt:                    &timestamp,
					BaseWarehouse:         "Zale144",
					Entities:              []Entity{&Product{}},
					Error:                 nil,
					PreviousActionMandate: 0,
					PreviousAction:        0,
					RepublishAttempt:      nil,
				},
			},
			want: []byte(`{"id":"1234","doctype":"doc1","metadata":{"last_updated":"2021-11-22 03:04:05 +0000 UTC","last_update_event_occurred":"2021-11-22 03:04:05 +0000 UTC","last_update_event_id":"2222","created":"2021-11-22 03:04:05 +0000 UTC","created_event_id":"1111"},"event":{"event_name":"CreateProduct","event_category":"Product","event_source":"src1","metadata":{},"id":"msg1","reference":"ref1","event_occurred_time":"2021-11-22 03:04:05 +0000 UTC","event_received_time":"2021-11-22 03:04:05 +0000 UTC"},"raw_data_event":["e30="],"pt":"2021-11-22T03:04:05Z","base_warehouse":"Zale144","Product":[{}]}`),
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fields.be.MarshalJSON()
			if !tt.wantErr(t, err) {
				t.Errorf("MarshalJSON() error = %v", err)
				return
			}
			assert.JSONEqf(t, string(tt.want), string(got), "MarshalJSON()")
		})
	}
}

func TestInputsToBusinessEvents(t *testing.T) {
	UUIDStr = func() string {
		return "123"
	}
	Now = func() time.Time { return time.Date(2021, 11, 22, 3, 4, 5, 0, time.UTC) }
	now := Now()

	type args struct {
		ins    []Input
		entity Entity
	}
	tests := []struct {
		name    string
		args    args
		want    []Medium
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			args: args{
				ins: []Input{
					&Message{
						ID:        "123",
						Reference: "ref1",
						SourceURI: "src1",
						Body:      []byte(`{}`),
					},
				},
				entity: &Product{},
			},
			want: []Medium{
				&BusinessEvent{
					ID:      "123",
					DocType: "",
					Metadata: &Metadata{
						LastUpdated:             "2021-11-22T03:04:05Z",
						LastUpdateEventOccurred: "2021-11-22T03:04:05Z",
						LastUpdateEventID:       "123",
						Created:                 "2021-11-22T03:04:05Z",
						CreatedEventID:          "123",
					},
					Event: &Event{
						EventHeader: EventHeader{
							EventSource: "src1",
						},
						ID:                "123",
						Reference:         "ref1",
						EventOccurredTime: "2021-11-22T03:04:05Z",
						EventReceivedTime: "2021-11-22T03:04:05Z",
					},
					entity:       &Product{},
					RawDataEvent: [][]byte{[]byte("{}")},
					Body:         []byte("{}"),
					Pt:           &now,
				},
			},
			wantErr: func(_ assert.TestingT, err error, _ ...interface{}) bool {
				return err == nil
			},
		},
		{
			name: "success: business event",
			args: args{
				ins: []Input{
					&Message{
						ID:        "123",
						Reference: "ref1",
						SourceURI: "src1",
						Body:      []byte(`{"id":"123","metadata":{"last_updated":"2021-11-22T03:04:05Z","last_update_event_occurred":"2021-11-22T03:04:05Z","last_update_event_id":"123","created":"2021-11-22T03:04:05Z","created_event_id":"123"},"event":{"event_name":"CreateProduct","event_category":"product","event_source":"src1","id":"123","reference":"ref1","event_occurred_time":"2021-11-22T03:04:05Z","event_received_time":"2021-11-22T03:04:05Z"},"product":{}, "raw_data_event":["e30="],"pt":"2021-11-22T03:04:05Z","base_warehouse":""}`),
					},
				},
				entity: &Product{},
			},
			want: []Medium{
				&BusinessEvent{
					ID:      "123",
					DocType: "",
					Metadata: &Metadata{
						LastUpdated:             "2021-11-22T03:04:05Z",
						LastUpdateEventOccurred: "2021-11-22T03:04:05Z",
						LastUpdateEventID:       "123",
						Created:                 "2021-11-22T03:04:05Z",
						CreatedEventID:          "123",
					},
					Event: &Event{
						EventHeader: EventHeader{
							EventName:     "CreateProduct",
							EventCategory: "product",
							EventSource:   "src1",
						},
						ID:                "123",
						Reference:         "ref1",
						EventOccurredTime: "2021-11-22T03:04:05Z",
						EventReceivedTime: "2021-11-22T03:04:05Z",
					},
					entity:       &Product{},
					RawDataEvent: [][]byte{[]byte("{}")},
					Pt:           &now,
				},
			},
			wantErr: func(_ assert.TestingT, err error, _ ...interface{}) bool {
				return err == nil
			},
		},
		{
			name: "failure: no entity",
			args: args{
				ins: []Input{
					&Message{
						ID:        "123",
						Reference: "ref1",
						SourceURI: "src1",
						Body:      []byte(`{"id":"123","metadata":{"last_updated":"2021-11-22T03:04:05Z","last_update_event_occurred":"2021-11-22T03:04:05Z","last_update_event_id":"123","created":"2021-11-22T03:04:05Z","created_event_id":"123"},"event":{"event_name":"CreateProduct","event_category":"product","event_source":"src1","id":"123","reference":"ref1","event_occurred_time":"2021-11-22T03:04:05Z","event_received_time":"2021-11-22T03:04:05Z"},"product":{}, "raw_data_event":["e30="],"pt":"2021-11-22T03:04:05Z","base_warehouse":""}`),
					},
				},
				entity: nil,
			},
			want: nil,
			wantErr: func(_ assert.TestingT, err error, _ ...interface{}) bool {
				return err != nil && err.Error() == "entity is not provided"
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := InputsToBusinessEvents(tt.args.ins, tt.args.entity)
			if !tt.wantErr(t, err) {
				t.Errorf("InputsToBusinessEvents(): unexpected error: %s", err)
				return
			}
			assert.Equalf(t, tt.want, got, "InputsToBusinessEvents(%v, %v)", tt.args.ins, tt.args.entity)

			jsn, _ := json.Marshal(tt.want)
			t.Log(string(jsn))

			gotj, _ := json.Marshal(got)
			t.Log(string(gotj))
		})
	}
}

// Product is the base product model
type Product struct {
	PBaseKey
	Product          string `json:"product,omitempty"`
	CreatedAt        string `json:"created_at,omitempty"`
	Name             string `json:"name,omitempty"`
	Description      string `json:"description,omitempty"`
	ShortDescription string `json:"short_description,omitempty"`
}

type PBaseKey struct {
	ProductID int32 `json:"product_id,omitempty"`
}

func (b PBaseKey) PK() string {
	return fmt.Sprint(b.ProductID)
}

// GetKey returns the key to the entity
func (p Product) GetKey() Key {
	return p.PBaseKey
}
