package client

import (
	"appscode/pkg/clients/kube"

	"k8s.io/kubernetes/pkg/api"
	rest "k8s.io/kubernetes/pkg/client/restclient"
	"k8s.io/kubernetes/pkg/watch"
)

type AlertNamespacer interface {
	Alert(namespace string) AlertInterface
}

type AlertInterface interface {
	List(opts api.ListOptions) (*kube.AlertList, error)
	Get(name string) (*kube.Alert, error)
	Create(Alert *kube.Alert) (*kube.Alert, error)
	Update(Alert *kube.Alert) (*kube.Alert, error)
	Delete(name string, options *api.DeleteOptions) error
	Watch(opts api.ListOptions) (watch.Interface, error)
	UpdateStatus(Alert *kube.Alert) (*kube.Alert, error)
}

type AlertImpl struct {
	r  rest.Interface
	ns string
}

func newAlert(c *AppsCodeExtensionsClient, namespace string) *AlertImpl {
	return &AlertImpl{c.restClient, namespace}
}

func (c *AlertImpl) List(opts api.ListOptions) (result *kube.AlertList, err error) {
	result = &kube.AlertList{}
	err = c.r.Get().
		Namespace(c.ns).
		Resource("alerts").
		VersionedParams(&opts, ExtendedCodec).
		Do().
		Into(result)
	return
}

func (c *AlertImpl) Get(name string) (result *kube.Alert, err error) {
	result = &kube.Alert{}
	err = c.r.Get().
		Namespace(c.ns).
		Resource("alerts").
		Name(name).
		Do().
		Into(result)
	return
}

func (c *AlertImpl) Create(alert *kube.Alert) (result *kube.Alert, err error) {
	result = &kube.Alert{}
	err = c.r.Post().
		Namespace(c.ns).
		Resource("alerts").
		Body(alert).
		Do().
		Into(result)
	return
}

func (c *AlertImpl) Update(alert *kube.Alert) (result *kube.Alert, err error) {
	result = &kube.Alert{}
	err = c.r.Put().
		Namespace(c.ns).
		Resource("alerts").
		Name(alert.Name).
		Body(alert).
		Do().
		Into(result)
	return
}

func (c *AlertImpl) Delete(name string, options *api.DeleteOptions) (err error) {
	return c.r.Delete().
		Namespace(c.ns).
		Resource("alerts").
		Name(name).
		Body(options).
		Do().
		Error()
}

func (c *AlertImpl) Watch(opts api.ListOptions) (watch.Interface, error) {
	return c.r.Get().
		Prefix("watch").
		Namespace(c.ns).
		Resource("alerts").
		VersionedParams(&opts, api.ParameterCodec).
		Watch()
}

func (c *AlertImpl) UpdateStatus(alert *kube.Alert) (result *kube.Alert, err error) {
	result = &kube.Alert{}
	err = c.r.Put().
		Namespace(c.ns).
		Resource("alerts").
		Name(alert.Name).
		SubResource("status").
		Body(alert).
		Do().
		Into(result)
	return
}
