package driver

import (
	"appscode/pkg/clients/kube"
	"os"

	"github.com/appscode/log"
	"github.com/appscode/searchlight/pkg/config"
)

func GetAlertInfo(namespace, alertName string) (*kube.Alert, error) {
	kubeClient, err := config.NewKubeClient()
	if err != nil {
		log.Errorln(err)
		os.Exit(1)
	}
	return kubeClient.AppscodeExtensionClient.Alert(namespace).Get(alertName)
}
