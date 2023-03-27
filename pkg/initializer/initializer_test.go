package initializer

import (
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd/api"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_initializer_GetNamespace(t *testing.T) {
	t.Run("should return default namespace if not present", func(t *testing.T) {
		// given
		sut := New()
		prevValue, present := os.LookupEnv(namespaceEnvName)
		defer resetEnv(t, namespaceEnvName, prevValue, present)
		err := os.Unsetenv("NAMESPACE")
		require.NoError(t, err)

		// when
		actual := sut.GetNamespace()

		// then
		assert.Equal(t, "default", actual)
	})
	t.Run("should return default namespace if empty string", func(t *testing.T) {
		// given
		sut := New()
		prevValue, present := os.LookupEnv(namespaceEnvName)
		defer resetEnv(t, namespaceEnvName, prevValue, present)
		err := os.Setenv("NAMESPACE", "")
		require.NoError(t, err)

		// when
		actual := sut.GetNamespace()

		// then
		assert.Equal(t, "default", actual)
	})
	t.Run("should return namespace from env", func(t *testing.T) {
		// given
		sut := New()
		prevValue, present := os.LookupEnv(namespaceEnvName)
		defer resetEnv(t, namespaceEnvName, prevValue, present)
		err := os.Setenv("NAMESPACE", "quark")
		require.NoError(t, err)

		// when
		actual := sut.GetNamespace()

		// then
		assert.Equal(t, "quark", actual)
	})
}

func resetEnv(t *testing.T, name, value string, present bool) {
	t.Helper()
	var err error
	if present {
		err = os.Setenv(name, value)
	} else {
		err = os.Unsetenv(name)
	}
	require.NoError(t, err)
}

func Test_initializer_CreateClientSet(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		defer func() {
			ctrl.GetConfigOrDie = config.GetConfigOrDie
		}()
		ctrl.GetConfigOrDie = func() *rest.Config {
			return &rest.Config{}
		}
		sut := initializer{}

		// when
		clientSet, err := sut.CreateClientSet()

		// then
		require.NoError(t, err)
		require.NotNil(t, clientSet)
	})

	t.Run("should return error on invalid config", func(t *testing.T) {
		// given
		defer func() {
			ctrl.GetConfigOrDie = config.GetConfigOrDie
		}()
		ctrl.GetConfigOrDie = func() *rest.Config {
			return &rest.Config{ExecProvider: &api.ExecConfig{}, AuthProvider: &api.AuthProviderConfig{}}
		}
		sut := initializer{}

		// when
		_, err := sut.CreateClientSet()

		// then
		require.Error(t, err)
	})
}

func Test_initializer_CreateCesRegistry(t *testing.T) {
	t.Run("should fail to create registry", func(t *testing.T) {
		// given
		sut := New()
		prevValue, present := os.LookupEnv(namespaceEnvName)
		defer resetEnv(t, namespaceEnvName, prevValue, present)
		err := os.Setenv(namespaceEnvName, "(!)//=)!%(?=(")
		require.NoError(t, err)

		// when
		actual, err := sut.CreateCesRegistry()

		// then
		require.Error(t, err)
		assert.Nil(t, actual)
		assert.ErrorContains(t, err, "parse \"http://etcd.(!)//=)!%(?=(.svc.cluster.local:4001\": invalid URL escape \"%(\"")
	})
	t.Run("should create registry", func(t *testing.T) {
		// given
		sut := New()
		prevValue, present := os.LookupEnv(namespaceEnvName)
		defer resetEnv(t, namespaceEnvName, prevValue, present)
		err := os.Setenv(namespaceEnvName, "ces")
		require.NoError(t, err)

		// when
		actual, err := sut.CreateCesRegistry()

		// then
		require.NoError(t, err)
		assert.NotNil(t, actual)
	})
}
