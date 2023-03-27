package dogu

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	fakeappsv1 "k8s.io/client-go/kubernetes/typed/apps/v1/fake"
	clienttest "k8s.io/client-go/testing"
)

const testNamespace = "ecosystem"

func TestNewDeploymentFetcher(t *testing.T) {
	// given
	clientSet := fake.NewSimpleClientset()

	// when
	fetcher := NewDeploymentFetcher(clientSet)

	// then
	require.NotNil(t, fetcher)
	assert.Equal(t, clientSet, fetcher.clientSet)
}

func Test_deploymentFetcher_FetchAll(t *testing.T) {
	type args struct {
		ctx       context.Context
		namespace string
	}
	tests := []struct {
		name      string
		clientSet kubernetes.Interface
		args      args
		want      []appsv1.Deployment
		wantErr   func(t *testing.T, err error)
	}{
		{
			name:      "should fail to list deployments",
			clientSet: failingClientSet(),
			args:      args{ctx: context.TODO(), namespace: testNamespace},
			want:      nil,
			wantErr: func(t *testing.T, err error) {
				require.Error(t, err)
				assert.ErrorIs(t, err, assert.AnError)
				assert.ErrorContains(t, err, "could not list deployments with selector 'dogu.name'")
			},
		},
		{
			name: "should find two deployments",
			clientSet: fake.NewSimpleClientset(
				&appsv1.Deployment{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "matching-deployment",
						Namespace: testNamespace,
						Labels: map[string]string{
							"app":          "ces",
							"dogu.name":    "cas",
							"dogu.version": "1.2.3-1",
						},
					},
				},
				&appsv1.Deployment{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "dog",
						Namespace: testNamespace,
						Labels: map[string]string{
							"app":       "ces",
							"dog.name":  "Hasso",
							"dog.breed": "Schäferhund",
						},
					},
				},
				&appsv1.Deployment{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "matching-deployment2",
						Namespace: testNamespace,
						Labels: map[string]string{
							"dogu.name": "redmine",
						},
					},
				},
			),
			args: args{ctx: context.TODO(), namespace: testNamespace},
			want: []appsv1.Deployment{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "matching-deployment",
						Namespace: testNamespace,
						Labels: map[string]string{
							"app":          "ces",
							"dogu.name":    "cas",
							"dogu.version": "1.2.3-1",
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "matching-deployment2",
						Namespace: testNamespace,
						Labels: map[string]string{
							"dogu.name": "redmine",
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
			f := &deploymentFetcher{
				clientSet: tt.clientSet,
			}
			got, err := f.FetchAll(tt.args.ctx, tt.args.namespace)
			tt.wantErr(t, err)
			assert.Equalf(t, tt.want, got, "FetchAll(%v, %v)", tt.args.ctx, tt.args.namespace)
		})
	}
}

func failingClientSet() *fake.Clientset {
	clientSet := fake.NewSimpleClientset()
	clientSet.AppsV1().(*fakeappsv1.FakeAppsV1).PrependReactor("list", "deployments", func(action clienttest.Action) (handled bool, ret runtime.Object, err error) {
		return true, nil, assert.AnError
	})
	return clientSet
}
