package main

import (
	"context"
	"os"

	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/cloudogu/k8s-host-change/pkg/alias"
	"github.com/cloudogu/k8s-host-change/pkg/hosts"
	"github.com/cloudogu/k8s-host-change/pkg/initializer"
	"github.com/cloudogu/k8s-host-change/pkg/logging"
	"github.com/cloudogu/k8s-registry-lib/registry"
)

var logger = ctrl.Log.WithName("k8s-host-change")

func init() {
	if err := logging.ConfigureLogger(); err != nil {
		panic(err.Error())
	}
}

func main() {
	err := run()
	if err != nil {
		handleError(err)
	}
}

func run() error {
	init := initializer.New()
	namespace := init.GetNamespace()

	clientSet, err := init.CreateClientSet()
	if err != nil {
		return err
	}

	globalConfig, err := registry.NewGlobalConfigReader(context.Background(), clientSet.CoreV1().ConfigMaps(namespace))
	if err != nil {
		return err
	}

	hostGenerator := alias.NewHostAliasGenerator(globalConfig)

	updater := hosts.NewHostAliasUpdater(clientSet, hostGenerator)
	err = updater.UpdateHosts(context.Background(), namespace)
	if err != nil {
		return err
	}

	return nil
}

func handleError(err error) {
	logger.Error(err, "exit k8s-host-change")
	os.Exit(1)
}
