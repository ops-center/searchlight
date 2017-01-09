package node_status

import (
	"fmt"
	"os"

	"github.com/appscode/searchlight/pkg/client/k8s"
	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/labels"
)

func GetKubernetesNodeName(kubeClient *k8s.KubeClient) string {
	nodeList, err := kubeClient.Client.Core().Nodes().List(
		kapi.ListOptions{
			LabelSelector: labels.Everything(),
		},
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(nodeList.Items) == 0 {
		fmt.Println("No node found")
		os.Exit(1)
	}
	return nodeList.Items[0].Name
}
