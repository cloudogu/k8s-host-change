package alias

import (
	"context"
	"fmt"
	"github.com/cloudogu/k8s-registry-lib/config"
	"net"
	"strconv"
	"strings"

	v1 "k8s.io/api/core/v1"
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

type HostAliasGenerator struct {
	globalConfigGetter globalConfigGetter
}

// NewHostAliasGenerator creates a generator with the ability to return host aliases from the configured internal ip, additional hosts and fqdn.
func NewHostAliasGenerator(globalConfigGetter globalConfigGetter) *HostAliasGenerator {
	return &HostAliasGenerator{
		globalConfigGetter: globalConfigGetter,
	}
}

// Generate patches the given deployment with the host configuration provided.
func (d *HostAliasGenerator) Generate(ctx context.Context) (hostAliases []v1.HostAlias, err error) {
	cfg, err := d.getGeneratorConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	if cfg.useInternalIP {
		splitDnsHostAlias := v1.HostAlias{
			IP:        cfg.internalIP.String(),
			Hostnames: []string{cfg.fqdn},
		}
		hostAliases = append(hostAliases, splitDnsHostAlias)
	}

	for hostName, ip := range cfg.additionalHosts {
		addHostAlias := v1.HostAlias{
			IP:        ip,
			Hostnames: []string{hostName},
		}
		hostAliases = append(hostAliases, addHostAlias)
	}

	return hostAliases, nil
}

// getGeneratorConfig reads hosts-specific keys from the global configuration and creates a generatorConfig object.
func (d *HostAliasGenerator) getGeneratorConfig(ctx context.Context) (*generatorConfig, error) {
	globalCfg, err := d.globalConfigGetter.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get global config: %w", err)
	}

	fqdn, err := d.getFQDN(globalCfg)
	if err != nil {
		return nil, err
	}

	hostsConfig := &generatorConfig{
		fqdn: fqdn,
	}

	hostsConfig.useInternalIP, err = d.isInternalIPUsed(globalCfg)
	if err != nil {
		return nil, err
	}

	if hostsConfig.useInternalIP {
		hostsConfig.internalIP, err = d.getInternalIP(globalCfg)
		if err != nil {
			return nil, err
		}
	}

	hostsConfig.additionalHosts, err = d.retrieveAdditionalHosts(globalCfg)
	if err != nil {
		return nil, err
	}

	return hostsConfig, nil
}

func (d *HostAliasGenerator) getFQDN(globalCfg config.GlobalConfig) (string, error) {
	fqdn, ok := globalCfg.Get(fqdnKey)
	if !ok {
		return "", fmt.Errorf("key: %s does not exist in global config", fqdnKey)
	}

	return fqdn.String(), nil
}

func (d *HostAliasGenerator) isInternalIPUsed(globalCfg config.GlobalConfig) (useInternalIP bool, err error) {
	useInternalIPRaw, ok := globalCfg.Get(useInternalIPKey)
	if !ok {
		return false, fmt.Errorf("key: %s does not exist in global config", useInternalIPKey)
	}

	useInternalIP, err = strconv.ParseBool(useInternalIPRaw.String())
	if err != nil {
		return false, fmt.Errorf("failed to parse value '%s' of field '%s' in global config: %w", useInternalIPRaw, useInternalIPKey, err)
	}

	return useInternalIP, nil
}

func (d *HostAliasGenerator) getInternalIP(globalCfg config.GlobalConfig) (net.IP, error) {
	internalIPRaw, ok := globalCfg.Get(internalIPKey)
	if !ok {
		return nil, fmt.Errorf("key: %s does not exist in global config", internalIPKey)
	}

	ip := net.ParseIP(internalIPRaw.String())
	if ip == nil {
		return nil, fmt.Errorf("failed to parse value '%s' of field '%s' in global config: not a valid ip", internalIPRaw, internalIPKey)
	}

	return ip, nil
}

func (d *HostAliasGenerator) retrieveAdditionalHosts(globalCfg config.GlobalConfig) (map[string]string, error) {
	globalCfgEntries := globalCfg.GetAll()

	additionalHosts := map[string]string{}
	for key, value := range globalCfgEntries {
		if strings.HasPrefix(key.String(), additionalHostsPrefix) {
			hostName := strings.TrimPrefix(key.String(), additionalHostsPrefix)
			additionalHosts[hostName] = value.String()
		}
	}

	return additionalHosts, nil
}
