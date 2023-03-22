package hosts

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

const testNamespace = "ecosystem"

var hostAliases = []corev1.HostAlias{{
	IP:        "1.2.3.4",
	Hostnames: []string{"www.example.com"},
}}

var doguDeployments = []appsv1.Deployment{{
	TypeMeta: metav1.TypeMeta{
		Kind:       "Deployment",
		APIVersion: "apps/v1",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "cas",
		Namespace: testNamespace,
		Labels: map[string]string{
			"app":       "ces",
			"dogu.name": "cas",
		},
	},
}}

func Test_hostAliasUpdater_UpdateHosts(t *testing.T) {
	type fields struct {
		generator hostAliasGenerator
		fetcher   doguDeploymentFetcher
		patcher   hostAliasPatcher
		updater   deploymentUpdater
	}
	type args struct {
		ctx       context.Context
		namespace string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr func(t *testing.T, err error)
	}{
		{
			name:   "should fail to generate host aliases",
			fields: fields{generator: failingHostAliasGenerator(t)},
			args: args{
				ctx:       context.TODO(),
				namespace: testNamespace,
			},
			wantErr: func(t *testing.T, err error) {
				require.Error(t, err)
				assert.ErrorIs(t, err, assert.AnError)
				assert.ErrorContains(t, err, "failed to generate host aliases")
			},
		},
		{
			name: "should fail to fetch dogu deployments",
			fields: fields{
				generator: succeedingHostAliasGenerator(t),
				fetcher:   failingDoguDeploymentFetcher(t),
			},
			args: args{
				ctx:       context.TODO(),
				namespace: testNamespace,
			},
			wantErr: func(t *testing.T, err error) {
				require.Error(t, err)
				assert.ErrorIs(t, err, assert.AnError)
				assert.ErrorContains(t, err, "failed to fetch dogu deployments")
			},
		},
		{
			name: "should fail to update dogu deployments",
			fields: fields{
				generator: succeedingHostAliasGenerator(t),
				fetcher:   succeedingDoguDeploymentFetcher(t),
				patcher:   createHostAliasPatcher(t),
				updater:   failingDeploymentUpdater(t),
			},
			args: args{
				ctx:       context.TODO(),
				namespace: testNamespace,
			},
			wantErr: func(t *testing.T, err error) {
				require.Error(t, err)
				assert.ErrorIs(t, err, assert.AnError)
				assert.ErrorContains(t, err, "failed to update host-aliases of dogu deployments in cluster")
			},
		},
		{
			name: "should succeed",
			fields: fields{
				generator: succeedingHostAliasGenerator(t),
				fetcher:   succeedingDoguDeploymentFetcher(t),
				patcher:   createHostAliasPatcher(t),
				updater:   succeedingDeploymentUpdater(t),
			},
			args: args{
				ctx:       context.TODO(),
				namespace: testNamespace,
			},
			wantErr: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hau := &hostAliasUpdater{
				generator: tt.fields.generator,
				fetcher:   tt.fields.fetcher,
				patcher:   tt.fields.patcher,
				updater:   tt.fields.updater,
			}
			err := hau.UpdateHosts(tt.args.ctx, tt.args.namespace)
			tt.wantErr(t, err)
		})
	}
}

func failingHostAliasGenerator(t *testing.T) hostAliasGenerator {
	t.Helper()
	generator := newMockHostAliasGenerator(t)
	generator.EXPECT().Generate().Return(nil, assert.AnError)
	return generator
}

func succeedingHostAliasGenerator(t *testing.T) hostAliasGenerator {
	t.Helper()
	generator := newMockHostAliasGenerator(t)
	generator.EXPECT().Generate().Return(hostAliases, nil)
	return generator
}

func failingDoguDeploymentFetcher(t *testing.T) doguDeploymentFetcher {
	t.Helper()
	fetcher := newMockDoguDeploymentFetcher(t)
	fetcher.EXPECT().FetchAll(context.TODO(), testNamespace).Return(nil, assert.AnError)
	return fetcher
}

func succeedingDoguDeploymentFetcher(t *testing.T) doguDeploymentFetcher {
	t.Helper()
	fetcher := newMockDoguDeploymentFetcher(t)
	fetcher.EXPECT().FetchAll(context.TODO(), testNamespace).Return(doguDeployments, nil)
	return fetcher
}

func createHostAliasPatcher(t *testing.T) hostAliasPatcher {
	t.Helper()
	patcher := newMockHostAliasPatcher(t)
	patcher.EXPECT().Patch(doguDeployments, hostAliases)
	return patcher
}

func failingDeploymentUpdater(t *testing.T) deploymentUpdater {
	t.Helper()
	updater := newMockDeploymentUpdater(t)
	updater.EXPECT().Update(context.TODO(), testNamespace, doguDeployments).Return(assert.AnError)
	return updater
}

func succeedingDeploymentUpdater(t *testing.T) deploymentUpdater {
	t.Helper()
	updater := newMockDeploymentUpdater(t)
	updater.EXPECT().Update(context.TODO(), testNamespace, doguDeployments).Return(nil)
	return updater
}

func TestNewHostAliasUpdater(t *testing.T) {
	// given
	clientSet := fake.NewSimpleClientset()
	cesReg := newMockCesRegistry(t)
	cesReg.EXPECT().GlobalConfig().Return(nil)

	// when
	updater := NewHostAliasUpdater(clientSet, cesReg)

	// then
	require.NotNil(t, updater)
}
