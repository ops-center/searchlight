package k8s

import (
	"github.com/appscode/errors"
	_env "github.com/appscode/go/env"
	_ "github.com/appscode/k8s-addons/api/install"
	acs "github.com/appscode/k8s-addons/client/clientset"
	clientset "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
	rest "k8s.io/kubernetes/pkg/client/restclient"
)

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
		host := ""
		username := ""
		password := ""
		config = &rest.Config{
			Host:     host,
			Insecure: true,
			Username: username,
			Password: password,
		}
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
		config:                  config,
		Client:                  client,
		AppscodeExtensionClient: appscodeClient,
	}, nil
}
