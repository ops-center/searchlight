package kube_event

import (
	"time"

	"github.com/appscode/go/crypto/rand"
	"github.com/appscode/searchlight/pkg/client/k8s"
	"github.com/appscode/searchlight/test/plugin"
	"github.com/appscode/searchlight/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/pkg/api"
	apiv1 "k8s.io/client-go/pkg/api/v1"
)

func getStatusCodeForEventCount(kubeClient *k8s.KubeClient, checkInterval, clockSkew time.Duration) (util.IcingaState, error) {

	now := time.Now()
	// Create some fake event
	for i := 0; i < 5; i++ {
		_, err := kubeClient.Client.CoreV1().Events(apiv1.NamespaceDefault).Create(&apiv1.Event{
			ObjectMeta: metav1.ObjectMeta{
				Name: rand.WithUniqSuffix("event"),
			},
			Type:           apiv1.EventTypeWarning,
			FirstTimestamp: metav1.NewTime(now),
			LastTimestamp:  metav1.NewTime(now),
		})
		if err != nil {
			return util.Unknown, err
		}
	}

	count := 0
	field := fields.OneTermEqualSelector(api.EventTypeField, apiv1.EventTypeWarning)
	checkTime := time.Now().Add(-(checkInterval + clockSkew))
	eventList, err := kubeClient.Client.CoreV1().Events(apiv1.NamespaceAll).List(metav1.ListOptions{
		FieldSelector: field.String(),
	})
	if err != nil {
		return util.Unknown, err
	}

	for _, event := range eventList.Items {
		if checkTime.Before(event.LastTimestamp.Time) {
			count = count + 1
		}
	}

	if count > 0 {
		return util.Warning, nil
	}
	return util.Ok, nil
}

func GetTestData(kubeClient *k8s.KubeClient, checkInterval, clockSkew time.Duration) ([]plugin.TestData, error) {
	expectedIcingaState, err := getStatusCodeForEventCount(kubeClient, checkInterval, clockSkew)
	if err != nil {
		return nil, err
	}
	testDataList := []plugin.TestData{
		{
			Data: map[string]interface{}{
				"CheckInterval": checkInterval,
				"ClockSkew":     clockSkew,
			},
			ExpectedIcingaState: expectedIcingaState,
		},
	}

	return testDataList, nil
}
