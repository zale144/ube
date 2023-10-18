// Code generated by MockGen. DO NOT EDIT.
// Source: ./interfaces.go

// Package actions is a generated GoMock package.
package actions

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	model "github.com/zale144/ube/model"
	io "io"
	reflect "reflect"
)

// MockIAcker is a mock of IAcker interface
type MockIAcker struct {
	ctrl     *gomock.Controller
	recorder *MockIAckerMockRecorder
}

// MockIAckerMockRecorder is the mock recorder for MockIAcker
type MockIAckerMockRecorder struct {
	mock *MockIAcker
}

// NewMockIAcker creates a new mock instance
func NewMockIAcker(ctrl *gomock.Controller) *MockIAcker {
	mock := &MockIAcker{ctrl: ctrl}
	mock.recorder = &MockIAckerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockIAcker) EXPECT() *MockIAckerMockRecorder {
	return m.recorder
}

// AckMessages mocks base method
func (m *MockIAcker) AckMessages(ctx context.Context, msg ...model.Input) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range msg {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "AckMessages", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// AckMessages indicates an expected call of AckMessages
func (mr *MockIAckerMockRecorder) AckMessages(ctx interface{}, msg ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, msg...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AckMessages", reflect.TypeOf((*MockIAcker)(nil).AckMessages), varargs...)
}

// MockIPublisher is a mock of IPublisher interface
type MockIPublisher struct {
	ctrl     *gomock.Controller
	recorder *MockIPublisherMockRecorder
}

// MockIPublisherMockRecorder is the mock recorder for MockIPublisher
type MockIPublisherMockRecorder struct {
	mock *MockIPublisher
}

// NewMockIPublisher creates a new mock instance
func NewMockIPublisher(ctrl *gomock.Controller) *MockIPublisher {
	mock := &MockIPublisher{ctrl: ctrl}
	mock.recorder = &MockIPublisherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockIPublisher) EXPECT() *MockIPublisherMockRecorder {
	return m.recorder
}

// PublishEvents mocks base method
func (m *MockIPublisher) PublishEvents(ctx context.Context, msg ...model.Input) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range msg {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "PublishEvents", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// PublishEvents indicates an expected call of PublishEvents
func (mr *MockIPublisherMockRecorder) PublishEvents(ctx interface{}, msg ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, msg...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PublishEvents", reflect.TypeOf((*MockIPublisher)(nil).PublishEvents), varargs...)
}

// MockIRepublisher is a mock of IRepublisher interface
type MockIRepublisher struct {
	ctrl     *gomock.Controller
	recorder *MockIRepublisherMockRecorder
}

// MockIRepublisherMockRecorder is the mock recorder for MockIRepublisher
type MockIRepublisherMockRecorder struct {
	mock *MockIRepublisher
}

// NewMockIRepublisher creates a new mock instance
func NewMockIRepublisher(ctrl *gomock.Controller) *MockIRepublisher {
	mock := &MockIRepublisher{ctrl: ctrl}
	mock.recorder = &MockIRepublisherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockIRepublisher) EXPECT() *MockIRepublisherMockRecorder {
	return m.recorder
}

// AckMessages mocks base method
func (m *MockIRepublisher) AckMessages(ctx context.Context, msg ...model.Input) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range msg {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "AckMessages", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// AckMessages indicates an expected call of AckMessages
func (mr *MockIRepublisherMockRecorder) AckMessages(ctx interface{}, msg ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, msg...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AckMessages", reflect.TypeOf((*MockIRepublisher)(nil).AckMessages), varargs...)
}

// PublishEvents mocks base method
func (m *MockIRepublisher) PublishEvents(ctx context.Context, msg ...model.Input) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range msg {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "PublishEvents", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// PublishEvents indicates an expected call of PublishEvents
func (mr *MockIRepublisherMockRecorder) PublishEvents(ctx interface{}, msg ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, msg...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PublishEvents", reflect.TypeOf((*MockIRepublisher)(nil).PublishEvents), varargs...)
}

// MockIUploader is a mock of IUploader interface
type MockIUploader struct {
	ctrl     *gomock.Controller
	recorder *MockIUploaderMockRecorder
}

// MockIUploaderMockRecorder is the mock recorder for MockIUploader
type MockIUploaderMockRecorder struct {
	mock *MockIUploader
}

// NewMockIUploader creates a new mock instance
func NewMockIUploader(ctrl *gomock.Controller) *MockIUploader {
	mock := &MockIUploader{ctrl: ctrl}
	mock.recorder = &MockIUploaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockIUploader) EXPECT() *MockIUploaderMockRecorder {
	return m.recorder
}

// UploadFile mocks base method
func (m *MockIUploader) UploadFile(ctx context.Context, key string, body io.Reader) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadFile", ctx, key, body)
	ret0, _ := ret[0].(error)
	return ret0
}

// UploadFile indicates an expected call of UploadFile
func (mr *MockIUploaderMockRecorder) UploadFile(ctx, key, body interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadFile", reflect.TypeOf((*MockIUploader)(nil).UploadFile), ctx, key, body)
}

// MockIDownloader is a mock of IDownloader interface
type MockIDownloader struct {
	ctrl     *gomock.Controller
	recorder *MockIDownloaderMockRecorder
}

// MockIDownloaderMockRecorder is the mock recorder for MockIDownloader
type MockIDownloaderMockRecorder struct {
	mock *MockIDownloader
}

// NewMockIDownloader creates a new mock instance
func NewMockIDownloader(ctrl *gomock.Controller) *MockIDownloader {
	mock := &MockIDownloader{ctrl: ctrl}
	mock.recorder = &MockIDownloaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockIDownloader) EXPECT() *MockIDownloaderMockRecorder {
	return m.recorder
}

// DownloadFile mocks base method
func (m *MockIDownloader) DownloadFile(ctx context.Context, key string, body io.Writer) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DownloadFile", ctx, key, body)
	ret0, _ := ret[0].(error)
	return ret0
}

// DownloadFile indicates an expected call of DownloadFile
func (mr *MockIDownloaderMockRecorder) DownloadFile(ctx, key, body interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DownloadFile", reflect.TypeOf((*MockIDownloader)(nil).DownloadFile), ctx, key, body)
}

// DownloadFileFromBucket mocks base method
func (m *MockIDownloader) DownloadFileFromBucket(ctx context.Context, bucket, key string, body io.Writer) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DownloadFileFromBucket", ctx, bucket, key, body)
	ret0, _ := ret[0].(error)
	return ret0
}

// DownloadFileFromBucket indicates an expected call of DownloadFileFromBucket
func (mr *MockIDownloaderMockRecorder) DownloadFileFromBucket(ctx, bucket, key, body interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DownloadFileFromBucket", reflect.TypeOf((*MockIDownloader)(nil).DownloadFileFromBucket), ctx, bucket, key, body)
}

// MockIService is a mock of IService interface
type MockIService struct {
	ctrl     *gomock.Controller
	recorder *MockIServiceMockRecorder
}

// MockIServiceMockRecorder is the mock recorder for MockIService
type MockIServiceMockRecorder struct {
	mock *MockIService
}

// NewMockIService creates a new mock instance
func NewMockIService(ctrl *gomock.Controller) *MockIService {
	mock := &MockIService{ctrl: ctrl}
	mock.recorder = &MockIServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockIService) EXPECT() *MockIServiceMockRecorder {
	return m.recorder
}

// Execute mocks base method
func (m *MockIService) Execute(ctx context.Context, bes ...model.Medium) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range bes {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Execute", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Execute indicates an expected call of Execute
func (mr *MockIServiceMockRecorder) Execute(ctx interface{}, bes ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, bes...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Execute", reflect.TypeOf((*MockIService)(nil).Execute), varargs...)
}

// MockIRepository is a mock of IRepository interface
type MockIRepository struct {
	ctrl     *gomock.Controller
	recorder *MockIRepositoryMockRecorder
}

// MockIRepositoryMockRecorder is the mock recorder for MockIRepository
type MockIRepositoryMockRecorder struct {
	mock *MockIRepository
}

// NewMockIRepository creates a new mock instance
func NewMockIRepository(ctrl *gomock.Controller) *MockIRepository {
	mock := &MockIRepository{ctrl: ctrl}
	mock.recorder = &MockIRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockIRepository) EXPECT() *MockIRepositoryMockRecorder {
	return m.recorder
}

// GetEntity mocks base method
func (m *MockIRepository) GetEntity(arg0 context.Context, arg1 model.Key, arg2 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEntity", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetEntity indicates an expected call of GetEntity
func (mr *MockIRepositoryMockRecorder) GetEntity(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEntity", reflect.TypeOf((*MockIRepository)(nil).GetEntity), arg0, arg1, arg2)
}

// EntityExists mocks base method
func (m *MockIRepository) EntityExists(arg0 context.Context, arg1 model.Key) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EntityExists", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EntityExists indicates an expected call of EntityExists
func (mr *MockIRepositoryMockRecorder) EntityExists(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EntityExists", reflect.TypeOf((*MockIRepository)(nil).EntityExists), arg0, arg1)
}

// SaveEntities mocks base method
func (m *MockIRepository) SaveEntities(ctx context.Context, entity ...model.Entity) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range entity {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SaveEntities", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveEntities indicates an expected call of SaveEntities
func (mr *MockIRepositoryMockRecorder) SaveEntities(ctx interface{}, entity ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, entity...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveEntities", reflect.TypeOf((*MockIRepository)(nil).SaveEntities), varargs...)
}
