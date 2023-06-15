package cmd

import (
	"fmt"
	"os"

	"github.com/mehix/sec-checklist/cmd/api"
	"github.com/mehix/sec-checklist/cmd/client"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "secctrls",
	Short: "Manage security controls",
	Long:  "Register an application and view and manage its security controls",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

func Execute() {
	rootCmd.AddCommand(client.Command(), api.Command())
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
