package model

// InputEvent is the message event implementation of the model.Input interface
type InputEvent struct {
	inputs []Input
}

// NewInputEvent returns the pointer to the new InputEvent
func NewInputEvent(inputs []Input) *InputEvent {
	return &InputEvent{inputs: inputs}
}

// Inputs returns the input event inputs
func (ie *InputEvent) Inputs() []Input {
	return ie.inputs
}

// Input is the abstract pipeline input
type Input interface {
	GetID() string
	GetReference() string
	GetSourceURI() string
	GetBody() string
}
