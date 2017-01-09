package client

import (
	_ "github.com/appscode/k8s-addons/api/install"
	"github.com/appscode/searchlight/pkg/client/icinga"
	"github.com/appscode/searchlight/pkg/client/k8s"
)

type Context struct {
	KubeClient   *k8s.KubeClient
	IcingaClient *icinga.IcingaClient
}
