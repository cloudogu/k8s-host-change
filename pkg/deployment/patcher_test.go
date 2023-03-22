package deployment

import (
	"testing"

	"github.com/stretchr/testify/assert"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

func Test_hostAliasPatcher_Patch(t *testing.T) {
	// given
	deployments := []appsv1.Deployment{
		{
			Spec: appsv1.DeploymentSpec{
				Template: corev1.PodTemplateSpec{
					Spec: corev1.PodSpec{},
				},
			},
		},
		{
			Spec: appsv1.DeploymentSpec{
				Template: corev1.PodTemplateSpec{
					Spec: corev1.PodSpec{
						HostAliases: []corev1.HostAlias{{
							IP:        "5.6.7.8",
							Hostnames: []string{"abc"},
						}},
					},
				},
			},
		},
	}
	aliases := []corev1.HostAlias{
		{
			IP:        "1.2.3.4",
			Hostnames: []string{"xyz"},
		},
		{
			IP:        "0.0.0.0",
			Hostnames: []string{"ads.example.com"},
		},
	}
	sut := NewHostAliasPatcher()

	// when
	sut.Patch(deployments, aliases)

	// then
	assert.Equal(t, aliases, deployments[0].Spec.Template.Spec.HostAliases)
	assert.Equal(t, aliases, deployments[1].Spec.Template.Spec.HostAliases)
}
