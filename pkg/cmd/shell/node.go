package shell

import (
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	k8stools "sdpctl/pkg/util/kube"
	"strings"
)

func nodeShell(cmd *cobra.Command, args []string) {

	kubeClientSet, kubeClientConfig := k8stools.KubeClientAndConfig(kubeConfigStr)

	cmdStr := strings.Join(args, " ")
	// 返回所有需要运行运行的Node列表
	pods := getShellPods(kubeClientSet, targetNode, toolName)
	for _, pod := range pods {
		var stdOut, stdErr bytes.Buffer
		shExecOps := k8stools.ExecOptions{
			Command:       cmdStr,
			ContainerName: "",
			In:            nil,
			Out:           &stdOut,
			Err:           &stdErr,
			Istty:         false,
			TimeOut:       httpTimeOutInSec,
			Pod:           &pod,
		}
		if err := k8stools.ExecCmd(kubeClientSet, kubeClientConfig, &shExecOps); err != nil {
			stdErr.Write([]byte(err.Error()))
		}
		printOutput(&pod, stdOut, stdErr)
	}

}

//打印输出
func printOutput(pod *v1.Pod, stdOut, stdErr bytes.Buffer) {

	switch format {
	case "raw":
		color.Blue("------------------------------> No. Shell on node: %s <------------------------------", pod.Status.HostIP)
		fmt.Printf(stdOut.String())
		color.HiRed(stdErr.String())
		break
	case "prefix":
		color.Blue("------------------------------> No. Shell on node: %s <------------------------------", pod.Status.HostIP)
		printWithPrefix(color.BlueString("[%s]", pod.Status.HostIP), stdOut.String())
		printWithPrefix(color.RedString("[%s]", pod.Status.HostIP), stdErr.String())
		break
	default:
		logrus.Error("不支持该格式输出")

	}

}

// 返回目标节点（Node List）的shell pod列表
func getShellPods(kubeClientSet *kubernetes.Clientset, shellNodeName, toolName string) []v1.Pod {

	shellPods, _ := k8stools.GetPodList(kubeClientSet, toolName, "name="+toolName)

	if strings.EqualFold(shellNodeName, "") {
		// 指定Node执行shell
		return shellPods.Items
	} else {
		// 所有Node执行Shell
		for _, pod := range shellPods.Items {
			if strings.EqualFold(pod.Status.HostIP, shellNodeName) {
				return []v1.Pod{pod}
			}

		}
	}

	return []v1.Pod{}
}
