package util

import (
	"errors"

	"github.com/appscode/searchlight/pkg/client/k8s"
	"github.com/appscode/searchlight/pkg/controller/host"
	"k8s.io/kubernetes/pkg/labels"
)

func GetLabels(client *k8s.KubeClient, namespace, objectType, objectName string) (labels.Selector, error) {
	var labelMap map[string]string
	switch objectType {
	case host.TypeServices:
		service, err := client.Client.Core().Services(namespace).Get(objectName)
		if err != nil {
			return nil, err
		}
		labelMap = service.Spec.Selector
	case host.TypeReplicationcontrollers:
		rc, err := client.Client.Core().ReplicationControllers(namespace).Get(objectName)
		if err != nil {
			return nil, err
		}
		labelMap = rc.Spec.Selector
	case host.TypeDaemonsets:
		daemonSet, err := client.Client.Extensions().DaemonSets(namespace).Get(objectName)
		if err != nil {
			return nil, err
		}
		labelMap = daemonSet.Spec.Selector.MatchLabels
	case host.TypeReplicasets:
		replicaSet, err := client.Client.Extensions().ReplicaSets(namespace).Get(objectName)
		if err != nil {
			return nil, err
		}
		labelMap = replicaSet.Spec.Selector.MatchLabels
	case host.TypeStatefulSet:
		stateFulSet, err := client.Client.Apps().StatefulSets(namespace).Get(objectName)
		if err != nil {
			return nil, err
		}
		labelMap = stateFulSet.Spec.Selector.MatchLabels
	case host.TypeDeployments:
		deployment, err := client.Client.Extensions().Deployments(namespace).Get(objectName)
		if err != nil {
			return nil, err
		}
		labelMap = deployment.Spec.Selector.MatchLabels
	default:
		return nil, errors.New("Invalid kubernetes object type")
	}
	return labels.SelectorFromSet(labelMap), nil
}
