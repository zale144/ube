package model

import (
	"encoding/json"
)

// Message is the event based communication medium
type Message struct {
	ID        string          `json:"id"`
	Reference string          `json:"reference"`
	SourceURI string          `json:"source_uri"`
	Body      json.RawMessage `json:"body"`
}

func (m *Message) MarshalJSON() ([]byte, error) {
	type tempMsg struct {
		ID        string          `json:"id,omitempty"`
		Reference string          `json:"reference,omitempty"`
		Body      json.RawMessage `json:"body,omitempty"`
		SourceURI string          `json:"source_uri,omitempty"`
	}
	tm := &tempMsg{
		ID:        m.ID,
		Reference: m.Reference,
		Body:      []byte(m.Body),
		SourceURI: m.SourceURI,
	}
	return json.MarshalIndent(tm, "", "	")
}

// GetID returns the message ID (implements model.Input)
func (m Message) GetID() string {
	return m.ID
}

// GetReference returns the message reference (implements model.Input)
func (m Message) GetReference() string {
	return m.Reference
}

// GetSourceURI returns the message source URI (implements model.Input)
func (m Message) GetSourceURI() string {
	return m.SourceURI
}

// GetBody returns the message body (implements model.Input)
func (m Message) GetBody() string {
	return string(m.Body)
}
