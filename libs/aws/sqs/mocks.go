// Code generated by MockGen. DO NOT EDIT.
// Source: interfaces.go

// Package sqs is a generated GoMock package.
package sqs

import (
	aws "github.com/aws/aws-sdk-go/aws"
	request "github.com/aws/aws-sdk-go/aws/request"
	sqs "github.com/aws/aws-sdk-go/service/sqs"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// Mockclient is a mock of client interface
type Mockclient struct {
	ctrl     *gomock.Controller
	recorder *MockclientMockRecorder
}

// MockclientMockRecorder is the mock recorder for Mockclient
type MockclientMockRecorder struct {
	mock *Mockclient
}

// NewMockclient creates a new mock instance
func NewMockclient(ctrl *gomock.Controller) *Mockclient {
	mock := &Mockclient{ctrl: ctrl}
	mock.recorder = &MockclientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *Mockclient) EXPECT() *MockclientMockRecorder {
	return m.recorder
}

// SendMessageBatchWithContext mocks base method
func (m *Mockclient) SendMessageBatchWithContext(ctx aws.Context, input *sqs.SendMessageBatchInput, opts ...request.Option) (*sqs.SendMessageBatchOutput, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, input}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SendMessageBatchWithContext", varargs...)
	ret0, _ := ret[0].(*sqs.SendMessageBatchOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendMessageBatchWithContext indicates an expected call of SendMessageBatchWithContext
func (mr *MockclientMockRecorder) SendMessageBatchWithContext(ctx, input interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, input}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMessageBatchWithContext", reflect.TypeOf((*Mockclient)(nil).SendMessageBatchWithContext), varargs...)
}

// CreateQueueWithContext mocks base method
func (m *Mockclient) CreateQueueWithContext(ctx aws.Context, input *sqs.CreateQueueInput, opts ...request.Option) (*sqs.CreateQueueOutput, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, input}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CreateQueueWithContext", varargs...)
	ret0, _ := ret[0].(*sqs.CreateQueueOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateQueueWithContext indicates an expected call of CreateQueueWithContext
func (mr *MockclientMockRecorder) CreateQueueWithContext(ctx, input interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, input}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateQueueWithContext", reflect.TypeOf((*Mockclient)(nil).CreateQueueWithContext), varargs...)
}

// DeleteMessageBatchWithContext mocks base method
func (m *Mockclient) DeleteMessageBatchWithContext(ctx aws.Context, input *sqs.DeleteMessageBatchInput, opts ...request.Option) (*sqs.DeleteMessageBatchOutput, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, input}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteMessageBatchWithContext", varargs...)
	ret0, _ := ret[0].(*sqs.DeleteMessageBatchOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteMessageBatchWithContext indicates an expected call of DeleteMessageBatchWithContext
func (mr *MockclientMockRecorder) DeleteMessageBatchWithContext(ctx, input interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, input}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteMessageBatchWithContext", reflect.TypeOf((*Mockclient)(nil).DeleteMessageBatchWithContext), varargs...)
}

// ReceiveMessageWithContext mocks base method
func (m *Mockclient) ReceiveMessageWithContext(ctx aws.Context, input *sqs.ReceiveMessageInput, opts ...request.Option) (*sqs.ReceiveMessageOutput, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, input}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ReceiveMessageWithContext", varargs...)
	ret0, _ := ret[0].(*sqs.ReceiveMessageOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReceiveMessageWithContext indicates an expected call of ReceiveMessageWithContext
func (mr *MockclientMockRecorder) ReceiveMessageWithContext(ctx, input interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, input}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReceiveMessageWithContext", reflect.TypeOf((*Mockclient)(nil).ReceiveMessageWithContext), varargs...)
}

// PurgeQueueWithContext mocks base method
func (m *Mockclient) PurgeQueueWithContext(ctx aws.Context, input *sqs.PurgeQueueInput, opts ...request.Option) (*sqs.PurgeQueueOutput, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, input}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "PurgeQueueWithContext", varargs...)
	ret0, _ := ret[0].(*sqs.PurgeQueueOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PurgeQueueWithContext indicates an expected call of PurgeQueueWithContext
func (mr *MockclientMockRecorder) PurgeQueueWithContext(ctx, input interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, input}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PurgeQueueWithContext", reflect.TypeOf((*Mockclient)(nil).PurgeQueueWithContext), varargs...)
}