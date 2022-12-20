// Code generated by MockGen. DO NOT EDIT.
// Source: rpc_module.go

// Package mock_modules is a generated GoMock package.
package mock_modules

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	modules "github.com/pokt-network/pocket/internal/shared/modules"
)

// MockRPCModule is a mock of RPCModule interface.
type MockRPCModule struct {
	ctrl     *gomock.Controller
	recorder *MockRPCModuleMockRecorder
}

// MockRPCModuleMockRecorder is the mock recorder for MockRPCModule.
type MockRPCModuleMockRecorder struct {
	mock *MockRPCModule
}

// NewMockRPCModule creates a new mock instance.
func NewMockRPCModule(ctrl *gomock.Controller) *MockRPCModule {
	mock := &MockRPCModule{ctrl: ctrl}
	mock.recorder = &MockRPCModuleMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRPCModule) EXPECT() *MockRPCModuleMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockRPCModule) Create(runtime modules.RuntimeMgr) (modules.Module, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", runtime)
	ret0, _ := ret[0].(modules.Module)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockRPCModuleMockRecorder) Create(runtime interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockRPCModule)(nil).Create), runtime)
}

// GetBus mocks base method.
func (m *MockRPCModule) GetBus() modules.Bus {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBus")
	ret0, _ := ret[0].(modules.Bus)
	return ret0
}

// GetBus indicates an expected call of GetBus.
func (mr *MockRPCModuleMockRecorder) GetBus() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBus", reflect.TypeOf((*MockRPCModule)(nil).GetBus))
}

// GetModuleName mocks base method.
func (m *MockRPCModule) GetModuleName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetModuleName")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetModuleName indicates an expected call of GetModuleName.
func (mr *MockRPCModuleMockRecorder) GetModuleName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetModuleName", reflect.TypeOf((*MockRPCModule)(nil).GetModuleName))
}

// SetBus mocks base method.
func (m *MockRPCModule) SetBus(arg0 modules.Bus) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetBus", arg0)
}

// SetBus indicates an expected call of SetBus.
func (mr *MockRPCModuleMockRecorder) SetBus(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetBus", reflect.TypeOf((*MockRPCModule)(nil).SetBus), arg0)
}

// Start mocks base method.
func (m *MockRPCModule) Start() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Start")
	ret0, _ := ret[0].(error)
	return ret0
}

// Start indicates an expected call of Start.
func (mr *MockRPCModuleMockRecorder) Start() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockRPCModule)(nil).Start))
}

// Stop mocks base method.
func (m *MockRPCModule) Stop() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stop")
	ret0, _ := ret[0].(error)
	return ret0
}

// Stop indicates an expected call of Stop.
func (mr *MockRPCModuleMockRecorder) Stop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockRPCModule)(nil).Stop))
}

// ValidateConfig mocks base method.
func (m *MockRPCModule) ValidateConfig(arg0 modules.Config) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateConfig", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// ValidateConfig indicates an expected call of ValidateConfig.
func (mr *MockRPCModuleMockRecorder) ValidateConfig(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateConfig", reflect.TypeOf((*MockRPCModule)(nil).ValidateConfig), arg0)
}
