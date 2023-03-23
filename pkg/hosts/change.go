package hosts

import (
	"context"
	"fmt"
	"github.com/cloudogu/k8s-host-change/pkg/alias"
	"github.com/cloudogu/k8s-host-change/pkg/deployment"
	"github.com/cloudogu/k8s-host-change/pkg/dogu"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type hostAliasUpdater struct {
	generator hostAliasGenerator
	fetcher   doguDeploymentFetcher
	updater   deploymentUpdater
}

// NewHostAliasUpdater is used to create a new instance of hostAliasUpdater.
var NewHostAliasUpdater = func(clientSet kubernetes.Interface, cesReg cesRegistry) *hostAliasUpdater {
	return &hostAliasUpdater{
		generator: alias.NewHostAliasGenerator(cesReg.GlobalConfig()),
		fetcher:   dogu.NewDeploymentFetcher(clientSet),
		updater:   deployment.NewUpdater(clientSet),
	}
}

// UpdateHosts updates all dogu deployments with host information like fqdn, internal ip and additional hosts from ces registry.
func (hau *hostAliasUpdater) UpdateHosts(ctx context.Context, namespace string) error {
	logger := log.FromContext(ctx)
	logger.Info("Update host entries in dogu deployments")
	hostAliases, err := hau.generator.Generate()
	if err != nil {
		return fmt.Errorf("failed to generate host aliases: %w", err)
	}
	if len(hostAliases) > 0 {
		logger.Info(fmt.Sprintf("Use aliases: %s", hostAliases))
	} else {
		logger.Info("Delete all aliases from dogu deployments")
	}

	logger.Info("Fetch all dogu deployments")
	deployments, err := hau.fetcher.FetchAll(ctx, namespace)
	if err != nil {
		return fmt.Errorf("failed to fetch dogu deployments: %w", err)
	}

	logger.Info("Update deployments with host aliases")
	err = hau.updater.UpdateHostAliases(ctx, namespace, deployments, hostAliases)
	if err != nil {
		return fmt.Errorf("failed to update host-aliases of dogu deployments in cluster: %w", err)
	}

	return nil
}
