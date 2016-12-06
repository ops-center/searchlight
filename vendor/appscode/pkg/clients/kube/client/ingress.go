package client

import (
	"appscode/pkg/clients/kube"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/watch"
)

// ExtendedIngressNamespacer has methods to work with ExtendedIngress resources in a namespace
type IngressNamespacer interface {
	Ingress(namespace string) IngressInterface
}

// ExtendedIngressInterface exposes methods to work on ExtendedIngress resources.
type IngressInterface interface {
	List(opts api.ListOptions) (*kube.IngressList, error)
	Get(name string) (*kube.Ingress, error)
	Create(ExtendedIngress *kube.Ingress) (*kube.Ingress, error)
	Update(ExtendedIngress *kube.Ingress) (*kube.Ingress, error)
	Delete(name string, options *api.DeleteOptions) error
	Watch(opts api.ListOptions) (watch.Interface, error)
	UpdateStatus(ExtendedIngress *kube.Ingress) (*kube.Ingress, error)
}

// ExtendedIngress implements ExtendedIngressNamespacer interface
type IngressImpl struct {
	r  *AppsCodeExtensionsClient
	ns string
}

// newExtendedIngress returns a ExtendedIngress
func newExtendedIngress(c *AppsCodeExtensionsClient, namespace string) *IngressImpl {
	return &IngressImpl{c, namespace}
}

// List returns a list of ExtendedIngress that match the label and field selectors.
func (c *IngressImpl) List(opts api.ListOptions) (result *kube.IngressList, err error) {
	result = &kube.IngressList{}
	err = c.r.Get().
		Namespace(c.ns).
		Resource("ingresses").
		VersionedParams(&opts, ExtendedCodec).
		Do().
		Into(result)
	return
}

// Get returns information about a particular ExtendedIngress.
func (c *IngressImpl) Get(name string) (result *kube.Ingress, err error) {
	result = &kube.Ingress{}
	err = c.r.Get().
		Namespace(c.ns).
		Resource("ingresses").
		Name(name).
		Do().
		Into(result)
	return
}

// Create creates a new ExtendedIngress.
func (c *IngressImpl) Create(extendedIngress *kube.Ingress) (result *kube.Ingress, err error) {
	result = &kube.Ingress{}
	err = c.r.Post().
		Namespace(c.ns).
		Resource("ingresses").
		Body(extendedIngress).
		Do().
		Into(result)
	return
}

// Update updates an existing ExtendedIngress.
func (c *IngressImpl) Update(extendedIngress *kube.Ingress) (result *kube.Ingress, err error) {
	result = &kube.Ingress{}
	err = c.r.Put().
		Namespace(c.ns).
		Resource("ingresses").
		Name(extendedIngress.Name).
		Body(extendedIngress).
		Do().
		Into(result)
	return
}

// Delete deletes a ExtendedIngress, returns error if one occurs.
func (c *IngressImpl) Delete(name string, options *api.DeleteOptions) (err error) {
	return c.r.Delete().
		Namespace(c.ns).
		Resource("ingresses").
		Name(name).
		Body(options).
		Do().
		Error()
}

// Watch returns a watch.Interface that watches the requested ExtendedIngress.
func (c *IngressImpl) Watch(opts api.ListOptions) (watch.Interface, error) {
	return c.r.Get().
		Prefix("watch").
		Namespace(c.ns).
		Resource("ingresses").
		VersionedParams(&opts, ExtendedCodec).
		Watch()
}

// UpdateStatus takes the name of the ExtendedIngress and the new status.  Returns the server's representation of the ExtendedIngress, and an error, if it occurs.
func (c *IngressImpl) UpdateStatus(extendedIngress *kube.Ingress) (result *kube.Ingress, err error) {
	result = &kube.Ingress{}
	err = c.r.Put().
		Namespace(c.ns).
		Resource("ingresses").
		Name(extendedIngress.Name).
		SubResource("status").
		Body(extendedIngress).
		Do().
		Into(result)
	return
}
