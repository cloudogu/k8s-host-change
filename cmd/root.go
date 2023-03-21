package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func RootCmd() *cobra.Command {
	return &cobra.Command{
		Use:          "k8s-host-change",
		Short:        "Sync hosts-specific changes",
		Long:         "Sync etcd values of FQDN, internal IP and additional hosts with the K8s Cloudogu EcoSystem.",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: add payload
			return nil
		},
	}
}

func InitAndExecute() {
	if err := RootCmd().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
