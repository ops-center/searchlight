package mini

import (
	"github.com/appscode/searchlight/pkg/testing"
	"github.com/appscode/searchlight/pkg/watcher"
	kapi "k8s.io/kubernetes/pkg/api"
)

func CreateService(w *watcher.Watcher, namespace string, selector map[string]string) (*kapi.Service, error) {
	service := &kapi.Service{
		ObjectMeta: kapi.ObjectMeta{
			Namespace: namespace,
		},
		Spec: kapi.ServiceSpec{
			Selector: selector,
		},
	}
	if err := testing.CreateKubernetesObject(w.KubeClient, service); err != nil {
		return nil, err
	}
	return service, nil
}

func DeleteService(w *watcher.Watcher, service *kapi.Service) error {
	// Delete Service
	if err := w.KubeClient.Core().Services(service.Namespace).Delete(service.Name, nil); err != nil {
		return err
	}
	return nil
}
