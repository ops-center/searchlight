package k8s

import (
	_ "github.com/appscode/k8s-addons/api/install"
	acs "github.com/appscode/k8s-addons/client/clientset"
	clientset "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
)

type KubeClient struct {
	Client                  *clientset.Clientset
	AppscodeExtensionClient *acs.AppsCodeExtensionsClient
}
