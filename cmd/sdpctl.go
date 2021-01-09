package main

import (
	"github.com/spf13/cobra"
	"os"
	"sdpctl/pkg/cmd/get"
	"sdpctl/pkg/cmd/ops"
	"sdpctl/pkg/cmd/shell"
	"sdpctl/pkg/util"
)

func main() {
	util.InitLogger()
	var rootCmd = &cobra.Command{
		Use:   "sdpctl",
		Short: "ND Kubernetes 运维工具",
		Run:   runHelp,
	}

	rootCmd.AddCommand(get.NewCmdGet())
	rootCmd.AddCommand(ops.NewCmdOps())
	rootCmd.AddCommand(shell.NewCmdSh())
	if err := execute(rootCmd); err != nil {
		os.Exit(1)
	}
}

func execute(cmd *cobra.Command) error {
	err := cmd.Execute()
	return err
}

func runHelp(cmd *cobra.Command, args []string) {
	cmd.Help()
}
