package hosts

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

type HostAliasUpdater interface {
	// UpdateHosts updates all dogu deployments with host information like fqdn, internal ip and additional hosts from ces registry.
	UpdateHosts(ctx context.Context, namespace string) error
}

type hostAliasGenerator interface {
	// Generate patches the given deployment with the host configuration provided.
	Generate() (hostAliases []corev1.HostAlias, err error)
}

type doguDeploymentFetcher interface {
	// FetchAll retrieves all dogu deployments in a given namespace.
	FetchAll(ctx context.Context, namespace string) ([]appsv1.Deployment, error)
}

type deploymentUpdater interface {
	// UpdateHostAliases replaces the host aliases in the given deployments.
	UpdateHostAliases(ctx context.Context, namespace string, deployments []appsv1.Deployment, aliases []corev1.HostAlias) error
}
