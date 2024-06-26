package alias

import (
	"context"
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

		globalConfigMock := newMockGlobalConfigValueGetter(t)
		globalConfigMock.EXPECT().Get(mock.Anything, "fqdn").Return(fqdn, nil)
		globalConfigMock.EXPECT().Get(mock.Anything, "k8s/use_internal_ip").Return("true", nil)
		globalConfigMock.EXPECT().Get(mock.Anything, "k8s/internal_ip").Return(internalIP, nil)

		additionalHostOne := "prod.cloudogu.com"
		additionalHostTwo := "11.11.11.22"
		additionalHosts := map[string]string{"containers/additional_hosts/host_one": additionalHostOne,
			"containers/additional_hosts/host_two": additionalHostTwo}

		globalConfigMock.EXPECT().GetAll(mock.Anything).Return(additionalHosts, nil)

		generator := hostAliasGenerator{
			globalConfig: globalConfigMock,
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

		globalConfigMock := newMockGlobalConfigValueGetter(t)
		globalConfigMock.EXPECT().Get(mock.Anything, "fqdn").Return(fqdn, nil)
		globalConfigMock.EXPECT().Get(mock.Anything, "k8s/use_internal_ip").Return("false", nil)

		additionalHostOne := "prod.cloudogu.com"
		additionalHostTwo := "11.11.11.22"
		additionalHosts := map[string]string{"containers/additional_hosts/host_one": additionalHostOne,
			"containers/additional_hosts/host_two": additionalHostTwo}

		globalConfigMock.EXPECT().GetAll(mock.Anything).Return(additionalHosts, nil)

		generator := hostAliasGenerator{
			globalConfig: globalConfigMock,
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
		globalConfigMock := newMockGlobalConfigValueGetter(t)
		globalConfigMock.EXPECT().Get(mock.Anything, "fqdn").Return("", assert.AnError)

		generator := hostAliasGenerator{
			globalConfig: globalConfigMock,
		}

		// when
		_, err := generator.Generate(context.TODO())

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to read config: failed to get value for 'fqdn' from global config")
	})

	t.Run("should fail on query internalIP flag error", func(t *testing.T) {
		// given
		fqdn := "ecosystem.cloudogu.com"

		globalConfigMock := newMockGlobalConfigValueGetter(t)
		globalConfigMock.EXPECT().Get(mock.Anything, "fqdn").Return(fqdn, nil)
		globalConfigMock.EXPECT().Get(mock.Anything, "k8s/use_internal_ip").Return("", assert.AnError)

		generator := hostAliasGenerator{
			globalConfig: globalConfigMock,
		}

		// when
		_, err := generator.Generate(context.TODO())

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to read config: failed to get value for 'k8s/use_internal_ip' from global config")
	})

	t.Run("should fail on parse internalIP flag error", func(t *testing.T) {
		// given
		fqdn := "ecosystem.cloudogu.com"

		globalConfigMock := newMockGlobalConfigValueGetter(t)
		globalConfigMock.EXPECT().Get(mock.Anything, "fqdn").Return(fqdn, nil)
		globalConfigMock.EXPECT().Get(mock.Anything, "k8s/use_internal_ip").Return("no string boolean", nil)

		generator := hostAliasGenerator{
			globalConfig: globalConfigMock,
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

		globalConfigMock := newMockGlobalConfigValueGetter(t)
		globalConfigMock.EXPECT().Get(mock.Anything, "fqdn").Return(fqdn, nil)
		globalConfigMock.EXPECT().Get(mock.Anything, "k8s/use_internal_ip").Return("true", nil)
		globalConfigMock.EXPECT().Get(mock.Anything, "k8s/internal_ip").Return("", assert.AnError)

		generator := hostAliasGenerator{
			globalConfig: globalConfigMock,
		}

		// when
		_, err := generator.Generate(context.TODO())

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to read config: failed to get value for field 'k8s/internal_ip' from global config")
	})

	t.Run("should fail on parse internalIP error", func(t *testing.T) {
		// given
		fqdn := "ecosystem.cloudogu.com"

		globalConfigMock := newMockGlobalConfigValueGetter(t)
		globalConfigMock.EXPECT().Get(mock.Anything, "fqdn").Return(fqdn, nil)
		globalConfigMock.EXPECT().Get(mock.Anything, "k8s/use_internal_ip").Return("true", nil)
		globalConfigMock.EXPECT().Get(mock.Anything, "k8s/internal_ip").Return("fdsd2131", nil)

		generator := hostAliasGenerator{
			globalConfig: globalConfigMock,
		}

		// when
		_, err := generator.Generate(context.TODO())

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to read config: failed to parse value 'fdsd2131' of field 'k8s/internal_ip' in global config: not a valid ip")
	})

	t.Run("should fail on query additional hosts error", func(t *testing.T) {
		// given
		fqdn := "ecosystem.cloudogu.com"

		globalConfigMock := newMockGlobalConfigValueGetter(t)
		globalConfigMock.EXPECT().Get(mock.Anything, "fqdn").Return(fqdn, nil)
		globalConfigMock.EXPECT().Get(mock.Anything, "k8s/use_internal_ip").Return("false", nil)
		globalConfigMock.EXPECT().GetAll(mock.Anything).Return(nil, assert.AnError)

		generator := hostAliasGenerator{
			globalConfig: globalConfigMock,
		}

		// when
		_, err := generator.Generate(context.TODO())

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to get all keys from config")
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
