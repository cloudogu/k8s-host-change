package dogu

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func NewDeploymentFetcher(client kubernetes.Interface) *deploymentFetcher {
	return &deploymentFetcher{client: client}
}

type deploymentFetcher struct {
	client kubernetes.Interface
}

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
	deploymentList, err := f.client.AppsV1().Deployments(namespace).List(ctx, options)
	if err != nil {
		return nil, err
	}

	return deploymentList.Items, nil
}
