package config

import (
	"k8s.io/kubernetes/pkg/client/restclient"
	kClient "k8s.io/kubernetes/pkg/client/unversioned"
)

func GetKubeClient() (*KubeClient, error) {
	// Set debugMode to "true" for testing from local
	debugMode := false
	if !debugMode {
		client, err := kClient.NewInCluster()
		if err != nil {
			return nil, err
		}
		return &KubeClient{client}, nil
	} else {
		host := ""
		username := ""
		password := ""
		config := &restclient.Config{
			Host:     host,
			Insecure: true,
			Username: username,
			Password: password,
		}

		client, err := kClient.New(config)
		if err != nil {
			return nil, err
		}
		return &KubeClient{client}, nil
	}
}
