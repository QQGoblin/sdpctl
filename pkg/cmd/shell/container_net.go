package shell

import (
	"bytes"
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"io"
	pb "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
	"sdpctl/pkg/util"
	"strings"
)

const (
	DefaultEndpoint = "unix:///var/run/crio/crio.sock"
)

func dockerNet(cmd *cobra.Command, args []string) {

	if len(args) < 1 {
		logrus.Error("请输入要执行的命令")
	}
	cmdStr := strings.Join(args, " ")

	cli, err := getRuntimeClient(DefaultEndpoint)
	if err != nil {
		logrus.Errorf("创建CRI客户端失败，$s", err.Error())
		return
	}

	reqContainer := &pb.ListContainersRequest{}
	containerRes, _ := cli.ListContainers(context.Background(), reqContainer)

	for _, container := range containerRes.Containers {

		// 获取容器详情
		statusReq := &pb.ContainerStatusRequest{
			ContainerId: container.Id,
			Verbose:     true,
		}
		containerStatusRes, _ := cli.ContainerStatus(context.Background(), statusReq)

		podname := container.Labels["io.kube.pod.name"]
		pid := containerStatusRes.Info["pid"]
		color.HiBlue("---------------------> POD: %s, PID: %s <---------------------", podname, pid)
		if outStr, errStr, err := util.CmdOutErr("/usr/bin/nsenter", "-t", pid, "-n", "/bin/sh", "-c", cmdStr); err != nil {
			printWithPrefix(color.BlueString("[%s]", pid), outStr)
			printWithPrefix(color.RedString("[%s]", pid), errStr)

		}
	}

}

// 创建gRPC连接
func getRuntimeClient(endPoint string) (pb.RuntimeServiceClient, error) {
	conn, err := grpc.Dial(endPoint, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		errMsg := errors.Wrapf(err, "connect endpoint '%s', make sure you are running as root and the endpoint has been started", endPoint)
		logrus.Error(errMsg)
	} else {
		logrus.Debugf("connected successfully using endpoint: %s", endPoint)
	}
	runtimeClient := pb.NewRuntimeServiceClient(conn)
	return runtimeClient, nil
}

// f
func printWithPrefix(prefixStr, s string) {
	buf := bytes.NewBufferString(s)
	for {
		line, err := buf.ReadString('\n')
		if !strings.EqualFold(line, "") {
			fmt.Printf("%s %s\n", prefixStr, strings.TrimRight(line, "\n"))
		}
		if err != nil || io.EOF == err {

			break
		}

	}
}
