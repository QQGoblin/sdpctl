package kubernetes

import (
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
	"strings"
)

// 获取Kube Client 客户端
func KubeClientAndConfig(configStr string) (*kubernetes.Clientset, *restclient.Config) {
	config, err := KubeConfg(configStr)
	if err != nil {
		logrus.Error(err)
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		logrus.Error(err)
		panic(err.Error())
	}
	return clientset, config
}

func KubeMetricsAndConfig(configStr string) (*metrics.Clientset, *restclient.Config) {
	config, err := KubeConfg(configStr)
	if err != nil {
		logrus.Error(err)
		panic(err)
	}
	metricsClient, err := metrics.NewForConfig(config)
	if err != nil {
		logrus.Error(err)
		panic(err.Error())
	}
	return metricsClient, config
}

func KubeConfg(configStr string) (config *restclient.Config, err error) {
	if strings.EqualFold(configStr, "") {
		config, err = restclient.InClusterConfig()
	} else {
		config, err = clientcmd.BuildConfigFromFlags("", configStr)
	}
	return config, err
}
