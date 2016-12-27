package node_status

import (
	"errors"

	config "github.com/appscode/searchlight/pkg/client/k8s"
	"github.com/appscode/searchlight/test/plugin"
	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/labels"
)

func GetKubernetesNodeName(kubeClient *config.KubeClient) string {
	nodeList, err := kubeClient.Client.Core().Nodes().List(
		kapi.ListOptions{
			LabelSelector: labels.Everything(),
		},
	)
	plugin.Fatalln(err)

	if len(nodeList.Items) == 0 {
		plugin.Fatalln(errors.New("No node found"))
	}
	return nodeList.Items[0].Name
}
