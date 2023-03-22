package main

import (
	"github.com/cloudogu/k8s-host-change/cmd"
	"github.com/cloudogu/k8s-host-change/pkg/logging"
)

func init() {
	if err := logging.ConfigureLogger(); err != nil {
		panic(err.Error())
	}
}

func main() {
	cmd.InitAndExecute()
}
