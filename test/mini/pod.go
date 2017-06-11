package mini

import (
	"errors"
	"time"

	"github.com/appscode/searchlight/pkg/testing"
	"github.com/appscode/searchlight/pkg/watcher"
	kapi "k8s.io/kubernetes/pkg/api"
)

func checkPod(w *watcher.Watcher, pod *kapi.Pod) (*kapi.Pod, error) {
	check := 0
	for {
		time.Sleep(time.Second * 30)
		nPod, err := w.Storage.PodStore.Pods(pod.Namespace).Get(pod.Name)
		if err != nil {
			return nil, err
		}
		if nPod.Status.Phase == kapi.PodRunning {
			return nPod, nil
		}

		if check > 6 {
			return nil, errors.New("Fail to create Pod")
		}
		check++
	}
}

func CreatePod(w *watcher.Watcher, namespace string) (*kapi.Pod, error) {
	pod := &kapi.Pod{}
	pod.Namespace = namespace
	if err := testing.CreateKubernetesObject(w.KubeClient, pod); err != nil {
		return nil, err
	}

	return checkPod(w, pod)
}

func ReCreatePod(w *watcher.Watcher, pod *kapi.Pod) (*kapi.Pod, error) {
	newPod := &kapi.Pod{
		ObjectMeta: kapi.ObjectMeta{
			Name:      pod.Name,
			Namespace: pod.Namespace,
		},
	}
	if err := testing.CreateKubernetesObject(w.KubeClient, newPod); err != nil {
		return nil, err
	}

	return checkPod(w, newPod)
}

func DeletePod(w *watcher.Watcher, pod *kapi.Pod) error {
	// Delete Pod
	if err := w.KubeClient.Core().Pods(pod.Namespace).Delete(pod.Name, nil); err != nil {
		return err
	}
	return nil
}
