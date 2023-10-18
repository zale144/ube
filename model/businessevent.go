package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/zale144/ube/libs/validate"
)

type BusinessEvent struct {
	ID                    string     `json:"id"`
	DocType               string     `json:"doctype,omitempty"`
	Metadata              *Metadata  `json:"metadata,omitempty"`
	Event                 *Event     `json:"event,omitempty"`
	RawDataEvent          [][]byte   `json:"raw_data_event,omitempty"`
	Body                  []byte     `json:"-"`
	Pt                    *time.Time `json:"pt,omitempty"`
	BaseWarehouse         string     `json:"base_warehouse"`
	entity                Entity
	Entities              []Entity      `json:"entities,omitempty"`
	Error                 error         `json:"-"`
	PreviousActionMandate ActionMandate `json:"-"`
	PreviousAction        int           `json:"-"`
	RepublishAttempt      *int          `json:"is_republish,omitempty"`
}

var (
	_ PipelineMedium    = (*BusinessEvent)(nil)
	_ InputActionMedium = (*BusinessEvent)(nil)
)

type ActionMandate int

const (
	StopFurtherProcessing ActionMandate = 1 << iota
	ProcessOnlyCriticalActions
	// StopAndPark
	LogFailureAndContinue
	StopAndRaiseError
	StopAndRetry
)

// UnmarshalJSON overrides the default method
func (be *BusinessEvent) UnmarshalJSON(data []byte) error {
	if be.entity == nil {
		return errors.New("'entity' type must be provided")
	}

	m := make(map[string]interface{})
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	event, ok := m["event"].(map[string]interface{})
	if !ok {
		return errors.New("'event' is missing")
	}

	eventCat, ok := event["event_category"].(string)
	if !ok {
		return errors.New("'event_category' is missing")
	}

	ent, ok := m[eventCat]
	if !ok {
		return fmt.Errorf("entity '%s' is missing", eventCat)
	}

	type alias BusinessEvent
	aux := &alias{}

	if ents, ok := ent.([]interface{}); ok {
		m["entities"] = ents
		for i := 0; i < len(ents); i++ {
			aux.Entities = append(aux.Entities, reflect.New(reflect.TypeOf(be.entity).Elem()).Interface().(Entity)) // TODO
		}
	}

	delete(m, eventCat)

	jsn, _ := json.Marshal(m)

	if err := json.Unmarshal(jsn, aux); err != nil {
		return err
	}

	*be = BusinessEvent(*aux)

	if pa, ok := m["previous_action"]; ok {
		be.PreviousAction = int(pa.(float64))
	}

	return nil
}

func (be *BusinessEvent) MarshalJSON() ([]byte, error) {
	value := reflect.ValueOf(*be)
	t := value.Type()
	sf := make([]reflect.StructField, 0)
	for i := 0; i < t.NumField(); i++ {
		sf = append(sf, t.Field(i))
		if len(be.Entities) > 0 && t.Field(i).Name == "Entities" {
			sf[i].Tag = reflect.StructTag(fmt.Sprintf(`json:%q`, be.Event.EventCategory))
		}
	}

	newType := reflect.StructOf(sf)
	newValue := value.Convert(newType)

	return json.Marshal(newValue.Interface())
}

// InputsToBusinessEvents converts the inputs to []*BusinessEvent
func InputsToBusinessEvents(ins []Input, entity Entity) ([]Medium, error) {
	if entity == nil {
		return nil, fmt.Errorf("entity is not provided")
	}

	entType := reflect.TypeOf(entity)
	if entType.Kind() == reflect.Ptr {
		entType = entType.Elem()
	}

	bes := make([]Medium, len(ins))
	for i, in := range ins {
		be, ok := isBusinessEvent([]byte(in.GetBody()), entType)
		if !ok {
			be.initBase(in)
		}

		if be.entity == nil {
			be.entity = entity
		}

		if len(be.Entities) == 0 {
			be.Entities = append(be.Entities, entity)
		}

		bes[i] = be
	}

	return bes, nil
}

func isBusinessEvent(body []byte, entType reflect.Type) (*BusinessEvent, bool) {
	be := &BusinessEvent{
		entity: reflect.New(entType).Interface().(Entity),
	}
	if err := json.Unmarshal(body, be); err != nil {
		be.entity = nil
		return be, false
	}
	return be, true
}

// InitEntity converts the input to *BusinessEvent
func (be *BusinessEvent) InitEntity(entType reflect.Type) error {
	be.entity = reflect.New(entType).Interface().(Entity)

	if err := json.Unmarshal(be.Body, &be.entity); err != nil {
		zap.L().Warn("input is not a single record, retrying to unmarshal as a list: %s", zap.Error(err))

		list := make([]interface{}, 0)
		if err = json.Unmarshal(be.Body, &list); err != nil {
			return fmt.Errorf("unmarshal input body fail: %w", err)
		}

		for i := 0; i < len(list); i++ {
			be.Entities = append(be.Entities, reflect.New(entType).Interface().(Entity))
		}

		if err = json.Unmarshal(be.Body, &be.Entities); err != nil {
			return fmt.Errorf("unmarshal input body fail: %w", err)
		}
	} else {
		be.Entities = []Entity{be.entity}
	}

	for _, ent := range be.Entities {
		if err := validate.Struct(ent); err != nil {
			return fmt.Errorf("invalid input event: %w", err)
		}
	}

	return nil
}

