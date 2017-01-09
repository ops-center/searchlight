package node_count

import (
	"fmt"
	"os"

	"github.com/appscode/searchlight/pkg/client/k8s"
	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/labels"
)

func GetKubernetesNodeCount(kubeClient *k8s.KubeClient) int {
	nodeList, err := kubeClient.Client.Core().
		Nodes().List(
		kapi.ListOptions{
			LabelSelector: labels.Everything(),
		},
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return len(nodeList.Items)
}
