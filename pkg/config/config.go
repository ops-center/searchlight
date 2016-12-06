package config

import (
	ext_client "appscode/pkg/clients/kube/client"

	kube_client "k8s.io/kubernetes/pkg/client/unversioned"
)

const (
	TypeServices               = "services"
	TypeReplicationcontrollers = "replicationcontrollers"
	TypeDaemonsets             = "daemonsets"
	TypePetsets                = "petsets"
	TypeReplicasets            = "replicasets"
	TypeDeployments            = "deployments"
	TypePods                   = "pods"
)

type KubeClient struct {
	*kube_client.Client
	*ext_client.AppsCodeExtensionsClient
}
