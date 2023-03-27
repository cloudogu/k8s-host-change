package main

import (
	"github.com/cloudogu/cesapp-lib/registry"
	"github.com/cloudogu/k8s-host-change/pkg/hosts"
	"github.com/cloudogu/k8s-host-change/pkg/initializer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"k8s.io/client-go/kubernetes"
	"testing"
)

func Test_run(t *testing.T) {
	t.Run("should fail to create clientSet", func(t *testing.T) {
		// given
		oldInitConstructor := initializer.New
		oldUpdaterConstructor := hosts.NewHostAliasUpdater
		defer func() {
			initializer.New = oldInitConstructor
			hosts.NewHostAliasUpdater = oldUpdaterConstructor
		}()
		initializerMock := newMockInitializer(t)
		initializerMock.EXPECT().GetNamespace().Return("default").Once()
		initializerMock.EXPECT().CreateClientSet().Return(nil, assert.AnError).Once()
		initializer.New = func() initializer.Initializer {
			return initializerMock
		}

		// when
		err := run()

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})
	t.Run("should fail to create registry", func(t *testing.T) {
		// given
		oldInitConstructor := initializer.New
		oldUpdaterConstructor := hosts.NewHostAliasUpdater
		defer func() {
			initializer.New = oldInitConstructor
			hosts.NewHostAliasUpdater = oldUpdaterConstructor
		}()
		initializerMock := newMockInitializer(t)
		initializerMock.EXPECT().GetNamespace().Return("default").Once()
		initializerMock.EXPECT().CreateClientSet().Return(nil, nil).Once()
		initializerMock.EXPECT().CreateCesRegistry().Return(nil, assert.AnError).Once()
		initializer.New = func() initializer.Initializer {
			return initializerMock
		}

		// when
		err := run()

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})
	t.Run("should fail to update host aliases", func(t *testing.T) {
		// given
		oldInitConstructor := initializer.New
		oldUpdaterConstructor := hosts.NewHostAliasUpdater
		defer func() {
			initializer.New = oldInitConstructor
			hosts.NewHostAliasUpdater = oldUpdaterConstructor
		}()
		initializerMock := newMockInitializer(t)
		initializerMock.EXPECT().GetNamespace().Return("default").Once()
		initializerMock.EXPECT().CreateClientSet().Return(nil, nil).Once()
		initializerMock.EXPECT().CreateCesRegistry().Return(nil, nil).Once()
		initializer.New = func() initializer.Initializer {
			return initializerMock
		}
		updaterMock := newMockHostAliasUpdater(t)
		updaterMock.EXPECT().UpdateHosts(mock.Anything, "default").Return(assert.AnError).Once()
		hosts.NewHostAliasUpdater = func(_ kubernetes.Interface, _ registry.Registry) hosts.HostAliasUpdater {
			return updaterMock
		}

		// when
		err := run()

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})
	t.Run("should succeed", func(t *testing.T) {
		// given
		oldInitConstructor := initializer.New
		oldUpdaterConstructor := hosts.NewHostAliasUpdater
		defer func() {
			initializer.New = oldInitConstructor
			hosts.NewHostAliasUpdater = oldUpdaterConstructor
		}()
		initializerMock := newMockInitializer(t)
		initializerMock.EXPECT().GetNamespace().Return("default").Once()
		initializerMock.EXPECT().CreateClientSet().Return(nil, nil).Once()
		initializerMock.EXPECT().CreateCesRegistry().Return(nil, nil).Once()
		initializer.New = func() initializer.Initializer {
			return initializerMock
		}
		updaterMock := newMockHostAliasUpdater(t)
		updaterMock.EXPECT().UpdateHosts(mock.Anything, "default").Return(nil).Once()
		hosts.NewHostAliasUpdater = func(_ kubernetes.Interface, _ registry.Registry) hosts.HostAliasUpdater {
			return updaterMock
		}

		// when
		err := run()

		// then
		require.NoError(t, err)
	})
}
