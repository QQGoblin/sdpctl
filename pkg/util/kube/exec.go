package kubernetes

import (
	"github.com/sirupsen/logrus"
	"io"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"strings"
	"time"
)

// 在Pod中执行指定命令

type ExecOptions struct {
	Command       string
	ContainerName string
	In            io.Reader
	Out           io.Writer
	Err           io.Writer
	Istty         bool
	TimeOut       int
}

type NodeShellError struct {
	ErrCode int
	ErrMsg  string
}

func (err *NodeShellError) Error() string {
	return err.ErrMsg
}

func ExecCmd(kubeClientSet *kubernetes.Clientset, kubeClientConfig *restclient.Config, pod *v1.Pod, execOptions ExecOptions) error {

	if pod.Status.Phase != v1.PodRunning {
		logrus.Println("Pod 没有就绪：", pod.Name, pod.Status.HostIP)
		err := NodeShellError{500, "Pod 没有就绪"}
		return &err
	}

	// 获取pod中的目标Container
	container := containerToExec(execOptions.ContainerName, pod)
	// 创建运行表达式
	podOptions := v1.PodExecOptions{
		Command:   strings.Fields(execOptions.Command),
		Container: container.Name,
		Stdin:     execOptions.In != nil,
		Stdout:    execOptions.Out != nil,
		Stderr:    execOptions.Err != nil,
		TTY:       execOptions.Istty,
	}

	// 创建客户端请求
	req := kubeClientSet.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(pod.Name).
		Namespace(pod.Namespace).
		SubResource("exec").
		Timeout(time.Duration(execOptions.TimeOut))

	req.VersionedParams(&podOptions, scheme.ParameterCodec)

	// 执行命令，并输出到标准输出
	streamOptions := remotecommand.StreamOptions{
		Stdin:  execOptions.In,
		Stdout: execOptions.Out,
		Stderr: execOptions.Err,
		Tty:    execOptions.Istty,
	}
	exec, err := remotecommand.NewSPDYExecutor(kubeClientConfig, "POST", req.URL())
	if err != nil {
		return err
	}
	return exec.Stream(streamOptions)
}

// 返回Pod内Container的名称
func containerToExec(container string, pod *v1.Pod) *v1.Container {
	if len(container) > 0 {
		for i := range pod.Spec.Containers {
			if pod.Spec.Containers[i].Name == container {
				return &pod.Spec.Containers[i]
			}
		}
		for i := range pod.Spec.InitContainers {
			if pod.Spec.InitContainers[i].Name == container {
				return &pod.Spec.InitContainers[i]
			}
		}
		logrus.Errorf("container not found (%s)", container)
		return nil
	}
	return &pod.Spec.Containers[0]
}
