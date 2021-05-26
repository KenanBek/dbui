// Code generated by MockGen. DO NOT EDIT.
// Source: dbui.go

// Package controller is a generated GoMock package.
package controller

import (
	internal "dbui/internal"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockAppConfig is a mock of AppConfig interface
type MockAppConfig struct {
	ctrl     *gomock.Controller
	recorder *MockAppConfigMockRecorder
}

// MockAppConfigMockRecorder is the mock recorder for MockAppConfig
type MockAppConfigMockRecorder struct {
	mock *MockAppConfig
}

// NewMockAppConfig creates a new mock instance
func NewMockAppConfig(ctrl *gomock.Controller) *MockAppConfig {
	mock := &MockAppConfig{ctrl: ctrl}
	mock.recorder = &MockAppConfigMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAppConfig) EXPECT() *MockAppConfigMockRecorder {
	return m.recorder
}

// DataSourceConfigs mocks base method
func (m *MockAppConfig) DataSourceConfigs() map[string]internal.DataSourceConfig {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DataSourceConfigs")
	ret0, _ := ret[0].(map[string]internal.DataSourceConfig)
	return ret0
}

// DataSourceConfigs indicates an expected call of DataSourceConfigs
func (mr *MockAppConfigMockRecorder) DataSourceConfigs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DataSourceConfigs", reflect.TypeOf((*MockAppConfig)(nil).DataSourceConfigs))
}

// Default mocks base method
func (m *MockAppConfig) Default() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Default")
	ret0, _ := ret[0].(string)
	return ret0
}

// Default indicates an expected call of Default
func (mr *MockAppConfigMockRecorder) Default() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Default", reflect.TypeOf((*MockAppConfig)(nil).Default))
}

// MockDataSourceConfig is a mock of DataSourceConfig interface
type MockDataSourceConfig struct {
	ctrl     *gomock.Controller
	recorder *MockDataSourceConfigMockRecorder
}

// MockDataSourceConfigMockRecorder is the mock recorder for MockDataSourceConfig
type MockDataSourceConfigMockRecorder struct {
	mock *MockDataSourceConfig
}

// NewMockDataSourceConfig creates a new mock instance
func NewMockDataSourceConfig(ctrl *gomock.Controller) *MockDataSourceConfig {
	mock := &MockDataSourceConfig{ctrl: ctrl}
	mock.recorder = &MockDataSourceConfigMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDataSourceConfig) EXPECT() *MockDataSourceConfigMockRecorder {
	return m.recorder
}

// Alias mocks base method
func (m *MockDataSourceConfig) Alias() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Alias")
	ret0, _ := ret[0].(string)
	return ret0
}

// Alias indicates an expected call of Alias
func (mr *MockDataSourceConfigMockRecorder) Alias() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Alias", reflect.TypeOf((*MockDataSourceConfig)(nil).Alias))
}

// Type mocks base method
func (m *MockDataSourceConfig) Type() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Type")
	ret0, _ := ret[0].(string)
	return ret0
}

// Type indicates an expected call of Type
func (mr *MockDataSourceConfigMockRecorder) Type() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Type", reflect.TypeOf((*MockDataSourceConfig)(nil).Type))
}

// DSN mocks base method
func (m *MockDataSourceConfig) DSN() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DSN")
	ret0, _ := ret[0].(string)
	return ret0
}

// DSN indicates an expected call of DSN
func (mr *MockDataSourceConfigMockRecorder) DSN() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DSN", reflect.TypeOf((*MockDataSourceConfig)(nil).DSN))
}

// MockDataSource is a mock of DataSource interface
type MockDataSource struct {
	ctrl     *gomock.Controller
	recorder *MockDataSourceMockRecorder
}

// MockDataSourceMockRecorder is the mock recorder for MockDataSource
type MockDataSourceMockRecorder struct {
	mock *MockDataSource
}

