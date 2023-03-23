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
	t.Run("should fail to generate host aliases", func(t *testing.T) {
		// given
		generator := failingHostAliasGenerator(t)
		ctx := context.TODO()
		sut := &hostAliasUpdater{generator: generator}

		// when
		err := sut.UpdateHosts(ctx, testNamespace)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to generate host aliases")
	})
	t.Run("should fail to fetch dogu deployments", func(t *testing.T) {
		// given
		generator := succeedingHostAliasGenerator(t)
		fetcher := failingDoguDeploymentFetcher(t)
		ctx := context.TODO()
		sut := &hostAliasUpdater{generator: generator, fetcher: fetcher}

		// when
		err := sut.UpdateHosts(ctx, testNamespace)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to fetch dogu deployments")
	})
	t.Run("should fail to update dogu deployments", func(t *testing.T) {
		// given
		generator := succeedingHostAliasGenerator(t)
		fetcher := succeedingDoguDeploymentFetcherOnRollback(t)
		patcher := createHostAliasPatcher(t)
		updater := failingDeploymentUpdater(t)
		ctx := context.TODO()
		sut := &hostAliasUpdater{
			generator: generator,
			fetcher:   fetcher,
			patcher:   patcher,
			updater:   updater,
		}

		// when
		err := sut.UpdateHosts(ctx, testNamespace)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to update host-aliases of dogu deployments in cluster")
	})
	t.Run("should fail to fetch dogu deployments on rollback", func(t *testing.T) {
		// given
		generator := succeedingHostAliasGenerator(t)
		fetcher := failingDoguDeploymentFetcherOnRollback(t)
		patcher := createHostAliasPatcher(t)
		updater := failingDeploymentUpdaterCallOnce(t)
		ctx := context.TODO()
		sut := &hostAliasUpdater{
			generator: generator,
			fetcher:   fetcher,
			patcher:   patcher,
			updater:   updater,
		}

		// when
		err := sut.UpdateHosts(ctx, testNamespace)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to update host-aliases of dogu deployments in cluster")
		assert.ErrorContains(t, err, "failed to fetch dogu deployments on rollback")
	})
	t.Run("should fail to update dogu deployments on rollback", func(t *testing.T) {
		// given
		generator := succeedingHostAliasGenerator(t)
		fetcher := succeedingDoguDeploymentFetcherOnRollback(t)
		patcher := createHostAliasPatcher(t)
		updater := failingDeploymentUpdaterOnRollback(t)
		ctx := context.TODO()
		sut := &hostAliasUpdater{
			generator: generator,
			fetcher:   fetcher,
			patcher:   patcher,
			updater:   updater,
		}

		// when
		err := sut.UpdateHosts(ctx, testNamespace)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to update host-aliases of dogu deployments in cluster")
		assert.ErrorContains(t, err, "failed to rollback dogu deployments")
	})
	t.Run("should fail to update dogu deployments on rollback", func(t *testing.T) {
		// given
		generator := succeedingHostAliasGenerator(t)
		fetcher := succeedingDoguDeploymentFetcher(t)
		patcher := createHostAliasPatcher(t)
		updater := succeedingDeploymentUpdater(t)
		ctx := context.TODO()
		sut := &hostAliasUpdater{
			generator: generator,
			fetcher:   fetcher,
			patcher:   patcher,
			updater:   updater,
		}

		// when
		err := sut.UpdateHosts(ctx, testNamespace)

		// then
		require.NoError(t, err)
	})
}

func failingHostAliasGenerator(t *testing.T) hostAliasGenerator {
	t.Helper()
	generator := newMockHostAliasGenerator(t)
	generator.EXPECT().Generate().Return(nil, assert.AnError).Once()
	return generator
}

func succeedingHostAliasGenerator(t *testing.T) hostAliasGenerator {
	t.Helper()
	generator := newMockHostAliasGenerator(t)
	generator.EXPECT().Generate().Return(hostAliases, nil).Once()
	return generator
}

func failingDoguDeploymentFetcher(t *testing.T) doguDeploymentFetcher {
	t.Helper()
	fetcher := newMockDoguDeploymentFetcher(t)
	fetcher.EXPECT().FetchAll(context.TODO(), testNamespace).Return(nil, assert.AnError).Once()
	return fetcher
}

func succeedingDoguDeploymentFetcher(t *testing.T) doguDeploymentFetcher {
	t.Helper()
	fetcher := newMockDoguDeploymentFetcher(t)
	fetcher.EXPECT().FetchAll(context.TODO(), testNamespace).Return(doguDeployments, nil).Once()
	return fetcher
}

func succeedingDoguDeploymentFetcherOnRollback(t *testing.T) doguDeploymentFetcher {
	t.Helper()
	fetcher := newMockDoguDeploymentFetcher(t)
	fetcher.EXPECT().FetchAll(context.TODO(), testNamespace).Return(doguDeployments, nil).Once()
	fetcher.EXPECT().FetchAll(context.TODO(), testNamespace).Return(doguDeployments, nil).Once()
	return fetcher
}

func failingDoguDeploymentFetcherOnRollback(t *testing.T) doguDeploymentFetcher {
	t.Helper()
	fetcher := newMockDoguDeploymentFetcher(t)
	fetcher.EXPECT().FetchAll(context.TODO(), testNamespace).Return(doguDeployments, nil).Once()
	fetcher.EXPECT().FetchAll(context.TODO(), testNamespace).Return(nil, assert.AnError).Once()
	return fetcher
}

func createHostAliasPatcher(t *testing.T) hostAliasPatcher {
	t.Helper()
	patcher := newMockHostAliasPatcher(t)
	patcher.EXPECT().Patch(doguDeployments, hostAliases).Return().Once()
	return patcher
}

func failingDeploymentUpdater(t *testing.T) deploymentUpdater {
	t.Helper()
	updater := newMockDeploymentUpdater(t)
	updater.EXPECT().Update(context.TODO(), testNamespace, doguDeployments).Return(assert.AnError).Once()
	updater.EXPECT().Update(context.TODO(), testNamespace, doguDeployments).Return(nil).Once()
	return updater
}

func failingDeploymentUpdaterCallOnce(t *testing.T) deploymentUpdater {
	t.Helper()
	updater := newMockDeploymentUpdater(t)
	updater.EXPECT().Update(context.TODO(), testNamespace, doguDeployments).Return(assert.AnError).Once()
	return updater
}

func failingDeploymentUpdaterOnRollback(t *testing.T) deploymentUpdater {
	t.Helper()
	updater := newMockDeploymentUpdater(t)
	updater.EXPECT().Update(context.TODO(), testNamespace, doguDeployments).Return(assert.AnError).Once()
	updater.EXPECT().Update(context.TODO(), testNamespace, doguDeployments).Return(assert.AnError).Once()
	return updater
}

func succeedingDeploymentUpdater(t *testing.T) deploymentUpdater {
	t.Helper()
	updater := newMockDeploymentUpdater(t)
	updater.EXPECT().Update(context.TODO(), testNamespace, doguDeployments).Return(nil).Once()
	return updater
}

func TestNewHostAliasUpdater(t *testing.T) {
	// given
	clientSet := fake.NewSimpleClientset()
	cesReg := newMockCesRegistry(t)
	cesReg.EXPECT().GlobalConfig().Return(nil).Once()

	// when
	updater := NewHostAliasUpdater(clientSet, cesReg)

	// then
	require.NotNil(t, updater)
}
