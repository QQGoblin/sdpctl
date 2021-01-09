package kubernetes

import (
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"strings"
)

// 根据Node返回当前Pod运行字典
func GetPodDict(kubeClientSet *kubernetes.Clientset, lableSelector string) (podDist map[string][]v1.Pod, err error) {
	listOptions := metav1.ListOptions{}
	if !strings.EqualFold(lableSelector, "") {
		listOptions = metav1.ListOptions{
			TypeMeta:      metav1.TypeMeta{},
			LabelSelector: lableSelector,
		}
	}

	pods, err := kubeClientSet.CoreV1().Pods("").List(listOptions)
	if err != nil {
		logrus.Error(err.Error())
	}
	podDist = make(map[string][]v1.Pod)
	for _, pod := range pods.Items {
		key := pod.Spec.NodeName
		podListOnNode := podDist[key]
		if podListOnNode == nil {
			podListOnNode = make([]v1.Pod, 0)
		}
		podListOnNode = append(podListOnNode, pod)
		podDist[key] = podListOnNode
	}
	return
}

// 返回指定Pod列表
func GetPodList(kubeClientSet *kubernetes.Clientset, namespaceStr, lableSelector string) (pods *v1.PodList, err error) {

	listOptions := metav1.ListOptions{}
	if !strings.EqualFold(lableSelector, "") {
		listOptions = metav1.ListOptions{
			TypeMeta:      metav1.TypeMeta{},
			LabelSelector: lableSelector,
		}
	}
	pods, err = kubeClientSet.CoreV1().Pods(namespaceStr).List(listOptions)
	if err != nil {
		logrus.Error("获取Pod列表异常")
	}
	return
}
