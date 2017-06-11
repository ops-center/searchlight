package mini

import (
	"errors"
	"time"

	"github.com/appscode/searchlight/pkg/controller/host"
	"github.com/appscode/searchlight/pkg/testing"
	"github.com/appscode/searchlight/pkg/watcher"
	"github.com/appscode/searchlight/util"
	kapi "k8s.io/kubernetes/pkg/api"
)

func CreateReplicationController(w *watcher.Watcher, namespace string) (*kapi.ReplicationController, error) {
	replicationController := &kapi.ReplicationController{}
	replicationController.Namespace = namespace
	if err := testing.CreateKubernetesObject(w.KubeClient, replicationController); err != nil {
		return nil, err
	}

	check := 0
	for {
		time.Sleep(time.Second * 30)
		nReplicationController, err := w.Storage.RcStore.ReplicationControllers(replicationController.Namespace).Get(replicationController.Name)
		if err != nil {
			return nil, err
		}
		if nReplicationController.Status.ReadyReplicas == nReplicationController.Status.Replicas {
			break
		}

		if check > 6 {
			return nil, errors.New("Fail to create ReplicationController")
		}
		check++
	}

	return replicationController, nil
}

func DeleteReplicationController(w *watcher.Watcher, replicationController *kapi.ReplicationController) error {
	replicationController, err := w.KubeClient.Core().ReplicationControllers(replicationController.Namespace).Get(replicationController.Name)
	if err != nil {
		return err
	}
	// Update ReplicationController
	replicationController.Spec.Replicas = 0
	if _, err := w.KubeClient.Core().ReplicationControllers(replicationController.Namespace).Update(replicationController); err != nil {
		return err
	}

	labelSelector, err := util.GetLabels(w.KubeClient, replicationController.Namespace, host.TypeReplicationcontrollers, replicationController.Name)
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
			return errors.New("Fail to delete ReplicationController Pods")
		}
		check++
	}

	// Delete ReplicationController
	if err := w.KubeClient.Core().ReplicationControllers(replicationController.Namespace).Delete(replicationController.Name, nil); err != nil {
		return err
	}
	return nil
}
