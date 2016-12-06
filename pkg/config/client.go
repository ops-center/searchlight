package config

import (
	ext_client "appscode/pkg/clients/kube/client"

	"k8s.io/kubernetes/pkg/client/restclient"
	kClient "k8s.io/kubernetes/pkg/client/unversioned"
)

func GetKubeClient() (*KubeClient, error) {
	var clientConfig *restclient.Config
	var err error
	// Set debugMode to "true" for testing from local
	debugMode := false
	if !debugMode {
		clientConfig, err = restclient.InClusterConfig()
		if err != nil {
			return nil, err
		}
	} else {
		host := ""
		username := ""
		password := ""
		clientConfig = &restclient.Config{
			Host:     host,
			Insecure: true,
			Username: username,
			Password: password,
		}
	}

	client, err := kClient.New(clientConfig)
	if err != nil {
		return nil, err
	}
	acExtClient, err := ext_client.NewAppsCodeExtensions(clientConfig)
	if err != nil {
		return nil, err
	}

	return &KubeClient{client, acExtClient}, nil
}
