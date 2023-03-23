package hosts

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/cloudogu/k8s-host-change/pkg/alias"
	"github.com/cloudogu/k8s-host-change/pkg/deployment"
	"github.com/cloudogu/k8s-host-change/pkg/dogu"
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
	logger := log.FromContext(ctx)
	logger.Info("Update host entries in dogu deployments")
	hostAliases, err := hau.generator.Generate()
	if err != nil {
		return fmt.Errorf("failed to generate host aliases: %w", err)
	}
	logger.Info("Use aliases: %s", hostAliases)

	deployments, err := hau.fetcher.FetchAll(ctx, namespace)
	if err != nil {
		return fmt.Errorf("failed to fetch dogu deployments: %w", err)
	}

	previousHostAliases := make(map[string][]corev1.HostAlias)
	for _, deploy := range deployments {
		previousHostAliases[deploy.Name] = deploy.Spec.Template.Spec.HostAliases
	}

	hau.patcher.Patch(deployments, hostAliases)

	logger.Info("Update deployments")
	err = hau.updater.Update(ctx, namespace, deployments)
	if err != nil {
		logger.Error(err, "Failed to update dogu deployments: rolling back")

		rollbackErr := hau.rollback(ctx, namespace, previousHostAliases)
		if rollbackErr != nil {
			err = multierror.Append(err, rollbackErr)
		}

		return fmt.Errorf("failed to update host-aliases of dogu deployments in cluster: %w", err)
	}

	return nil
}

func (hau *hostAliasUpdater) rollback(ctx context.Context, namespace string, previousHostAliases map[string][]corev1.HostAlias) error {
	deployments, err := hau.fetcher.FetchAll(ctx, namespace)
	if err != nil {
		return fmt.Errorf("failed to fetch dogu deployments on rollback: %w", err)
	}

	for i, deploy := range deployments {
		aliases, exists := previousHostAliases[deploy.Name]
		if exists {
			deployments[i].Spec.Template.Spec.HostAliases = aliases
		}
	}

	err = hau.updater.Update(ctx, namespace, deployments)
	if err != nil {
		return fmt.Errorf("failed to rollback dogu deployments: %w", err)
	}

	return nil
}
