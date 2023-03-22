package hosts

import (
	"context"
	"fmt"
	"github.com/cloudogu/k8s-host-change/pkg/alias"
	"github.com/cloudogu/k8s-host-change/pkg/deployment"
	"github.com/cloudogu/k8s-host-change/pkg/dogu"
	"k8s.io/client-go/kubernetes"
)

type hostAliasUpdater struct {
	generator hostAliasGenerator
	fetcher   doguDeploymentFetcher
	patcher   hostAliasPatcher
	updater   deploymentUpdater
}

var NewHostAliasUpdater = func(clientSet kubernetes.Interface, cesReg cesRegistry) *hostAliasUpdater {
	return &hostAliasUpdater{
		generator: alias.NewHostAliasGenerator(cesReg.GlobalConfig()),
		fetcher:   dogu.NewDeploymentFetcher(clientSet),
		patcher:   deployment.NewHostAliasPatcher(),
		updater:   deployment.NewUpdater(clientSet),
	}
}

// UpdateHosts updates all dogu deployments with host information like fqdn, internal ip and additional hosts from ces registry.
func (hau *hostAliasUpdater) UpdateHosts(ctx context.Context, namespace string) error {
	hostAliases, err := hau.generator.Generate()
	if err != nil {
		return fmt.Errorf("failed to generate host aliases: %w", err)
	}

	deployments, err := hau.fetcher.FetchAll(ctx, namespace)
	if err != nil {
		return fmt.Errorf("failed to fetch dogu deployments: %w", err)
	}

	hau.patcher.Patch(deployments, hostAliases)

	err = hau.updater.Update(ctx, namespace, deployments)
	if err != nil {
		return fmt.Errorf("failed to update host-aliases of dogu deployments in cluster: %w", err)
	}

	return nil
}
