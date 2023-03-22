package cmd

import (
	"context"
	"fmt"
	"github.com/cloudogu/k8s-host-change/pkg/initializer"
	"os"

	"github.com/spf13/cobra"

	"github.com/cloudogu/k8s-host-change/pkg/hosts"
)

func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "k8s-host-change",
		Short:        "Sync hosts-specific changes",
		Long:         "Sync etcd values of FQDN, internal IP and additional hosts with the K8s Cloudogu EcoSystem.",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			init := initializer.New()
			namespace := init.GetNamespace()

			clientSet, err := init.CreateClientSet()
			if err != nil {
				return err
			}

			cesReg, err := init.CreateCesRegistry()
			if err != nil {
				return err
			}

			updater, err := hosts.NewHostAliasUpdater(clientSet, cesReg)
			err = updater.UpdateHosts(ctx, namespace)
			if err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}

func InitAndExecute() {
	ctx := context.Background()
	if err := RootCmd().ExecuteContext(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
