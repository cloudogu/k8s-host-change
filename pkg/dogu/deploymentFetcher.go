package dogu

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// NewDeploymentFetcher creates a new instance of a deployment fetcher which is used for retrieving dogu deployments.
func NewDeploymentFetcher(clientSet kubernetes.Interface) *deploymentFetcher {
	return &deploymentFetcher{clientSet: clientSet}
}

type deploymentFetcher struct {
	clientSet kubernetes.Interface
}

// FetchAll retrieves all dogu deployments in a given namespace.
// The 'dogu.name' label key is used for identifying dogu deployments.
func (f *deploymentFetcher) FetchAll(ctx context.Context, namespace string) ([]appsv1.Deployment, error) {
	selector := &metav1.LabelSelector{
		MatchExpressions: []metav1.LabelSelectorRequirement{{
			Key:      "dogu.name",
			Operator: metav1.LabelSelectorOpExists,
			Values:   nil,
		}},
	}

	options := metav1.ListOptions{
		LabelSelector: metav1.FormatLabelSelector(selector),
	}
	deploymentList, err := f.clientSet.AppsV1().Deployments(namespace).List(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("could not list deployments with selector 'dogu.name': %w", err)
	}

	return deploymentList.Items, nil
}
