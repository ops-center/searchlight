package client

import (
	"k8s.io/kubernetes/pkg/apimachinery/registered"
	"k8s.io/kubernetes/pkg/client/restclient"
)

const (
	defaultAPIPath = "/apis"
)

type AppsCodeExtensionInterface interface {
	IngressNamespacer
	AlertNamespacer
	CertificateNamespacer
}

// AppsCodeExtensionsClient is used to interact with experimental Kubernetes features.
// Features of Extensions group are not supported and may be changed or removed in
// incompatible ways at any time.
type AppsCodeExtensionsClient struct {
	*restclient.RESTClient
}

func (a *AppsCodeExtensionsClient) Ingress(namespace string) IngressInterface {
	return newExtendedIngress(a, namespace)
}

func (a *AppsCodeExtensionsClient) Alert(namespace string) AlertInterface {
	return newAlert(a, namespace)
}

func (a *AppsCodeExtensionsClient) Certificate(namespace string) CertificateInterface {
	return newCertificate(a, namespace)
}

// NewAppsCodeExtensions creates a new AppsCodeExtensionsClient for the given config. This client
// provides access to experimental Kubernetes features.
// Features of Extensions group are not supported and may be changed or removed in
// incompatible ways at any time.
func NewAppsCodeExtensions(c *restclient.Config) (*AppsCodeExtensionsClient, error) {
	config := *c
	if err := setExtensionsDefaults(&config); err != nil {
		return nil, err
	}
	client, err := restclient.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return &AppsCodeExtensionsClient{client}, nil
}

// NewAppsCodeExtensionsOrDie creates a new AppsCodeExtensionsClient for the given config and
// panics if there is an error in the config.
// Features of Extensions group are not supported and may be changed or removed in
// incompatible ways at any time.
func NewAppsCodeExtensionsOrDie(c *restclient.Config) *AppsCodeExtensionsClient {
	client, err := NewAppsCodeExtensions(c)
	if err != nil {
		panic(err)
	}
	return client
}

func setExtensionsDefaults(config *restclient.Config) error {
	config.APIPath = defaultAPIPath
	if config.UserAgent == "" {
		config.UserAgent = restclient.DefaultKubernetesUserAgent()
	}

	contentConfig := ContentConfig()
	if config.NegotiatedSerializer == nil {
		config.NegotiatedSerializer = contentConfig.NegotiatedSerializer
	}
	config.ContentConfig = contentConfig

	if config.GroupVersion == nil || config.GroupVersion.Group != "appscode.com" {
		g, err := registered.Group("appscode.com")
		if err != nil {
			return err
		}
		copyGroupVersion := g.GroupVersion
		config.GroupVersion = &copyGroupVersion
	}

	if config.QPS == 0 {
		config.QPS = 5
	}
	if config.Burst == 0 {
		config.Burst = 10
	}
	return nil
}
