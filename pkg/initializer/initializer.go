package initializer

import (
	"fmt"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/cesapp-lib/registry"
	"k8s.io/client-go/kubernetes"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
)

type initializer struct {
}

var New = func() *initializer {
	return &initializer{}
}

func (i *initializer) GetNamespace() string {
	env, present := os.LookupEnv("NAMESPACE")
	if present {
		return env
	}

	return "default"
}

func (i *initializer) CreateClientSet() (kubernetes.Interface, error) {
	restConfig := ctrl.GetConfigOrDie()
	clientSet, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}

	return clientSet, nil
}

func (i *initializer) CreateCesRegistry() (registry.Registry, error) {
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
