package component_status

import (
	config "github.com/appscode/searchlight/pkg/client/k8s"
	"github.com/appscode/searchlight/test/plugin"
	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/labels"
)

func GetStatusCodeForComponentStatus(kubeClient *config.KubeClient) int {
	components, err := kubeClient.Client.Core().ComponentStatuses().
		List(kapi.ListOptions{LabelSelector: labels.Everything()})
	plugin.Fatalln(err)

	for _, component := range components.Items {
		for _, condition := range component.Conditions {
			if condition.Type == kapi.ComponentHealthy && condition.Status == kapi.ConditionFalse {
				return plugin.CRITICAL
			}
		}
	}
	return 0
}