// Override it for testing
var (
	Now     = time.Now
	UUIDStr = uuid.NewString
)

// initBase initialises the base part of the business event structure
func (be *BusinessEvent) initBase(in Input) {
	beID := UUIDStr()
	t := Now()
	nowStr := t.UTC().Format(time.RFC3339Nano)

	be.ID = beID

	if be.Event == nil {
		be.Event = &Event{}
	}

	be.Event.ID = in.GetID()
	be.Event.EventSource = in.GetSourceURI()
	be.Event.EventOccurredTime = nowStr
	be.Event.EventReceivedTime = nowStr
	be.Event.Reference = in.GetReference()

	if be.Metadata == nil {
		be.Metadata = &Metadata{}
	}

	be.Metadata.Created = nowStr
	be.Metadata.CreatedEventID = beID
	be.Metadata.LastUpdated = nowStr
	be.Metadata.LastUpdateEventOccurred = nowStr
	be.Metadata.LastUpdateEventID = beID
	be.Pt = &t
	be.Body = []byte(in.GetBody())
	be.RawDataEvent = [][]byte{[]byte(in.GetBody())}
}

// GetID gets the ID
func (be *BusinessEvent) GetID() string {
	return be.ID
}

// GetEvent gets the event if exists
func (be *BusinessEvent) GetEvent() *Event {
	if be.Event == nil {
		zap.L().Error("no event")
		return &Event{}
	}

	return be.Event
}

// GetRawData returns raw data
func (be *BusinessEvent) GetRawData() [][]byte {
	return be.RawDataEvent
}

// UpdateMetadata updates metadata
func (be *BusinessEvent) UpdateMetadata(now func() time.Time) {
	if be.Event.Metadata == nil {
		be.Event.Metadata = &Metadata{}
	}

	be.Event.Metadata.LastUpdated = now().Format(time.RFC3339Nano)
	be.Event.Metadata.LastUpdateEventID = be.GetID()
	be.Event.Metadata.LastUpdateEventOccurred = be.GetEvent().EventOccurredTime
}

func (be *BusinessEvent) GetEntities() []Entity {
	return be.Entities
}

func (be *BusinessEvent) SetEntities(entities []Entity) {
	be.Entities = entities
}

func (be *BusinessEvent) GetEventName() string {
	if be.Event == nil {
		return ""
	}
	return be.Event.EventName
}

func (be *BusinessEvent) GetEventCategory() string {
	if be.Event == nil {
		return ""
	}
	return be.Event.EventCategory
}

func (be *BusinessEvent) GetError() error {
	return be.Error
}

func (be *BusinessEvent) SetError(err error) {
	be.Error = err
}

func (be *BusinessEvent) GetPreviousAction() int {
	return be.PreviousAction
}

func (be *BusinessEvent) SetPreviousAction(prevAct int) {
	be.PreviousAction = prevAct
}

func (be *BusinessEvent) SetEventProcessedTime(t time.Time) {
	if be.Event == nil {
		return
	}
	be.Event.EventProcessedTime = t.String()
}

func (be *BusinessEvent) GetPreviousActionMandate() ActionMandate {
	return be.PreviousActionMandate
}

func (be *BusinessEvent) GetRepublishAttempt() *int {
	return be.RepublishAttempt
}

func (be *BusinessEvent) SetRepublishAttempt(republishAttempt *int) {
	be.RepublishAttempt = republishAttempt
}

func (be *BusinessEvent) IncrementRepublishAttempt() {
	if be.RepublishAttempt == nil {
		be.RepublishAttempt = new(int)
	}
	*be.RepublishAttempt++
}

func (be *BusinessEvent) SetPreviousActionMandate(mandate ActionMandate) {
	be.PreviousActionMandate = mandate
}

func (be *BusinessEvent) GetEventID() string {
	if be.Event == nil {
		return ""
	}
	return be.Event.ID
}

func (be *BusinessEvent) GetEventReference() string {
	if be.Event == nil {
		return ""
	}
	return be.Event.Reference
}

func (be *BusinessEvent) SetEventID(id string) {
	if be.Event == nil {
		be.Event = &Event{}
	}
	be.Event.ID = id
}

func (be *BusinessEvent) SetEventReference(ref string) {
	if be.Event == nil {
		be.Event = &Event{}
	}
	be.Event.Reference = ref
}

func (be *BusinessEvent) GetBody() []byte {
	return be.Body
}

func (be *BusinessEvent) SetBody(body []byte) {
	be.Body = body
}

func (be *BusinessEvent) SetEventName(evName string) {
	if be.Event == nil {
		be.Event = &Event{}
	}
	be.Event.EventName = evName
}

func (be *BusinessEvent) SetEventCategory(evCat string) {
	if be.Event == nil {
		be.Event = &Event{}
	}
	be.Event.EventCategory = evCat
}

func (be *BusinessEvent) SetSource(src string) {
	be.BaseWarehouse = src // TODO ?
}

func (be *BusinessEvent) GetEventSource() string {
	if be.Event == nil {
		return ""
	}
	return be.Event.EventSource
}