// NewMockDataSource creates a new mock instance
func NewMockDataSource(ctrl *gomock.Controller) *MockDataSource {
	mock := &MockDataSource{ctrl: ctrl}
	mock.recorder = &MockDataSourceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDataSource) EXPECT() *MockDataSourceMockRecorder {
	return m.recorder
}

// Ping mocks base method
func (m *MockDataSource) Ping() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping")
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping
func (mr *MockDataSourceMockRecorder) Ping() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockDataSource)(nil).Ping))
}

// ListSchemas mocks base method
func (m *MockDataSource) ListSchemas() ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListSchemas")
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListSchemas indicates an expected call of ListSchemas
func (mr *MockDataSourceMockRecorder) ListSchemas() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListSchemas", reflect.TypeOf((*MockDataSource)(nil).ListSchemas))
}

// ListTables mocks base method
func (m *MockDataSource) ListTables(schema string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListTables", schema)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListTables indicates an expected call of ListTables
func (mr *MockDataSourceMockRecorder) ListTables(schema interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListTables", reflect.TypeOf((*MockDataSource)(nil).ListTables), schema)
}

// PreviewTable mocks base method
func (m *MockDataSource) PreviewTable(schema, table string) ([][]*string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PreviewTable", schema, table)
	ret0, _ := ret[0].([][]*string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PreviewTable indicates an expected call of PreviewTable
func (mr *MockDataSourceMockRecorder) PreviewTable(schema, table interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PreviewTable", reflect.TypeOf((*MockDataSource)(nil).PreviewTable), schema, table)
}

// DescribeTable mocks base method
func (m *MockDataSource) DescribeTable(schema, table string) ([][]*string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DescribeTable", schema, table)
	ret0, _ := ret[0].([][]*string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DescribeTable indicates an expected call of DescribeTable
func (mr *MockDataSourceMockRecorder) DescribeTable(schema, table interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DescribeTable", reflect.TypeOf((*MockDataSource)(nil).DescribeTable), schema, table)
}

// Query mocks base method
func (m *MockDataSource) Query(schema, query string) ([][]*string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Query", schema, query)
	ret0, _ := ret[0].([][]*string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Query indicates an expected call of Query
func (mr *MockDataSourceMockRecorder) Query(schema, query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Query", reflect.TypeOf((*MockDataSource)(nil).Query), schema, query)
}

// MockDataController is a mock of DataController interface
type MockDataController struct {
	ctrl     *gomock.Controller
	recorder *MockDataControllerMockRecorder
}

// MockDataControllerMockRecorder is the mock recorder for MockDataController
type MockDataControllerMockRecorder struct {
	mock *MockDataController
}

// NewMockDataController creates a new mock instance
func NewMockDataController(ctrl *gomock.Controller) *MockDataController {
	mock := &MockDataController{ctrl: ctrl}
	mock.recorder = &MockDataControllerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDataController) EXPECT() *MockDataControllerMockRecorder {
	return m.recorder
}

// List mocks base method
func (m *MockDataController) List() [][]string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List")
	ret0, _ := ret[0].([][]string)
	return ret0
}

// List indicates an expected call of List
func (mr *MockDataControllerMockRecorder) List() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockDataController)(nil).List))
}

// Switch mocks base method
func (m *MockDataController) Switch(alias string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Switch", alias)
	ret0, _ := ret[0].(error)
	return ret0
}

// Switch indicates an expected call of Switch
func (mr *MockDataControllerMockRecorder) Switch(alias interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Switch", reflect.TypeOf((*MockDataController)(nil).Switch), alias)
}

// Current mocks base method
func (m *MockDataController) Current() internal.DataSource {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Current")
	ret0, _ := ret[0].(internal.DataSource)
	return ret0
}

// Current indicates an expected call of Current
func (mr *MockDataControllerMockRecorder) Current() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Current", reflect.TypeOf((*MockDataController)(nil).Current))
}