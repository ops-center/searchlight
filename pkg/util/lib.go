package util

import (
	"errors"

	"github.com/appscode/searchlight/pkg/config"
	"k8s.io/kubernetes/pkg/labels"
	"k8s.io/kubernetes/pkg/selection"
	"k8s.io/kubernetes/pkg/util/sets"
)

func GetLabels(client *config.KubeClient, namespace, objectType, objectName string) (labels.Selector, error) {
	label := labels.NewSelector()
	labelsMap := make(map[string]string, 0)
	if objectType == config.TypeServices {
		service, err := client.Services(namespace).Get(objectName)
		if err != nil {
			return nil, err
		}
		labelsMap = service.Spec.Selector

	} else if objectType == config.TypeReplicationcontrollers {
		rc, err := client.ReplicationControllers(namespace).Get(objectName)
		if err != nil {
			return nil, err
		}
		labelsMap = rc.Spec.Selector
	} else if objectType == config.TypeDaemonsets {
		daemonSet, err := client.DaemonSets(namespace).Get(objectName)
		if err != nil {
			return nil, err
		}
		labelsMap = daemonSet.Spec.Selector.MatchLabels
	} else if objectType == config.TypeReplicasets {
		replicaSet, err := client.ReplicaSets(namespace).Get(objectName)
		if err != nil {
			return nil, err
		}
		labelsMap = replicaSet.Spec.Selector.MatchLabels
	} else if objectType == config.TypePetsets {
		petSet, err := client.PetSets(namespace).Get(objectName)
		if err != nil {
			return nil, err
		}
		labelsMap = petSet.Spec.Selector.MatchLabels
	} else if objectType == config.TypeDeployments {
		deployment, err := client.Deployments(namespace).Get(objectName)
		if err != nil {
			return nil, err
		}
		labelsMap = deployment.Spec.Selector.MatchLabels
	} else {
		return label, errors.New("Invalid kubernetes object type")
	}

	for key, value := range labelsMap {
		s := sets.NewString(value)
		ls, err := labels.NewRequirement(key, selection.Equals, s)
		if err != nil {
			return nil, err
		}
		label = label.Add(*ls)
	}

	return label, nil
}
