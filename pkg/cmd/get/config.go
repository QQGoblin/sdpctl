package get

import (
	"github.com/spf13/pflag"
	"os"
	"path/filepath"
)

var (
	kubeConfigStr string
)

func addGetNodeFlags(flags *pflag.FlagSet) {
	flags.StringVar(&kubeConfigStr, "kubeconfig", filepath.Join(os.Getenv("HOME"), ".kube", "config"), "Kubernete集群配置文件。")

}
