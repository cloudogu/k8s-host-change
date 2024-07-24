package alias

import (
	"context"
	"github.com/cloudogu/k8s-registry-lib/config"
	"github.com/stretchr/testify/mock"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	v1 "k8s.io/api/core/v1"
	"k8s.io/utils/strings/slices"
)

func Test_hostAliasGenerator_Generate(t *testing.T) {
	t.Run("if internalIP is used: fqdn should point to internalIP additional hosts are present", func(t *testing.T) {
		// given
		fqdn := "ecosystem.cloudogu.com"
		internalIP := "23.24.12.99"

		additionalHostOne := "prod.cloudogu.com"
		additionalHostTwo := "11.11.11.22"

		entries := config.Entries{
			"fqdn":                                 config.Value(fqdn),
			"k8s/use_internal_ip":                  config.Value("true"),
			"k8s/internal_ip":                      config.Value(internalIP),
			"containers/additional_hosts/host_one": config.Value(additionalHostOne),
			"containers/additional_hosts/host_two": config.Value(additionalHostTwo),
		}

		globalConfigRepoMock := newMockGlobalConfigGetter(t)
		globalConfigRepoMock.EXPECT().Get(mock.Anything).Return(config.CreateGlobalConfig(entries), nil)

		generator := HostAliasGenerator{
			globalConfigGetter: globalConfigRepoMock,
		}

		aliasFqdn := v1.HostAlias{IP: internalIP, Hostnames: []string{fqdn}}
		aliasOne := v1.HostAlias{IP: additionalHostOne, Hostnames: []string{"host_one"}}
		aliasTwo := v1.HostAlias{IP: additionalHostTwo, Hostnames: []string{"host_two"}}

		// when
		aliases, err := generator.Generate(context.TODO())

		// then
		require.NoError(t, err)
		assert.Equal(t, 3, len(aliases))
		assert.True(t, hasAlias(aliases, aliasFqdn))
		assert.True(t, hasAlias(aliases, aliasOne))
		assert.True(t, hasAlias(aliases, aliasTwo))
	})

	t.Run("if internalIP is not used: only additional hosts should be in hosts", func(t *testing.T) {
		// given
		fqdn := "ecosystem.cloudogu.com"

		additionalHostOne := "prod.cloudogu.com"
		additionalHostTwo := "11.11.11.22"

		entries := config.Entries{
			"fqdn":                                 config.Value(fqdn),
			"k8s/use_internal_ip":                  config.Value("false"),
			"containers/additional_hosts/host_one": config.Value(additionalHostOne),
			"containers/additional_hosts/host_two": config.Value(additionalHostTwo),
		}

		globalConfigRepoMock := newMockGlobalConfigGetter(t)
		globalConfigRepoMock.EXPECT().Get(mock.Anything).Return(config.CreateGlobalConfig(entries), nil)

		generator := HostAliasGenerator{
			globalConfigGetter: globalConfigRepoMock,
		}

		aliasOne := v1.HostAlias{IP: additionalHostOne, Hostnames: []string{"host_one"}}
		aliasTwo := v1.HostAlias{IP: additionalHostTwo, Hostnames: []string{"host_two"}}

		// when
		aliases, err := generator.Generate(context.TODO())

		// then
		require.NoError(t, err)
		assert.Equal(t, 2, len(aliases))
		assert.True(t, hasAlias(aliases, aliasOne))
		assert.True(t, hasAlias(aliases, aliasTwo))
	})

	t.Run("should fail on query fqdn error ", func(t *testing.T) {
		// given
		entries := config.Entries{}

		globalConfigRepoMock := newMockGlobalConfigGetter(t)
		globalConfigRepoMock.EXPECT().Get(mock.Anything).Return(config.CreateGlobalConfig(entries), nil)

		generator := HostAliasGenerator{
			globalConfigGetter: globalConfigRepoMock,
		}

		// when
		_, err := generator.Generate(context.TODO())

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "fqdn does not exist in global config")
	})

	t.Run("should fail on query internalIP flag error", func(t *testing.T) {
		// given
		fqdn := "ecosystem.cloudogu.com"

		entries := config.Entries{
			"fqdn": config.Value(fqdn),
		}

		globalConfigRepoMock := newMockGlobalConfigGetter(t)
		globalConfigRepoMock.EXPECT().Get(mock.Anything).Return(config.CreateGlobalConfig(entries), nil)

		generator := HostAliasGenerator{
			globalConfigGetter: globalConfigRepoMock,
		}

		// when
		_, err := generator.Generate(context.TODO())

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "k8s/use_internal_ip does not exist in global config")
	})

	t.Run("should fail on parse internalIP flag error", func(t *testing.T) {
		// given
		fqdn := "ecosystem.cloudogu.com"

		entries := config.Entries{
			"fqdn":                config.Value(fqdn),
			"k8s/use_internal_ip": config.Value("no string boolean"),
		}

		globalConfigRepoMock := newMockGlobalConfigGetter(t)
		globalConfigRepoMock.EXPECT().Get(mock.Anything).Return(config.CreateGlobalConfig(entries), nil)

		generator := HostAliasGenerator{
			globalConfigGetter: globalConfigRepoMock,
		}

		// when
		_, err := generator.Generate(context.TODO())

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to read config: failed to parse value 'no string boolean' of field 'k8s/use_internal_ip' in global config")
	})

	t.Run("should fail on query internalIP error", func(t *testing.T) {
		// given
		fqdn := "ecosystem.cloudogu.com"

		entries := config.Entries{
			"fqdn":                config.Value(fqdn),
			"k8s/use_internal_ip": config.Value("true"),
		}

		globalConfigRepoMock := newMockGlobalConfigGetter(t)
		globalConfigRepoMock.EXPECT().Get(mock.Anything).Return(config.CreateGlobalConfig(entries), nil)

		generator := HostAliasGenerator{
			globalConfigGetter: globalConfigRepoMock,
		}

		// when
		_, err := generator.Generate(context.TODO())

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "k8s/internal_ip does not exist in global config")
	})

	t.Run("should fail on parse internalIP error", func(t *testing.T) {
		// given
		fqdn := "ecosystem.cloudogu.com"

		entries := config.Entries{
			"fqdn":                config.Value(fqdn),
			"k8s/use_internal_ip": config.Value("true"),
			"k8s/internal_ip":     config.Value("fdsd2131"),
		}

		globalConfigRepoMock := newMockGlobalConfigGetter(t)
		globalConfigRepoMock.EXPECT().Get(mock.Anything).Return(config.CreateGlobalConfig(entries), nil)

		generator := HostAliasGenerator{
			globalConfigGetter: globalConfigRepoMock,
		}

		// when
		_, err := generator.Generate(context.TODO())

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to read config: failed to parse value 'fdsd2131' of field 'k8s/internal_ip' in global config: not a valid ip")
	})

	t.Run("should fail on query global config", func(t *testing.T) {
		// given
		globalConfigRepoMock := newMockGlobalConfigGetter(t)
		globalConfigRepoMock.EXPECT().Get(mock.Anything).Return(config.GlobalConfig{}, assert.AnError)

		generator := HostAliasGenerator{
			globalConfigGetter: globalConfigRepoMock,
		}

		// when
		_, err := generator.Generate(context.TODO())

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to get global config")
	})
}

func hasAlias(aliases []v1.HostAlias, alias v1.HostAlias) bool {
	for _, a := range aliases {
		if a.IP == alias.IP && slices.Equal(a.Hostnames, alias.Hostnames) {
			return true
		}
	}

	return false
}

func TestNewHostAliasGenerator(t *testing.T) {
	// when
	generator := NewHostAliasGenerator(nil)

	// then
	require.NotNil(t, generator)
}
