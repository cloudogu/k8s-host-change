package deployment

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

type hostAliasPatcher struct {
}

func NewHostAliasPatcher() *hostAliasPatcher {
	return &hostAliasPatcher{}
}

func (hap *hostAliasPatcher) Patch(deployments []appsv1.Deployment, aliases []corev1.HostAlias) {
	for i := range deployments {
		deployments[i].Spec.Template.Spec.HostAliases = aliases
	}
}
