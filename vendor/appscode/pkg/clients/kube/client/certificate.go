package client

import (
	"appscode/pkg/clients/kube"

	"k8s.io/kubernetes/pkg/api"
	rest "k8s.io/kubernetes/pkg/client/restclient"
	"k8s.io/kubernetes/pkg/watch"
)

type CertificateNamespacer interface {
	Certificate(namespace string) CertificateInterface
}

type CertificateInterface interface {
	List(opts api.ListOptions) (*kube.CertificateList, error)
	Get(name string) (*kube.Certificate, error)
	Create(certificate *kube.Certificate) (*kube.Certificate, error)
	Update(certificate *kube.Certificate) (*kube.Certificate, error)
	Delete(name string, options *api.DeleteOptions) error
	Watch(opts api.ListOptions) (watch.Interface, error)
	UpdateStatus(certificate *kube.Certificate) (*kube.Certificate, error)
}

type CertificateImpl struct {
	r  rest.Interface
	ns string
}

func newCertificate(c *AppsCodeExtensionsClient, namespace string) *CertificateImpl {
	return &CertificateImpl{c.restClient, namespace}
}

func (c *CertificateImpl) List(opts api.ListOptions) (result *kube.CertificateList, err error) {
	result = &kube.CertificateList{}
	err = c.r.Get().
		Namespace(c.ns).
		Resource("certificates").
		VersionedParams(&opts, ExtendedCodec).
		Do().
		Into(result)
	return
}

func (c *CertificateImpl) Get(name string) (result *kube.Certificate, err error) {
	result = &kube.Certificate{}
	err = c.r.Get().
		Namespace(c.ns).
		Resource("certificates").
		Name(name).
		Do().
		Into(result)
	return
}

func (c *CertificateImpl) Create(certificate *kube.Certificate) (result *kube.Certificate, err error) {
	result = &kube.Certificate{}
	err = c.r.Post().
		Namespace(c.ns).
		Resource("certificates").
		Body(certificate).
		Do().
		Into(result)
	return
}

func (c *CertificateImpl) Update(certificate *kube.Certificate) (result *kube.Certificate, err error) {
	result = &kube.Certificate{}
	err = c.r.Put().
		Namespace(c.ns).
		Resource("certificates").
		Name(certificate.Name).
		Body(certificate).
		Do().
		Into(result)
	return
}

func (c *CertificateImpl) Delete(name string, options *api.DeleteOptions) (err error) {
	return c.r.Delete().
		Namespace(c.ns).
		Resource("certificates").
		Name(name).
		Body(options).
		Do().
		Error()
}

func (c *CertificateImpl) Watch(opts api.ListOptions) (watch.Interface, error) {
	return c.r.Get().
		Prefix("watch").
		Namespace(c.ns).
		Resource("certificates").
		VersionedParams(&opts, api.ParameterCodec).
		Watch()
}

func (c *CertificateImpl) UpdateStatus(certificate *kube.Certificate) (result *kube.Certificate, err error) {
	result = &kube.Certificate{}
	err = c.r.Put().
		Namespace(c.ns).
		Resource("certificates").
		Name(certificate.Name).
		SubResource("status").
		Body(certificate).
		Do().
		Into(result)
	return
}
