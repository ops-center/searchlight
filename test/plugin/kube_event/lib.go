package kube_event

import (
	"time"

	"fmt"
	"github.com/appscode/searchlight/pkg/client/k8s"
	"github.com/appscode/searchlight/test/plugin"
	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/fields"
	"k8s.io/kubernetes/pkg/labels"
	"os"
)

func GetStatusCodeForEventCount(kubeClient *k8s.KubeClient, checkInterval, clockSkew time.Duration) int {
	var namespaceList *kapi.NamespaceList
	namespaceList, err := kubeClient.Client.Core().Namespaces().List(
		kapi.ListOptions{
			LabelSelector: labels.Everything(),
		},
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	count := 0
	field := fields.OneTermEqualSelector(kapi.EventTypeField, kapi.EventTypeWarning)
	checkTime := time.Now().Add(-(checkInterval + clockSkew))
	for _, ns := range namespaceList.Items {
		eventList, err := kubeClient.Client.Core().Events(ns.Name).List(
			kapi.ListOptions{
				FieldSelector: field,
			},
		)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		for _, event := range eventList.Items {
			if checkTime.Before(event.LastTimestamp.Time) {
				count = count + 1
			}
		}
	}

	if count > 0 {
		return plugin.WARNING
	}

	return plugin.OK
}
