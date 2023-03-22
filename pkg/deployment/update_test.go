package deployment

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"testing"
)

const testNamespace = "ecosystem"

func TestNewUpdater(t *testing.T) {
	// given
	clientSet := fake.NewSimpleClientset()

	// when
	updater := NewUpdater(clientSet)

	// then
	require.NotNil(t, updater)
	assert.Equal(t, clientSet, updater.clientSet)
}

func Test_updater_Update(t *testing.T) {
	type args struct {
		ctx         context.Context
		namespace   string
		deployments []appsv1.Deployment
	}
	tests := []struct {
		name      string
		clientSet kubernetes.Interface
		args      args
		wantErr   func(t *testing.T, err error)
	}{
		{
			name:      "should fail once",
			clientSet: fake.NewSimpleClientset(),
			args: args{
				ctx:       context.TODO(),
				namespace: testNamespace,
				deployments: []appsv1.Deployment{{
					ObjectMeta: metav1.ObjectMeta{
						Name: "will-not-be-found",
					},
				}},
			},
			wantErr: func(t *testing.T, err error) {
				require.Error(t, err)
				assert.ErrorContains(t, err, "1 error occurred")
				assert.ErrorContains(t, err, "failed to update deployment 'will-not-be-found': deployments.apps \"will-not-be-found\" not found")
			},
		},
		{
			name: "should fail twice",
			clientSet: fake.NewSimpleClientset(&appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "will-be-found",
					Namespace: testNamespace,
				},
			}),
			args: args{
				ctx:       context.TODO(),
				namespace: testNamespace,
				deployments: []appsv1.Deployment{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name: "will-not-be-found",
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name: "will-be-found",
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name: "will-not-be-found-either",
						},
					},
				},
			},
			wantErr: func(t *testing.T, err error) {
				require.Error(t, err)
				assert.ErrorContains(t, err, "2 errors occurred")
				assert.ErrorContains(t, err, "failed to update deployment 'will-not-be-found': deployments.apps \"will-not-be-found\" not found")
				assert.ErrorContains(t, err, "failed to update deployment 'will-not-be-found-either': deployments.apps \"will-not-be-found-either\" not found")
			},
		},
		{
			name: "should succeed",
			clientSet: fake.NewSimpleClientset(
				&appsv1.Deployment{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "will-be-found",
						Namespace: testNamespace,
					},
				},
				&appsv1.Deployment{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "will-be-found-as-well",
						Namespace: testNamespace,
					},
				},
			),
			args: args{
				ctx:       context.TODO(),
				namespace: testNamespace,
				deployments: []appsv1.Deployment{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name: "will-be-found",
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name: "will-be-found-as-well",
						},
					},
				},
			},
			wantErr: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &updater{
				clientSet: tt.clientSet,
			}
			err := u.Update(tt.args.ctx, tt.args.namespace, tt.args.deployments)
			tt.wantErr(t, err)
		})
	}
}
