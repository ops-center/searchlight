package config

import kube_client "k8s.io/kubernetes/pkg/client/unversioned"

const (
	TypeServices               = "services"
	TypeReplicationcontrollers = "replicationcontrollers"
	TypeDaemonsets             = "daemonsets"
	TypePetsets                = "petsets"
	TypeReplicasets            = "replicasets"
	TypeDeployments            = "deployments"
	TypePods                   = "pods"
	TypeNodes                  = "nodes"
	TypeCluster                = "cluster"
)

type KubeClient struct {
	*kube_client.Client
}
