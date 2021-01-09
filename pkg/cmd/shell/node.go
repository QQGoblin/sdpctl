package shell

import (
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	k8stools "sdpctl/pkg/util/kube"
	"strings"
	"sync"
)

func nodeShell(cmd *cobra.Command, args []string) {

	kubeClientSet, kubeClientConfig := k8stools.KubeClientAndConfig(kubeConfigStr)

	cmdStr := strings.Join(args, " ")
	// 返回所有需要运行运行的Node列表
	pods := getShellPodDict(kubeClientSet, targetNode, toolName)
	threadNum := 0
	total := len(pods)
	outPuts := make([]OutPut, len(pods))
	var wg sync.WaitGroup
	for i, pod := range pods {
		outPut := OutPut{
			NodeName: pod.Status.HostIP,
			StdOut:   bytes.NewBufferString(""),
			StdErr:   bytes.NewBufferString(""),
		}

		outPuts[i] = outPut
		wg.Add(1)

		if pod.Status.Phase != v1.PodRunning {
			outPuts[i].StdOut.WriteString(" Shell容器异常")
		} else {
			shExecOps := k8stools.ExecOptions{
				Command:       cmdStr,
				ContainerName: "",
				In:            nil,
				Out:           outPut.StdOut,
				Err:           outPut.StdErr,
				Istty:         false,
				TimeOut:       httpTimeOutInSec,
			}
			go k8stools.ExecCmd(kubeClientSet, kubeClientConfig, &pod, shExecOps)
			threadNum += 1
			if threadNum%currentThreadNum == 0 || total == i+1 {
				wg.Wait()
			}
		}

	}
	printOutput(outPuts)
}

func printOutput(outPuts []OutPut) {
	for i, output := range outPuts {
		switch format {
		case "raw":
			color.Blue("------------------------------> No.%d  Shell on node: %s <------------------------------", i, output.NodeName)
			fmt.Printf(output.StdOut.String())
			color.HiRed(output.StdErr.String())
			break
		case "prefix":
			color.Blue("------------------------------> No.%d  Shell on node: %s <------------------------------", i, output.NodeName)
			prefixStr := color.BlueString("[%s]", output.NodeName)
			for {
				line, err := output.StdOut.ReadString('\n')
				if err != nil || io.EOF == err {
					break
				}
				fmt.Printf("%s %s", prefixStr, line)
			}
			for {
				line, err := output.StdErr.ReadString('\n')
				if err != nil || io.EOF == err {
					break
				}
				fmt.Printf("%s %s", prefixStr, color.RedString(line))
			}
			break
		default:
			logrus.Error("不支持该格式输出")

		}

	}
}

// 返回目标节点（Node List）的shell pod列表
func getShellPodDict(kubeClientSet *kubernetes.Clientset, shellNodeName, toolName string) []v1.Pod {

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
