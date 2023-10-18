package actions

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"go.uber.org/zap"

	"github.com/zale144/ube/model"
)

// Publish is a wrapper for publishing the business event
type Publish struct {
	publisher IPublisher
	Base
}

// MarshalIndent Override it for testing
var (
	MarshalIndent func(v interface{}, prefix, indent string) ([]byte, error)
)

// Publisher constructs a new Publish
func Publisher(publisher IPublisher, options ...BaseOption) *Publish {
	pub := &Publish{
		publisher: publisher,
		Base: Base{
			batchSize:      100,
			failureMandate: model.StopFurtherProcessing,
		},
	}

	for _, opt := range options {
		opt(&pub.Base)
	}

	return pub
}

// Process implements the action interface in UBE, executes the underlying embedded device
func (e Publish) Process(ctx context.Context, bes ...model.Medium) {
	var counter int
	msgs := make([]model.Input, len(bes))

	for i := range bes {
		be := bes[i]
		if be.GetID() == "" {
			bes[i].SetError(errors.New("publish can't handle empty business event"))
			continue
		}

		if _, ok := e.skips[be.GetEventName()]; ok {
			continue
		}

		var (
			msg *model.Message
			err error
		)
		msg, err = toMessage(be, false)
		if err != nil {
			be.SetError(fmt.Errorf("convert business event %s to message fail: %w", be.GetID(), err))
			continue
		}

		msgs[i] = msg
		counter++
	}

	if counter > 0 {
		var err error
		if err = e.publisher.PublishEvents(ctx, msgs...); err != nil {
			err = fmt.Errorf("publish message fail: %w", err)
		}
		for i := range bes {
			bes[i].SetError(err)
		}
	}

	zap.L().Info("messages published", zap.Int("size", counter))
}

// asMsg converts the business event into a message taking a bool to
// determine whether to compress the message
func toMessage(be model.Medium, compressed bool) (*model.Message, error) {
	var (
		err  error
		body []byte
	)
	if MarshalIndent != nil {
		body, err = MarshalIndent(be, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("marshal business event fail: %w", err)
		}
	} else {
		body, err = json.Marshal(be)
		if err != nil {
			return nil, fmt.Errorf("marshal business event fail: %w", err)
		}
	}

	if compressed {
		var buf bytes.Buffer
		// TODO:	err = compression.GzipData(&buf, dta)
		if err != nil {
			return nil, fmt.Errorf("compress business event fail: %w", err)
		}

		body = buf.Bytes()
	}

	return &model.Message{ID: be.GetID(), Body: body}, nil
}

func (Publish) Name() string {
	return "Publisher"
}

func (e Publish) DepCallNames() []string {
	return []string{"PublishEvents"}
}
