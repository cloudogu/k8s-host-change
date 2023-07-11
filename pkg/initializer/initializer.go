package initializer

import (
	"fmt"
	"os"

	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/cesapp-lib/registry"
)

const namespaceEnvName = "NAMESPACE"

type Initializer interface {
	GetNamespace() string
	CreateClientSet() (kubernetes.Interface, error)
	CreateCesRegistry() (registry.Registry, error)
}

type defaultInitializer struct {
}

var New = func() Initializer {
	return &defaultInitializer{}
}

func (i *defaultInitializer) GetNamespace() string {
	env, present := os.LookupEnv(namespaceEnvName)
	if present && env != "" {
		return env
	}

	return "default"
}

func (i *defaultInitializer) CreateClientSet() (kubernetes.Interface, error) {
	restConfig := ctrl.GetConfigOrDie()
	clientSet, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}

	return clientSet, nil
}

func (i *defaultInitializer) CreateCesRegistry() (registry.Registry, error) {
	namespace := i.GetNamespace()
	cesReg, err := registry.New(core.Registry{
		Type:      "etcd",
		Endpoints: []string{fmt.Sprintf("http://etcd.%s.svc.cluster.local:4001", namespace)},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create CES registry: %w", err)
	}

	return cesReg, nil
}
