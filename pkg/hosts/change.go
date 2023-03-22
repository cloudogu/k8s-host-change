package hosts

import (
	"context"
	"github.com/cloudogu/cesapp-lib/registry"
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

func NewHostAliasUpdater(clientSet kubernetes.Interface, cesReg registry.Registry) (*hostAliasUpdater, error) {
	return &hostAliasUpdater{
		generator: alias.NewHostAliasGenerator(cesReg),
		fetcher:   dogu.NewDeploymentFetcher(clientSet),
		patcher:   deployment.NewHostAliasPatcher(),
		updater:   deployment.NewUpdater(clientSet),
	}, nil
}

func (hau *hostAliasUpdater) UpdateHosts(ctx context.Context, namespace string) error {
	hostAliases, err := hau.generator.Generate()
	if err != nil {
		return err
	}

	deployments, err := hau.fetcher.FetchAll(ctx, namespace)
	if err != nil {
		return err
	}

	hau.patcher.Patch(deployments, hostAliases)

	err = hau.updater.Update(ctx, namespace, deployments)
	if err != nil {
		return err
	}

	return nil
}
