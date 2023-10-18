package model

import (
	"reflect"
	"time"
)

type Medium interface {
	GetID() string
	GetEventName() string
	GetRawData() [][]byte
	GetEntities() []Entity
	SetEntities([]Entity)

	// maybe do without these?
	GetError() error // maybe handle errors differently?
	SetError(error)  // maybe handle errors differently?
}

type InputActionMedium interface {
	Medium
	GetBody() []byte
	SetBody([]byte)
	SetEventName(string)
	SetEventCategory(string)
	SetSource(string)
	GetEventSource() string
	InitEntity(reflect.Type) error
}

type PipelineMedium interface {
	Medium
	GetPreviousAction() int
	SetPreviousAction(int)
	SetEventProcessedTime(time.Time)
	GetPreviousActionMandate() ActionMandate
	GetRepublishAttempt() *int
	SetRepublishAttempt(*int)
	IncrementRepublishAttempt()
	SetPreviousActionMandate(ActionMandate)
	GetEventID() string
	GetEventReference() string
	SetEventID(string)
	SetEventReference(string)
}
