package deployment

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/retry"
)

type updater struct {
	clientSet kubernetes.Interface
}

// NewUpdater creates a new instance of updater.
func NewUpdater(clientSet kubernetes.Interface) *updater {
	return &updater{clientSet: clientSet}
}

// UpdateHostAliases replaces the host aliases in the given deployments.
// Every deployment will be fetched again from the api with a retry mechanism to prevent
// conflict api errors.
func (u *updater) UpdateHostAliases(ctx context.Context, namespace string, deployments []appsv1.Deployment, aliases []corev1.HostAlias) error {
	var multiErr error
	for _, deploy := range deployments {
		err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
			deployment, err := u.clientSet.AppsV1().Deployments(namespace).Get(ctx, deploy.Name, metav1.GetOptions{})
			if err != nil {
				return fmt.Errorf("failed to get deployment '%s': %w", deploy.Name, err)
			}
			deployment.Spec.Template.Spec.HostAliases = aliases

			_, err = u.clientSet.AppsV1().Deployments(namespace).Update(ctx, deployment, metav1.UpdateOptions{})
			if err != nil {
				return fmt.Errorf("failed to update deployment '%s': %w", deploy.Name, err)
			}

			return nil
		})

		if err != nil {
			multiErr = multierror.Append(multiErr, err)
		}
	}
	if multiErr != nil {
		return multiErr
	}

	return nil
}
