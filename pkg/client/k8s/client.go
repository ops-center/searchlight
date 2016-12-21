package k8s

import (
	extclient "appscode/pkg/clients/kube/client"
	_ "appscode/pkg/clients/kube/install"

	"github.com/appscode/errors"
	_env "github.com/appscode/go-env"
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

	appscodeClient, err := extclient.NewACExtensionsForConfig(config)
	if err != nil {
		return nil, errors.New().WithCause(err).Internal()
	}

	return &KubeClient{
		config:                  config,
		Client:                  client,
		AppscodeExtensionClient: appscodeClient,
	}, nil
}
