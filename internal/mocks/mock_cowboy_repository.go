// Code generated by MockGen. DO NOT EDIT.
// Source: cowboy-app/internal/domain (interfaces: CowboyRepository)

// Package mocks is a generated GoMock package.
package mocks

import (
	domain "cowboy-app/internal/domain"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockCowboyRepository is a mock of CowboyRepository interface
type MockCowboyRepository struct {
	ctrl     *gomock.Controller
	recorder *MockCowboyRepositoryMockRecorder
}

// MockCowboyRepositoryMockRecorder is the mock recorder for MockCowboyRepository
type MockCowboyRepositoryMockRecorder struct {
	mock *MockCowboyRepository
}

// NewMockCowboyRepository creates a new mock instance
func NewMockCowboyRepository(ctrl *gomock.Controller) *MockCowboyRepository {
	mock := &MockCowboyRepository{ctrl: ctrl}
	mock.recorder = &MockCowboyRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCowboyRepository) EXPECT() *MockCowboyRepositoryMockRecorder {
	return m.recorder
}

// FindByName mocks base method
func (m *MockCowboyRepository) FindByName(arg0 string) (*domain.Cowboy, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByName", arg0)
	ret0, _ := ret[0].(*domain.Cowboy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByName indicates an expected call of FindByName
func (mr *MockCowboyRepositoryMockRecorder) FindByName(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByName", reflect.TypeOf((*MockCowboyRepository)(nil).FindByName), arg0)
}

// List mocks base method
func (m *MockCowboyRepository) List(arg0 *domain.CowboySearchCriteria) ([]*domain.Cowboy, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0)
	ret0, _ := ret[0].([]*domain.Cowboy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List
func (mr *MockCowboyRepositoryMockRecorder) List(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockCowboyRepository)(nil).List), arg0)
}

// UpdateHealthPoints mocks base method
func (m *MockCowboyRepository) UpdateHealthPoints(arg0 string, arg1 int32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateHealthPoints", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateHealthPoints indicates an expected call of UpdateHealthPoints
func (mr *MockCowboyRepositoryMockRecorder) UpdateHealthPoints(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateHealthPoints", reflect.TypeOf((*MockCowboyRepository)(nil).UpdateHealthPoints), arg0, arg1)
}