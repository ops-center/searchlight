package pod_status

import (
	config "github.com/appscode/searchlight/pkg/client/k8s"
	"github.com/appscode/searchlight/pkg/controller/host"
	"github.com/appscode/searchlight/test/plugin"
	"github.com/appscode/searchlight/util"
	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/labels"
)

func GetStatusCodeForPodStatus(kubeClient *config.KubeClient, hostname string) int {
	objectType, objectName, namespace := plugin.GetKubeObjectInfo(hostname)

	var err error
	if objectType == host.TypePods {
		pod, err := kubeClient.Client.Core().Pods(namespace).Get(objectName)
		plugin.Fatalln(err)
		if !(pod.Status.Phase == kapi.PodSucceeded || pod.Status.Phase == kapi.PodRunning) {
			return plugin.CRITICAL
		}

	} else {
		labelSelector := labels.Everything()
		if objectType != "" {
			labelSelector, err = util.GetLabels(kubeClient, namespace, objectType, objectName)
			plugin.Fatalln(err)
		}
		var podList *kapi.PodList
		podList, err = kubeClient.Client.Core().Pods(namespace).List(kapi.ListOptions{LabelSelector: labelSelector})
		plugin.Fatalln(err)

		for _, pod := range podList.Items {
			if !(pod.Status.Phase == kapi.PodSucceeded || pod.Status.Phase == kapi.PodRunning) {
				return plugin.CRITICAL
			}
		}
	}
	return plugin.OK
}
