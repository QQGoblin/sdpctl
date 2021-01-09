package get

import (
	mapset "github.com/deckarep/golang-set"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	metricsv1beta1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"
	"sdpctl/pkg/util"
	k8stools "sdpctl/pkg/util/kube"
	"sdpctl/pkg/util/table"
	"sort"
	"strconv"
	"strings"
)

var labelFilter = mapset.NewSet(
	"beta.kube.io/arch",
	"beta.kube.io/os",
	"kube.io/arch",
	"kube.io/hostname",
	"kube.io/os",
)

func NewCmdNode() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "node",
		Short:                 "打印节点信息",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			node(cmd, args)
		},
	}
	return cmd
}

func node(cmd *cobra.Command, args []string) {

	kubeClientSet, _ := k8stools.KubeClientAndConfig(kubeConfigStr)
	nodes, _ := kubeClientSet.CoreV1().Nodes().List(metav1.ListOptions{})
	nodeMetricsDict := nodeUsage()
	nodeBriefInfo(kubeClientSet, nodes, nodeMetricsDict)
}

func nodeUsage() map[string]metricsv1beta1.NodeMetrics {
	metricsCli, _ := k8stools.KubeMetricsAndConfig(kubeConfigStr)
	nodeMetricsDict := make(map[string]metricsv1beta1.NodeMetrics)
	if nodeMetricsList, err := metricsCli.MetricsV1beta1().NodeMetricses().List(metav1.ListOptions{}); err != nil {
		logrus.Error(err)
		return nil
	} else {
		for _, nodeMetrics := range nodeMetricsList.Items {
			nodeMetricsDict[nodeMetrics.Name] = nodeMetrics
		}
	}
	return nodeMetricsDict
}

func nodeBriefInfo(kubeClientSet *kubernetes.Clientset, nodes *v1.NodeList, nodeMetricsDict map[string]metricsv1beta1.NodeMetrics) {

	nodeInfoList := make([]NodeBriefInfo, len(nodes.Items))
	allPodDist, _ := k8stools.GetPodDict(kubeClientSet, "")
	for i, node := range nodes.Items {
		// 获取Role以及Label信息
		role := ""
		envLabel := make([]string, 0)
		typeLabel := make([]string, 0)
		commonLabel := make([]string, 0)
		for k, v := range node.Labels {
			if strings.EqualFold(v, "env") {
				envLabel = append(envLabel, k)
				continue
			}
			if strings.EqualFold(v, "type") {
				typeLabel = append(typeLabel, k)
				continue
			}
			if strings.HasPrefix(k, "node-role.kube.io") {
				role = strings.Split(k, "/")[1]
				continue
			}
			if !labelFilter.Contains(k) {
				commonLabel = append(commonLabel, k+"="+v)
				continue
			}
		}

		// 列出获取该节点的所有Pod
		podListOnNode := allPodDist[node.Name]

		// 获取节点的状态
		var unschedulable = ""
		if node.Spec.Unschedulable {
			unschedulable = "Y"
		}

		// 获取节点CPU/内存使用情况
		// 计算Pod申请的内存资源
		var memoryRequest int64 = 0
		var memoryLimits int64 = 0
		for _, pod := range podListOnNode {
			for _, c := range pod.Spec.Containers {
				memoryRequest += c.Resources.Requests.Memory().Value()
				memoryLimits += c.Resources.Limits.Memory().Value()
			}
		}

		nodeMetrics := nodeMetricsDict[node.Name]
		memoryUsageStr := util.FormatByte(nodeMetrics.Usage.Memory().Value())
		memoryCapacityStr := util.FormatByte(node.Status.Capacity.Memory().Value())
		memoryRequestStr := util.FormatByte(memoryRequest)
		memoryLimitsStr := util.FormatByte(memoryLimits)

		nodeInfo := NodeBriefInfo{
			Name:          node.Name,
			Role:          role,
			UnSche:        unschedulable,
			Env:           strings.Join(envLabel, ","),
			Type:          strings.Join(typeLabel, ","),
			Label:         strings.Join(commonLabel, ","),
			CPU:           node.Status.Capacity.Cpu().String(),
			Memory:        memoryCapacityStr,
			MemoryUsage:   memoryUsageStr + "(" + util.FormatPercentage(nodeMetrics.Usage.Memory().Value(), node.Status.Capacity.Memory().Value()) + ")",
			MemoryRequest: memoryRequestStr + "(" + util.FormatPercentage(memoryRequest, node.Status.Capacity.Memory().Value()) + ")",
			MemoryLimits:  memoryLimitsStr + "(" + util.FormatPercentage(memoryLimits, node.Status.Capacity.Memory().Value()) + ")",
			Pod:           strconv.Itoa(len(podListOnNode)) + "/" + node.Status.Capacity.Pods().String(),
		}

		if len(podListOnNode) > 80 {
			nodeInfo.Pod = color.HiYellowString(nodeInfo.Pod)
		} else if len(podListOnNode) < 30 {
			nodeInfo.Pod = color.HiGreenString(nodeInfo.Pod)
		}

		nodeInfoList[i] = nodeInfo
	}
	sort.Slice(nodeInfoList, func(i, j int) bool {
		return nodeInfoList[i].Role < nodeInfoList[j].Role
	})
	table.Output(nodeInfoList)
}
