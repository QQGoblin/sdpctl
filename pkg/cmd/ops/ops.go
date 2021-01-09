package ops

import (
	"fmt"
	"github.com/spf13/cobra"
)

func NewCmdOps() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "ops",
		Short:                 "SDP相关依赖工具",
		DisableFlagsInUseLine: true,
		Run:                   runHelp,
	}
	return cmd
}

func runHelp(cmd *cobra.Command, args []string) {
	fmt.Println("#############################################################")
	fmt.Println("#                        部署辅助工具                        #")
	fmt.Println("#############################################################")
	fmt.Println("kubectl apply -f http://goblin.lqingcloud.cn:9000/cicd/yaml/node-shell.yml")
}
