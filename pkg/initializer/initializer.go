package initializer

import (
	"os"

	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
)

const namespaceEnvName = "NAMESPACE"

// Initializer is used for populating this program with configuration values.
type Initializer interface {
	// GetNamespace retrieves the namespace this program should work in.
	GetNamespace() string
	// CreateClientSet creates a client set from a kubernetes rest config.
	CreateClientSet() (kubernetes.Interface, error)
}

type defaultInitializer struct {
}

var New = func() Initializer {
	return &defaultInitializer{}
}

// GetNamespace retrieves the namespace this program should work in from the NAMESPACE environment variable.
// If the NAMESPACE var is not set or contains an empty string, the 'default' namespace is returned instead.
func (i *defaultInitializer) GetNamespace() string {
	env, present := os.LookupEnv(namespaceEnvName)
	if present && env != "" {
		return env
	}

	return "default"
}

// CreateClientSet creates a client set from a kubernetes rest config.
func (i *defaultInitializer) CreateClientSet() (kubernetes.Interface, error) {
	restConfig := ctrl.GetConfigOrDie()
	clientSet, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}

	return clientSet, nil
}
