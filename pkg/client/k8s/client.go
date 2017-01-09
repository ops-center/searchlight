package k8s

import (
	"fmt"

	"github.com/appscode/errors"
	_env "github.com/appscode/go/env"
	_ "github.com/appscode/k8s-addons/api/install"
	acs "github.com/appscode/k8s-addons/client/clientset"
	clientset "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
	rest "k8s.io/kubernetes/pkg/client/restclient"
	"k8s.io/kubernetes/pkg/client/unversioned/clientcmd"
)

// NewClient() should only be used to create kube client for plugins.
func NewClient() (*KubeClient, error) {
	var config *rest.Config
	var err error

	debugEnabled := _env.FromHost().DebugEnabled()
	if !debugEnabled {
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	} else {
		rules := clientcmd.NewDefaultClientConfigLoadingRules()
		rules.DefaultClientConfig = &clientcmd.DefaultClientConfig
		overrides := &clientcmd.ConfigOverrides{ClusterDefaults: clientcmd.ClusterDefaults}
		config, err = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, overrides).ClientConfig()
		if err != nil {
			return nil, fmt.Errorf("Could not get kubernetes config: %s", err)
		}
		fmt.Println("Using cluster:", config.Host)
	}

	client, err := clientset.NewForConfig(config)
	if err != nil {
		return nil, errors.New().WithCause(err).Internal()
	}

	appscodeClient, err := acs.NewACExtensionsForConfig(config)
	if err != nil {
		return nil, errors.New().WithCause(err).Internal()
	}

	return &KubeClient{
		Client:                  client,
		AppscodeExtensionClient: appscodeClient,
	}, nil
}
