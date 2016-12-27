package e2e

import (
	"errors"
	"fmt"

	"github.com/appscode/k8s-addons/pkg/testing"
	"github.com/appscode/searchlight/pkg/controller/host"
	"github.com/appscode/searchlight/test/plugin"
	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/apps"
	ext "k8s.io/kubernetes/pkg/apis/extensions"
	clientset "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
	"k8s.io/kubernetes/pkg/labels"
)

const (
	GOHOSTOS         string = "linux"
	GOHOSTARCH       string = "amd64"
	DefaultNamespace string = "default"
)

type DataConfig struct {
	ObjectType   string
	CheckCommand string
	Namespace    string
}

func fixNamespace(ns string) string {
	if ns == "" {
		return DefaultNamespace
	}
	return ns
}

func getHostName(objectType, objectName string, namespace ...string) string {
	object := objectName
	if objectType != "" {
		object = fmt.Sprintf("%s|%s", objectType, objectName)
	}

	if len(namespace) == 1 {
		object = fmt.Sprintf("%s@%s", object, namespace[0])
	} else {
		object = fmt.Sprintf("%s@default", object)
	}
	return object
}

func getClusterCheckData(kubeClient clientset.Interface, checkCommand, namespace string) (name string, count int32, err error) {
	var podList *kapi.PodList
	if podList, err = kubeClient.Core().Pods(fixNamespace(namespace)).List(
		kapi.ListOptions{LabelSelector: labels.Everything()}); err != nil {
		return
	}
	count = int32(len(podList.Items))
	name = getHostName("", checkCommand, namespace)
	return
}

func getKubernetesObjectData(kubeClient clientset.Interface, objectType, namespace string) (name string, count int32, err error) {
	switch objectType {
	case host.TypeReplicationcontrollers:
		replicationController := &kapi.ReplicationController{}
		replicationController.Namespace = namespace
		if err = testing.CreateKubernetesObject(kubeClient, replicationController); err != nil {
			return
		}
		name = getHostName(host.TypeReplicationcontrollers, replicationController.Name, replicationController.Namespace)
		count = replicationController.Spec.Replicas
	case host.TypeDaemonsets:
		daemonSet := &ext.DaemonSet{}
		daemonSet.Namespace = namespace
		if err = testing.CreateKubernetesObject(kubeClient, daemonSet); err != nil {
			return
		}

		if daemonSet, err = kubeClient.Extensions().
			DaemonSets(daemonSet.Namespace).Get(daemonSet.Name); err != nil {
			return
		}
		name = getHostName(host.TypeDaemonsets, daemonSet.Name, daemonSet.Namespace)
		count = daemonSet.Status.DesiredNumberScheduled
	case host.TypeStatefulSet:
		statefulSet := &apps.StatefulSet{}
		statefulSet.Namespace = namespace
		if err = testing.CreateKubernetesObject(kubeClient, statefulSet); err != nil {
			return
		}
		name = getHostName(host.TypeStatefulSet, statefulSet.Name, statefulSet.Namespace)
		count = statefulSet.Spec.Replicas
	case host.TypeReplicasets:
		replicaSet := &ext.ReplicaSet{}
		replicaSet.Namespace = namespace
		if err = testing.CreateKubernetesObject(kubeClient, replicaSet); err != nil {
			return
		}
		name = getHostName(host.TypeReplicasets, replicaSet.Name, replicaSet.Namespace)
		count = replicaSet.Spec.Replicas
	case host.TypeDeployments:
		deployment := &ext.Deployment{}
		deployment.Namespace = namespace
		if err = testing.CreateKubernetesObject(kubeClient, deployment); err != nil {
			return
		}
		name = getHostName(host.TypeDeployments, deployment.Name, deployment.Namespace)
		count = deployment.Spec.Replicas
	case host.TypePods:
		pod := &kapi.Pod{}
		pod.Namespace = namespace
		if err = testing.CreateKubernetesObject(kubeClient, pod); err != nil {
			return
		}
		name = getHostName(host.TypePods, pod.Name, pod.Namespace)

	case host.TypeServices:
		replicaSet := &ext.ReplicaSet{}
		replicaSet.Namespace = namespace
		if err = testing.CreateKubernetesObject(kubeClient, replicaSet); err != nil {
			return
		}

		service := &kapi.Service{
			ObjectMeta: kapi.ObjectMeta{
				Namespace: replicaSet.Namespace,
			},
			Spec: kapi.ServiceSpec{
				Selector: replicaSet.Spec.Selector.MatchLabels,
			},
		}
		if err = testing.CreateKubernetesObject(kubeClient, service); err != nil {
			return
		}
		name = getHostName(host.TypeServices, service.Name, service.Namespace)
		count = replicaSet.Spec.Replicas
	default:
		err = errors.New("Unknown objectType")
	}
	return
}

func GetTestData(kubeClient clientset.Interface, dataConfig *DataConfig) (name string, count int32) {
	var err error
	if dataConfig.ObjectType == host.TypeCluster {
		name, count, err = getClusterCheckData(kubeClient, dataConfig.CheckCommand, dataConfig.Namespace)
		plugin.Fatalln(err)
	} else {
		name, count, err = getKubernetesObjectData(kubeClient, dataConfig.ObjectType, dataConfig.Namespace)
		plugin.Fatalln(err)
	}
	return
}

func CreateNewNamespace(kubeClient clientset.Interface, name string) {
	ns := &kapi.Namespace{
		ObjectMeta: kapi.ObjectMeta{
			Name: name,
		},
	}
	_, err := kubeClient.Core().Namespaces().Create(ns)
	plugin.Fatalln(err)
}

func deleteNewNamespace(kubeClient clientset.Interface, name string) {
	err := kubeClient.Core().Namespaces().Delete(name, nil)
	plugin.Fatalln(err)
}
