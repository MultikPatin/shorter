// Code generated by MockGen. DO NOT EDIT.
// Source: main/internal/interfaces (interfaces: HealthRepository,LinksRepository,FileStorageProducer,FileStorageConsumer,UsersRepository)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	models "main/internal/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockHealthRepository is a mock of HealthRepository interface.
type MockHealthRepository struct {
	ctrl     *gomock.Controller
	recorder *MockHealthRepositoryMockRecorder
}

// MockHealthRepositoryMockRecorder is the mock recorder for MockHealthRepository.
type MockHealthRepositoryMockRecorder struct {
	mock *MockHealthRepository
}

// NewMockHealthRepository creates a new mock instance.
func NewMockHealthRepository(ctrl *gomock.Controller) *MockHealthRepository {
	mock := &MockHealthRepository{ctrl: ctrl}
	mock.recorder = &MockHealthRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHealthRepository) EXPECT() *MockHealthRepositoryMockRecorder {
	return m.recorder
}

// Ping mocks base method.
func (m *MockHealthRepository) Ping() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping")
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockHealthRepositoryMockRecorder) Ping() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockHealthRepository)(nil).Ping))
}

// MockLinksRepository is a mock of LinksRepository interface.
type MockLinksRepository struct {
	ctrl     *gomock.Controller
	recorder *MockLinksRepositoryMockRecorder
}

// MockLinksRepositoryMockRecorder is the mock recorder for MockLinksRepository.
type MockLinksRepositoryMockRecorder struct {
	mock *MockLinksRepository
}

// NewMockLinksRepository creates a new mock instance.
func NewMockLinksRepository(ctrl *gomock.Controller) *MockLinksRepository {
	mock := &MockLinksRepository{ctrl: ctrl}
	mock.recorder = &MockLinksRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLinksRepository) EXPECT() *MockLinksRepositoryMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockLinksRepository) Add(arg0 context.Context, arg1 models.AddedLink) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Add indicates an expected call of Add.
func (mr *MockLinksRepositoryMockRecorder) Add(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockLinksRepository)(nil).Add), arg0, arg1)
}

// AddBatch mocks base method.
func (m *MockLinksRepository) AddBatch(arg0 context.Context, arg1 []models.AddedLink) ([]models.Result, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddBatch", arg0, arg1)
	ret0, _ := ret[0].([]models.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddBatch indicates an expected call of AddBatch.
func (mr *MockLinksRepositoryMockRecorder) AddBatch(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddBatch", reflect.TypeOf((*MockLinksRepository)(nil).AddBatch), arg0, arg1)
}

// Get mocks base method.
func (m *MockLinksRepository) Get(arg0 context.Context, arg1 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockLinksRepositoryMockRecorder) Get(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockLinksRepository)(nil).Get), arg0, arg1)
}

// MockFileStorageProducer is a mock of FileStorageProducer interface.
type MockFileStorageProducer struct {
	ctrl     *gomock.Controller
	recorder *MockFileStorageProducerMockRecorder
}

// MockFileStorageProducerMockRecorder is the mock recorder for MockFileStorageProducer.
type MockFileStorageProducerMockRecorder struct {
	mock *MockFileStorageProducer
}

// NewMockFileStorageProducer creates a new mock instance.
func NewMockFileStorageProducer(ctrl *gomock.Controller) *MockFileStorageProducer {
	mock := &MockFileStorageProducer{ctrl: ctrl}
	mock.recorder = &MockFileStorageProducerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFileStorageProducer) EXPECT() *MockFileStorageProducerMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockFileStorageProducer) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockFileStorageProducerMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockFileStorageProducer)(nil).Close))
}

// WriteEvent mocks base method.
func (m *MockFileStorageProducer) WriteEvent(arg0 *models.Event) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteEvent", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteEvent indicates an expected call of WriteEvent.
func (mr *MockFileStorageProducerMockRecorder) WriteEvent(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteEvent", reflect.TypeOf((*MockFileStorageProducer)(nil).WriteEvent), arg0)
}

// MockFileStorageConsumer is a mock of FileStorageConsumer interface.
type MockFileStorageConsumer struct {
	ctrl     *gomock.Controller
	recorder *MockFileStorageConsumerMockRecorder
}

// MockFileStorageConsumerMockRecorder is the mock recorder for MockFileStorageConsumer.
type MockFileStorageConsumerMockRecorder struct {
	mock *MockFileStorageConsumer
}

// NewMockFileStorageConsumer creates a new mock instance.
func NewMockFileStorageConsumer(ctrl *gomock.Controller) *MockFileStorageConsumer {
	mock := &MockFileStorageConsumer{ctrl: ctrl}
	mock.recorder = &MockFileStorageConsumerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFileStorageConsumer) EXPECT() *MockFileStorageConsumerMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockFileStorageConsumer) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockFileStorageConsumerMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockFileStorageConsumer)(nil).Close))
}

// ReadAllEvents mocks base method.
func (m *MockFileStorageConsumer) ReadAllEvents() ([]*models.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadAllEvents")
	ret0, _ := ret[0].([]*models.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadAllEvents indicates an expected call of ReadAllEvents.
func (mr *MockFileStorageConsumerMockRecorder) ReadAllEvents() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadAllEvents", reflect.TypeOf((*MockFileStorageConsumer)(nil).ReadAllEvents))
}

// MockUsersRepository is a mock of UsersRepository interface.
type MockUsersRepository struct {
	ctrl     *gomock.Controller
	recorder *MockUsersRepositoryMockRecorder
}

// MockUsersRepositoryMockRecorder is the mock recorder for MockUsersRepository.
type MockUsersRepositoryMockRecorder struct {
	mock *MockUsersRepository
}

// NewMockUsersRepository creates a new mock instance.
func NewMockUsersRepository(ctrl *gomock.Controller) *MockUsersRepository {
	mock := &MockUsersRepository{ctrl: ctrl}
	mock.recorder = &MockUsersRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUsersRepository) EXPECT() *MockUsersRepositoryMockRecorder {
	return m.recorder
}

// DeleteLinks mocks base method.
func (m *MockUsersRepository) DeleteLinks(arg0 context.Context, arg1 []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteLinks", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteLinks indicates an expected call of DeleteLinks.
func (mr *MockUsersRepositoryMockRecorder) DeleteLinks(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteLinks", reflect.TypeOf((*MockUsersRepository)(nil).DeleteLinks), arg0, arg1)
}

// GetLinks mocks base method.
func (m *MockUsersRepository) GetLinks(arg0 context.Context) ([]models.UserLinks, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLinks", arg0)
	ret0, _ := ret[0].([]models.UserLinks)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLinks indicates an expected call of GetLinks.
func (mr *MockUsersRepositoryMockRecorder) GetLinks(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLinks", reflect.TypeOf((*MockUsersRepository)(nil).GetLinks), arg0)
}

// Login mocks base method.
func (m *MockUsersRepository) Login(arg0 context.Context) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", arg0)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Login indicates an expected call of Login.
func (mr *MockUsersRepositoryMockRecorder) Login(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockUsersRepository)(nil).Login), arg0)
}
