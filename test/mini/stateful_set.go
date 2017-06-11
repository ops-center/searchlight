package mini

import (
	"errors"
	"time"

	"github.com/appscode/searchlight/pkg/controller/host"
	"github.com/appscode/searchlight/pkg/testing"
	"github.com/appscode/searchlight/pkg/watcher"
	"github.com/appscode/searchlight/util"
	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/apps"
)

func CreateStatefulSet(w *watcher.Watcher, namespace string) (*apps.StatefulSet, error) {
	// Create Service
	service, err := CreateService(w, namespace, nil)
	if err != nil {
		return nil, err
	}

	statefulSet := &apps.StatefulSet{
		ObjectMeta: kapi.ObjectMeta{
			Namespace: namespace,
		},
		Spec: apps.StatefulSetSpec{
			ServiceName: service.Name,
		},
	}

	if err := testing.CreateKubernetesObject(w.KubeClient, statefulSet); err != nil {
		return nil, err
	}

	check := 0
	for {
		time.Sleep(time.Second * 30)
		nStatefulSet, exists, err := w.Storage.StatefulSetStore.Get(statefulSet)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, errors.New("StatefulSet not found")
		}

		if nStatefulSet.(*apps.StatefulSet).Status.Replicas == statefulSet.Spec.Replicas {
			return nStatefulSet.(*apps.StatefulSet), nil
		}

		if check > 6 {
			return nil, errors.New("Fail to create StatefulSet")
		}
		check++
	}
}

func DeleteStatefulSet(w *watcher.Watcher, statefulSet *apps.StatefulSet) error {
	statefulSet, err := w.KubeClient.Apps().StatefulSets(statefulSet.Namespace).Get(statefulSet.Name)
	if err != nil {
		return err
	}
	// Update StatefulSet
	statefulSet.Spec.Replicas = 0
	if _, err := w.KubeClient.Apps().StatefulSets(statefulSet.Namespace).Update(statefulSet); err != nil {
		return err
	}

	labelSelector, err := util.GetLabels(w.KubeClient, statefulSet.Namespace, host.TypeStatefulSet, statefulSet.Name)
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
			return errors.New("Fail to delete StatefulSet Pods")
		}
		check++
	}

	// Delete StatefulSet
	if err := w.KubeClient.Apps().StatefulSets(statefulSet.Namespace).Delete(statefulSet.Name, nil); err != nil {
		return err
	}

	return DeleteService(w, &kapi.Service{
		ObjectMeta: kapi.ObjectMeta{
			Name:      statefulSet.Spec.ServiceName,
			Namespace: statefulSet.Namespace,
		},
	})
}
