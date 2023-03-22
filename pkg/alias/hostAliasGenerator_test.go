package alias

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	"k8s.io/utils/strings/slices"
	"testing"
)

func Test_hostAliasGenerator_Generate(t *testing.T) {
	t.Run("if internalIP is used: fqdn should point to internalIP additional hosts are present", func(t *testing.T) {
		// given
		globalConfigMock := newMockRegistryContext(t)
		expect := globalConfigMock.EXPECT()

		getExpect := expect.Get
		fqdn := "ecosystem.cloudogu.com"
		internalIP := "23.24.12.99"
		getExpect("fqdn").Return(fqdn, nil)
		getExpect("k8s/use_internal_ip").Return("true", nil)
		getExpect("k8s/internal_ip").Return(internalIP, nil)

		additionalHostOne := "prod.cloudogu.com"
		additionalHostTwo := "11.11.11.22"
		additionalHosts := map[string]string{"containers/additional_hosts/host_one": additionalHostOne,
			"containers/additional_hosts/host_two": additionalHostTwo}
		expect.GetAll().Return(additionalHosts, nil)

		generator := hostAliasGenerator{
			globalConfig: globalConfigMock,
		}

		aliasFqdn := v1.HostAlias{IP: internalIP, Hostnames: []string{fqdn}}
		aliasOne := v1.HostAlias{IP: additionalHostOne, Hostnames: []string{"host_one"}}
		aliasTwo := v1.HostAlias{IP: additionalHostTwo, Hostnames: []string{"host_two"}}

		// when
		aliases, err := generator.Generate()

		// then
		require.NoError(t, err)
		assert.Equal(t, 3, len(aliases))
		assert.True(t, hasAlias(aliases, aliasFqdn))
		assert.True(t, hasAlias(aliases, aliasOne))
		assert.True(t, hasAlias(aliases, aliasTwo))
	})

	t.Run("if internalIP is not used: only additional hosts should be in hosts", func(t *testing.T) {
		// given
		globalConfigMock := newMockRegistryContext(t)
		expect := globalConfigMock.EXPECT()

		getExpect := expect.Get
		fqdn := "ecosystem.cloudogu.com"
		getExpect("fqdn").Return(fqdn, nil)
		getExpect("k8s/use_internal_ip").Return("false", nil)

		additionalHostOne := "prod.cloudogu.com"
		additionalHostTwo := "11.11.11.22"
		additionalHosts := map[string]string{"containers/additional_hosts/host_one": additionalHostOne,
			"containers/additional_hosts/host_two": additionalHostTwo}
		expect.GetAll().Return(additionalHosts, nil)

		generator := hostAliasGenerator{
			globalConfig: globalConfigMock,
		}

		aliasOne := v1.HostAlias{IP: additionalHostOne, Hostnames: []string{"host_one"}}
		aliasTwo := v1.HostAlias{IP: additionalHostTwo, Hostnames: []string{"host_two"}}

		// when
		aliases, err := generator.Generate()

		// then
		require.NoError(t, err)
		assert.Equal(t, 2, len(aliases))
		assert.True(t, hasAlias(aliases, aliasOne))
		assert.True(t, hasAlias(aliases, aliasTwo))
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
