package node_count

import (
	"github.com/appscode/searchlight/pkg/client/k8s"
	"github.com/appscode/searchlight/test/plugin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func getKubernetesNodeCount(kubeClient *k8s.KubeClient) (int, error) {
	nodeList, err := kubeClient.Client.CoreV1().Nodes().List(metav1.ListOptions{
		LabelSelector: labels.Everything().String(),
	})
	if err != nil {
		return 0, err
	}
	return len(nodeList.Items), nil
}

func GetTestData(kubeClient *k8s.KubeClient) ([]plugin.TestData, error) {
	actualNodeCount, err := getKubernetesNodeCount(kubeClient)
	if err != nil {
		return nil, err
	}

	testDataList := []plugin.TestData{
		{
			Data: map[string]interface{}{
				"Count": actualNodeCount,
			},
			ExpectedIcingaState: 0,
		},
		{
			Data: map[string]interface{}{
				"Count": actualNodeCount + 1,
			},
			ExpectedIcingaState: 2,
		},
	}

	return testDataList, nil
}
