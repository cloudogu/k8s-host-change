package hosts

import (
	"context"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

type HostAliasUpdater interface {
	UpdateHosts(ctx context.Context, namespace string) error
}

type hostAliasGenerator interface {
	Generate() (hostAliases []corev1.HostAlias, err error)
}

type doguDeploymentFetcher interface {
	FetchAll(ctx context.Context, namespace string) ([]appsv1.Deployment, error)
}

type deploymentUpdater interface {
	UpdateHostAliases(ctx context.Context, namespace string, deployments []appsv1.Deployment, aliases []corev1.HostAlias) error
}
