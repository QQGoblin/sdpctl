package get

import (
	"github.com/spf13/cobra"
)

func NewCmdGet() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "get",
		Short:                 "打印信息",
		DisableFlagsInUseLine: true,
	}
	addGetFlags(cmd.Flags())
	cmd.AddCommand(NewCmdNode())
	return cmd
}
