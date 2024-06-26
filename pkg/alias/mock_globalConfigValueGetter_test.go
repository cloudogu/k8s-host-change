// Code generated by mockery v2.42.1. DO NOT EDIT.

package alias

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// mockGlobalConfigValueGetter is an autogenerated mock type for the globalConfigValueGetter type
type mockGlobalConfigValueGetter struct {
	mock.Mock
}

type mockGlobalConfigValueGetter_Expecter struct {
	mock *mock.Mock
}

func (_m *mockGlobalConfigValueGetter) EXPECT() *mockGlobalConfigValueGetter_Expecter {
	return &mockGlobalConfigValueGetter_Expecter{mock: &_m.Mock}
}

// Get provides a mock function with given fields: ctx, key
func (_m *mockGlobalConfigValueGetter) Get(ctx context.Context, key string) (string, error) {
	ret := _m.Called(ctx, key)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (string, error)); ok {
		return rf(ctx, key)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, key)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockGlobalConfigValueGetter_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type mockGlobalConfigValueGetter_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - ctx context.Context
//   - key string
func (_e *mockGlobalConfigValueGetter_Expecter) Get(ctx interface{}, key interface{}) *mockGlobalConfigValueGetter_Get_Call {
	return &mockGlobalConfigValueGetter_Get_Call{Call: _e.mock.On("Get", ctx, key)}
}

func (_c *mockGlobalConfigValueGetter_Get_Call) Run(run func(ctx context.Context, key string)) *mockGlobalConfigValueGetter_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *mockGlobalConfigValueGetter_Get_Call) Return(_a0 string, _a1 error) *mockGlobalConfigValueGetter_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockGlobalConfigValueGetter_Get_Call) RunAndReturn(run func(context.Context, string) (string, error)) *mockGlobalConfigValueGetter_Get_Call {
	_c.Call.Return(run)
	return _c
}

// GetAll provides a mock function with given fields: ctx
func (_m *mockGlobalConfigValueGetter) GetAll(ctx context.Context) (map[string]string, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetAll")
	}

	var r0 map[string]string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (map[string]string, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) map[string]string); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]string)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockGlobalConfigValueGetter_GetAll_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAll'
type mockGlobalConfigValueGetter_GetAll_Call struct {
	*mock.Call
}

// GetAll is a helper method to define mock.On call
//   - ctx context.Context
func (_e *mockGlobalConfigValueGetter_Expecter) GetAll(ctx interface{}) *mockGlobalConfigValueGetter_GetAll_Call {
	return &mockGlobalConfigValueGetter_GetAll_Call{Call: _e.mock.On("GetAll", ctx)}
}

func (_c *mockGlobalConfigValueGetter_GetAll_Call) Run(run func(ctx context.Context)) *mockGlobalConfigValueGetter_GetAll_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *mockGlobalConfigValueGetter_GetAll_Call) Return(_a0 map[string]string, _a1 error) *mockGlobalConfigValueGetter_GetAll_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockGlobalConfigValueGetter_GetAll_Call) RunAndReturn(run func(context.Context) (map[string]string, error)) *mockGlobalConfigValueGetter_GetAll_Call {
	_c.Call.Return(run)
	return _c
}

// newMockGlobalConfigValueGetter creates a new instance of mockGlobalConfigValueGetter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockGlobalConfigValueGetter(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockGlobalConfigValueGetter {
	mock := &mockGlobalConfigValueGetter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}