package config

import (
	extclient "appscode/pkg/clients/kube/client"
	_ "appscode/pkg/clients/kube/install"

	clientset "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
	rest "k8s.io/kubernetes/pkg/client/restclient"
)

const (
	TypeServices               = "services"
	TypeReplicationcontrollers = "replicationcontrollers"
	TypeDaemonsets             = "daemonsets"
	TypeStatefulSet            = "statefulsets"
	TypeReplicasets            = "replicasets"
	TypeDeployments            = "deployments"
	TypePods                   = "pods"
)

type KubeClient struct {
	Client                  clientset.Interface
	AppscodeExtensionClient extclient.AppsCodeExtensionInterface

	config *rest.Config
}
