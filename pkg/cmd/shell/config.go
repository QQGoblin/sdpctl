package shell

import (
	"github.com/spf13/pflag"
	"os"
	"path/filepath"
)

var (
	httpTimeOutInSec int
	currentThreadNum int
	targetNode       string
	shellMode        string
	toolName         string
	format           string
	kubeConfigStr    string
)

func AddShellFlags(flags *pflag.FlagSet) {
	flags.IntVar(&httpTimeOutInSec, "kubelet-timeout", 30, "连接Kubelet超时时间。")
	flags.StringVarP(&targetNode, "node", "n", "", "在指定宿主机节点执行操作。")
	flags.StringVar(&shellMode, "shell-mode", "k8s-node", "执行shell的模式：k8s-node，container-net")
	flags.StringVar(&toolName, "shell-tool-name", "node-shell", "Shell客户端工具名称。")
	flags.StringVarP(&format, "format", "f", "prefix", "输出格式：prefix 或者 raw")
	flags.StringVar(&kubeConfigStr, "kubeconfig", filepath.Join(os.Getenv("HOME"), ".kube", "config"), "Kubernete集群配置文件。")
}
