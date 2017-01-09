package general

import (
	"errors"
	"strings"

	aci "github.com/appscode/k8s-addons/api"
	"github.com/appscode/searchlight/pkg/client"
	"github.com/appscode/searchlight/pkg/client/k8s"
	"github.com/appscode/searchlight/pkg/controller/host"
	"github.com/appscode/searchlight/test/util"
)

func UpdateAlert(client *k8s.KubeClient, alert *aci.Alert) error {

	namespace := alert.Namespace
	alert, err := client.AppscodeExtensionClient.Alert(namespace).Get(alert.Name)
	if err != nil {
		return err
	}

	alert.Spec.IcingaParam.CheckIntervalSec = 600

	alert, err = client.AppscodeExtensionClient.Alert(namespace).Update(alert)
	if err != nil {
		return err
	}
	return nil
}

func CheckAlertServiceData(context *client.Context, alert *aci.Alert) error {
	checkCommand := alert.Spec.CheckCommand
	objectType := alert.Labels["alert.appscode.com/objectType"]
	objectName := alert.Labels["alert.appscode.com/objectName"]
	namespace := alert.Namespace

	hostType, err := util.GetIcingaHostType(checkCommand, objectType)
	if err != nil {
		return err
	}
	objectList, err := host.GetObjectList(context.KubeClient.Client, checkCommand, hostType, namespace, objectType, objectName, "")
	if err != nil {
		return err
	}

	serviceName := strings.Replace(alert.Name, "_", "-", -1)
	serviceName = strings.Replace(serviceName, ".", "-", -1)

	in := util.IcingaServiceSearchQuery(serviceName, objectList)
	var respService *host.ResponseObject

	if _, err := context.IcingaClient.Objects().Service("").Get([]string{}, in).Do().Into(&respService); err != nil {
		return errors.New("can't check icinga service")
	}

	if len(respService.Results) == 0 {
		return errors.New("No icinga service found")
	}

	for _, service := range respService.Results {
		if service.Attrs.CheckInterval != 600 {
			return errors.New("Service is not updated")
		}
	}
	return nil
}
