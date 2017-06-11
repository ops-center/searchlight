package mini

import (
	"errors"
	"time"

	"github.com/appscode/searchlight/pkg/controller/host"
	"github.com/appscode/searchlight/pkg/testing"
	"github.com/appscode/searchlight/pkg/watcher"
	"github.com/appscode/searchlight/util"
	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/apis/extensions"
	"k8s.io/kubernetes/pkg/labels"
)

func checkReplicaSet(w *watcher.Watcher, replicaSet *extensions.ReplicaSet) (*extensions.ReplicaSet, error) {
	check := 0
	for {
		time.Sleep(time.Second * 30)
		nReplicaset, err := w.Storage.ReplicaSetStore.ReplicaSets(replicaSet.Namespace).Get(replicaSet.Name)
		if err != nil {
			return nil, err
		}
		if nReplicaset.Status.ReadyReplicas == nReplicaset.Status.Replicas {
			return nReplicaset, nil
		}

		if check > 6 {
			return nil, errors.New("Fail to create ReplicaSet")
		}
		check++
	}
}

func CreateReplicaSet(w *watcher.Watcher, namespace string) (*extensions.ReplicaSet, error) {
	replicaSet := &extensions.ReplicaSet{}
	replicaSet.Namespace = namespace
	if err := testing.CreateKubernetesObject(w.KubeClient, replicaSet); err != nil {
		return nil, err
	}

	return checkReplicaSet(w, replicaSet)
}

func ReCreateReplicaSet(w *watcher.Watcher, replicaSet *extensions.ReplicaSet) (*extensions.ReplicaSet, error) {
	newReplicaSet := &extensions.ReplicaSet{
		ObjectMeta: kapi.ObjectMeta{
			Name:      replicaSet.Name,
			Namespace: replicaSet.Namespace,
		},
		Spec: extensions.ReplicaSetSpec{
			Replicas: replicaSet.Spec.Replicas,
		},
	}
	if err := testing.CreateKubernetesObject(w.KubeClient, newReplicaSet); err != nil {
		return nil, err
	}

	return checkReplicaSet(w, newReplicaSet)
}

func GetLastReplica(w *watcher.Watcher, replicaSet *extensions.ReplicaSet) (*kapi.Pod, error) {
	podList, err := w.Storage.PodStore.List(labels.Set(replicaSet.Spec.Selector.MatchLabels).AsSelector())
	if err != nil {
		return nil, err
	}
	if len(podList) == 0 {
		return nil, errors.New("Pod Not Fount")
	}

	var lastCreationTime unversioned.Time
	var lastPod *kapi.Pod

	for _, pod := range podList {
		if lastCreationTime.Before(pod.CreationTimestamp) {
			lastCreationTime = pod.CreationTimestamp
			lastPod = pod
		}

	}
	return lastPod, nil
}

func DeleteReplicaSet(w *watcher.Watcher, replicaSet *extensions.ReplicaSet) error {
	// Update ReplicaSet
	replicaSet, err := w.KubeClient.Extensions().ReplicaSets(replicaSet.Namespace).Get(replicaSet.Name)
	if err != nil {
		return err
	}

	replicaSet.Spec.Replicas = 0
	if _, err := w.KubeClient.Extensions().ReplicaSets(replicaSet.Namespace).Update(replicaSet); err != nil {
		return err
	}

	labelSelector, err := util.GetLabels(w.KubeClient, replicaSet.Namespace, host.TypeReplicasets, replicaSet.Name)
	if err != nil {
		return err
	}

	check := 0
	for {
		time.Sleep(time.Second * 30)
		podList, err := w.Storage.PodStore.List(labelSelector)
		if err != nil {
			return err
		}
		if len(podList) == 0 {
			break
		}

		if check > 6 {
			return errors.New("Fail to delete ReplicaSet Pods")
		}
		check++
	}

	// Delete ReplicaSet
	if err := w.KubeClient.Extensions().ReplicaSets(replicaSet.Namespace).Delete(replicaSet.Name, nil); err != nil {
		return err
	}
	return nil
}

func UpdateReplicaSet(w *watcher.Watcher, replicaSet *extensions.ReplicaSet) (*extensions.ReplicaSet, error) {
	if _, err := w.KubeClient.Extensions().ReplicaSets(replicaSet.Namespace).Update(replicaSet); err != nil {
		return nil, err
	}

	check := 0
	for {
		time.Sleep(time.Second * 30)
		nReplicaset, err := w.Storage.ReplicaSetStore.ReplicaSets(replicaSet.Namespace).Get(replicaSet.Name)
		if err != nil {
			return nil, err
		}
		if nReplicaset.Status.ReadyReplicas == nReplicaset.Status.Replicas {
			return nReplicaset, nil
		}

		if check > 6 {
			return nil, errors.New("Fail to create ReplicaSet")
		}
		check++
	}
}
