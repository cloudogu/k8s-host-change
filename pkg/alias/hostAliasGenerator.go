package alias

import (
	"context"
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
	globalConfig globalConfigValueGetter
}

// NewHostAliasGenerator creates a generator with the ability to return host aliases from the configured internal ip, additional hosts and fqdn.
func NewHostAliasGenerator(globalConfig globalConfigValueGetter) *hostAliasGenerator {
	return &hostAliasGenerator{
		globalConfig: globalConfig,
	}
}

// Generate patches the given deployment with the host configuration provided.
func (d *hostAliasGenerator) Generate(ctx context.Context) (hostAliases []v1.HostAlias, err error) {
	config, err := d.getGeneratorConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
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
func (d *hostAliasGenerator) getGeneratorConfig(ctx context.Context) (*generatorConfig, error) {
	fqdn, err := d.getFQDN(ctx)
	if err != nil {
		return nil, err
	}

	hostsConfig := &generatorConfig{
		fqdn: fqdn,
	}

	hostsConfig.useInternalIP, err = d.isInternalIPUsed(ctx)
	if err != nil {
		return nil, err
	}

	if hostsConfig.useInternalIP {
		hostsConfig.internalIP, err = d.getInternalIP(ctx)
		if err != nil {
			return nil, err
		}
	}

	hostsConfig.additionalHosts, err = d.retrieveAdditionalHosts(ctx)
	if err != nil {
		return nil, err
	}

	return hostsConfig, nil
}

func (d *hostAliasGenerator) getFQDN(ctx context.Context) (string, error) {
	fqdn, err := d.globalConfig.Get(ctx, fqdnKey)
	if err != nil {
		return "", fmt.Errorf("failed to get value for '%s' from global config: %w", fqdnKey, err)
	}

	return fqdn, nil
}

func (d *hostAliasGenerator) isInternalIPUsed(ctx context.Context) (useInternalIP bool, err error) {
	useInternalIPRaw, err := d.globalConfig.Get(ctx, useInternalIPKey)
	if err != nil && !registry.IsKeyNotFoundError(err) {
		return false, fmt.Errorf("failed to get value for '%s' from global config: %w", useInternalIPKey, err)
	} else if err == nil {
		useInternalIP, err = strconv.ParseBool(useInternalIPRaw)
		if err != nil {
			return false, fmt.Errorf("failed to parse value '%s' of field '%s' "+
				"in global config: %w", useInternalIPRaw, useInternalIPKey, err)
		}
	}

	return useInternalIP, nil
}

func (d *hostAliasGenerator) getInternalIP(ctx context.Context) (internalIP net.IP, err error) {
	internalIPRaw, err := d.globalConfig.Get(ctx, internalIPKey)
	if err != nil && !registry.IsKeyNotFoundError(err) {
		return nil, fmt.Errorf("failed to get value for field '%s' from global config: %w", internalIPKey, err)
	} else if err == nil {
		internalIP, err = parseInternalIP(internalIPRaw)
		if err != nil {
			return nil, err
		}
	}

	return internalIP, nil
}

func (d *hostAliasGenerator) retrieveAdditionalHosts(ctx context.Context) (map[string]string, error) {
	globalConfig, err := d.globalConfig.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all keys from config: %w", err)
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
		return nil, fmt.Errorf("failed to parse value '%s' of field '%s' in global config: "+
			"not a valid ip", raw, internalIPKey)
	}

	return ip, nil
}
