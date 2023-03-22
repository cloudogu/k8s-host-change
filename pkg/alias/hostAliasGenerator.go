package alias

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	v1 "k8s.io/api/core/v1"

	"github.com/cloudogu/cesapp-lib/registry"
)

const (
	useInternalIPKey      = "k8s/use_internal_ip"
	internalIPKey         = "k8s/internal_ip"
	fqdnKey               = "fqdn"
	additionalHostsPrefix = "containers/additional_hosts/"
)

type generatorConfig struct {
	fqdn            string
	useInternalIP   bool
	internalIP      net.IP
	additionalHosts map[string]string
}

type hostAliasGenerator struct {
	globalConfig registryContext
}

// NewHostAliasGenerator creates a generator with the ability to return host aliases from the configured internal ip, additional hosts and fqdn.
func NewHostAliasGenerator(globalConfig registryContext) *hostAliasGenerator {
	return &hostAliasGenerator{
		globalConfig: globalConfig,
	}
}

// Generate patches the given deployment with the host configuration provided.
func (d *hostAliasGenerator) Generate() (hostAliases []v1.HostAlias, err error) {
	config, err := d.getGeneratorConfig()
	if err != nil {
		return nil, err
	}

	if config.useInternalIP {
		splitDnsHostAlias := v1.HostAlias{
			IP:        config.internalIP.String(),
			Hostnames: []string{config.fqdn},
		}
		hostAliases = append(hostAliases, splitDnsHostAlias)
	}

	for hostName, ip := range config.additionalHosts {
		addHostAlias := v1.HostAlias{
			IP:        ip,
			Hostnames: []string{hostName},
		}
		hostAliases = append(hostAliases, addHostAlias)
	}

	return hostAliases, nil
}

// getGeneratorConfig reads hosts-specific keys from the global configuration and creates a generatorConfig object.
func (d *hostAliasGenerator) getGeneratorConfig() (*generatorConfig, error) {
	fqdn, err := d.getFQDN()
	if err != nil {
		return nil, err
	}

	hostsConfig := &generatorConfig{
		fqdn: fqdn,
	}

	hostsConfig.useInternalIP, err = d.isInternalIPUsed()
	if err != nil {
		return nil, err
	}

	if hostsConfig.useInternalIP {
		hostsConfig.internalIP, err = d.getInternalIP()
		if err != nil {
			return nil, err
		}
	}

	hostsConfig.additionalHosts, err = d.retrieveAdditionalHosts()
	if err != nil {
		return nil, err
	}

	return hostsConfig, nil
}

func (d *hostAliasGenerator) getFQDN() (string, error) {
	fqdn, err := d.globalConfig.Get(fqdnKey)
	if err != nil {
		return "", err
	}

	return fqdn, nil
}

func (d *hostAliasGenerator) isInternalIPUsed() (useInternalIP bool, err error) {
	useInternalIPRaw, err := d.globalConfig.Get(useInternalIPKey)
	if err != nil && !registry.IsKeyNotFoundError(err) {
		return false, err
	} else if err == nil {
		useInternalIP, err = strconv.ParseBool(useInternalIPRaw)
		if err != nil {
			return false, fmt.Errorf("failed to parse value '%s' of field 'k8s/use_internal_ip' in global generatorConfig: %w", useInternalIPRaw, err)
		}
	}

	return useInternalIP, nil
}

func (d *hostAliasGenerator) getInternalIP() (internalIP net.IP, err error) {
	internalIPRaw, err := d.globalConfig.Get(internalIPKey)
	if err != nil && !registry.IsKeyNotFoundError(err) {
		return nil, err
	} else if err == nil {
		internalIP, err = parseInternalIP(internalIPRaw)
		if err != nil {
			return nil, err
		}
	}

	return internalIP, nil
}

func (d *hostAliasGenerator) retrieveAdditionalHosts() (map[string]string, error) {
	globalConfig, err := d.globalConfig.GetAll()
	if err != nil {
		return nil, err
	}

	additionalHosts := map[string]string{}
	for key, value := range globalConfig {
		if strings.HasPrefix(key, additionalHostsPrefix) {
			hostName := strings.TrimPrefix(key, additionalHostsPrefix)
			additionalHosts[hostName] = value
		}
	}
	return additionalHosts, nil
}

func parseInternalIP(raw string) (net.IP, error) {
	ip := net.ParseIP(raw)
	if ip == nil {
		return nil, fmt.Errorf("failed to parse value '%s' of field 'k8s/internal_ip' in global generatorConfig: not a valid ip", raw)
	}

	return ip, nil
}
