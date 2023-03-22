package deployment

import (
	"context"
	"github.com/hashicorp/go-multierror"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type updater struct {
	clientSet kubernetes.Interface
}

func NewUpdater(clientSet kubernetes.Interface) *updater {
	return &updater{clientSet: clientSet}
}

func (u *updater) Update(ctx context.Context, namespace string, deployments []appsv1.Deployment) error {
	var multiErr error
	for _, deploy := range deployments {
		_, err := u.clientSet.AppsV1().Deployments(namespace).Update(ctx, &deploy, metav1.UpdateOptions{})
		if err != nil {
			multiErr = multierror.Append(multiErr, err)
		}
	}
	if multiErr != nil {
		return multiErr
	}

	return nil
}