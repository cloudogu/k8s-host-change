// Code generated by mockery v2.20.0. DO NOT EDIT.

package main

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// mockHostAliasUpdater is an autogenerated mock type for the hostAliasUpdater type
type mockHostAliasUpdater struct {
	mock.Mock
}

type mockHostAliasUpdater_Expecter struct {
	mock *mock.Mock
}

func (_m *mockHostAliasUpdater) EXPECT() *mockHostAliasUpdater_Expecter {
	return &mockHostAliasUpdater_Expecter{mock: &_m.Mock}
}

// UpdateHosts provides a mock function with given fields: ctx, namespace
func (_m *mockHostAliasUpdater) UpdateHosts(ctx context.Context, namespace string) error {
	ret := _m.Called(ctx, namespace)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, namespace)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockHostAliasUpdater_UpdateHosts_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateHosts'
type mockHostAliasUpdater_UpdateHosts_Call struct {
	*mock.Call
}

// UpdateHosts is a helper method to define mock.On call
//   - ctx context.Context
//   - namespace string
func (_e *mockHostAliasUpdater_Expecter) UpdateHosts(ctx interface{}, namespace interface{}) *mockHostAliasUpdater_UpdateHosts_Call {
	return &mockHostAliasUpdater_UpdateHosts_Call{Call: _e.mock.On("UpdateHosts", ctx, namespace)}
}

func (_c *mockHostAliasUpdater_UpdateHosts_Call) Run(run func(ctx context.Context, namespace string)) *mockHostAliasUpdater_UpdateHosts_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *mockHostAliasUpdater_UpdateHosts_Call) Return(_a0 error) *mockHostAliasUpdater_UpdateHosts_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockHostAliasUpdater_UpdateHosts_Call) RunAndReturn(run func(context.Context, string) error) *mockHostAliasUpdater_UpdateHosts_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTnewMockHostAliasUpdater interface {
	mock.TestingT
	Cleanup(func())
}

// newMockHostAliasUpdater creates a new instance of mockHostAliasUpdater. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockHostAliasUpdater(t mockConstructorTestingTnewMockHostAliasUpdater) *mockHostAliasUpdater {
	mock := &mockHostAliasUpdater{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}