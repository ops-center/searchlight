package e2e

import (
	"fmt"
	"os/exec"
	"syscall"
	"testing"
	"time"

	"github.com/appscode/go/crypto/rand"
	config "github.com/appscode/searchlight/pkg/client/k8s"
	"github.com/appscode/searchlight/pkg/controller/host"
	"github.com/appscode/searchlight/test/plugin"
	"github.com/appscode/searchlight/test/plugin/component_status"
	"github.com/appscode/searchlight/test/plugin/node_count"
	"github.com/appscode/searchlight/test/plugin/node_status"
	"github.com/appscode/searchlight/test/plugin/pod_status"
	"github.com/stretchr/testify/assert"
)

type testData struct {
	data         map[string]interface{}
	expectedCode int
	deleteObject bool
}

func getKubernetesClient() *config.KubeClient {
	kubeClient, err := config.NewClient()
	plugin.Fatalln(err)
	return kubeClient
}

func execCheckCommand(cmdName string, cmdArgs ...string) (code int) {
	cmd := fmt.Sprintf("../dist/%s/%s-%s-%s", cmdName, cmdName, GOHOSTOS, GOHOSTARCH)
	command := exec.Command(cmd, cmdArgs...)
	fmt.Println(fmt.Sprintf("Running: %v", command.Args))
	cmdOut, err := command.Output()
	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				code = status.ExitStatus()
			}
		}
	}
	fmt.Println(string(cmdOut))
	return
}

func TestComponentStatus(t *testing.T) {
	fmt.Println("== Testing >", host.CheckComponentStatus)

	statusCode := execCheckCommand("hyperalert", "check_component_status")
	kubeClient := getKubernetesClient()
	expectedStatusCode := component_status.GetStatusCodeForComponentStatus(kubeClient)
	assert.EqualValues(t, expectedStatusCode, statusCode)
}

func TestJsonPath(t *testing.T) {
	fmt.Println("== Testing >", host.CheckJsonPath)

	url := "https://api.appscode.info"
	uri := "/health"
	testDataList := []testData{
		testData{
			data: map[string]interface{}{
				"url":     url + uri,
				"query":   ".status",
				"warning": `.status!="OK"`,
			},
			expectedCode: 0,
		},
		testData{
			data: map[string]interface{}{
				"url":     url + uri,
				"query":   ".version",
				"warning": `.version=="0.3"`,
			},
			expectedCode: 1,
		},
		testData{
			data: map[string]interface{}{
				"url":     url + uri,
				"query":   ".version",
				"warning": `.version=="0.2.6"`, "critical": `.version=="0.3"`,
			},
			expectedCode: 2,
		},
		testData{
			data: map[string]interface{}{
				"url":     url + "/healthz",
				"query":   ".version",
				"warning": `.version=="0.2.6"`,
			},
			expectedCode: 3,
		},
	}

	for _, testData := range testDataList {
		argList := []string{
			"check_json_path",
		}
		for key, val := range testData.data {
			argList = append(argList, fmt.Sprintf(`--%s=%s`, key, val))
		}
		statusCode := execCheckCommand("hyperalert", argList...)
		assert.EqualValues(t, testData.expectedCode, statusCode)
	}
}

func TestNodeCount(t *testing.T) {
	fmt.Println("== Testing >", host.CheckNodeCount)

	kubeClient := getKubernetesClient()
	actualNodeCount := node_count.GetKubernetesNodeCount(kubeClient)

	testDataList := []testData{
		testData{
			data: map[string]interface{}{
				"count": actualNodeCount,
			},
			expectedCode: 0,
		},
		testData{
			data: map[string]interface{}{
				"count": actualNodeCount + 1,
			},
			expectedCode: 2,
		},
		testData{
			data:         map[string]interface{}{},
			expectedCode: 3,
		},
	}

	for _, testData := range testDataList {
		argList := []string{
			"check_node_count",
		}
		for key, val := range testData.data {
			argList = append(argList, fmt.Sprintf("--%s=%v", key, val))
		}
		statusCode := execCheckCommand("hyperalert", argList...)
		assert.EqualValues(t, testData.expectedCode, statusCode)
	}
}

func TestNodeStatus(t *testing.T) {
	fmt.Println("== Testing >", host.CheckNodeStatus)

	kubeClient := getKubernetesClient()
	actualNodeName := node_status.GetKubernetesNodeName(kubeClient)
	hostname := actualNodeName + "@default"

	testDataList := []testData{
		testData{
			data: map[string]interface{}{
				"host": hostname,
			},
			expectedCode: 0,
		},
		testData{
			data: map[string]interface{}{
				// make node name invalid using random 2 character.
				// Added as prefix because 1st part of hostname is nodename. (<node-name>@<alert-namespace>)
				"host": rand.Characters(2) + hostname,
			},
			expectedCode: 3,
		},
	}

	for _, testData := range testDataList {
		argList := []string{
			"check_node_status",
		}
		for key, val := range testData.data {
			argList = append(argList, fmt.Sprintf("--%s=%v", key, val))
		}
		statusCode := execCheckCommand("hyperalert", argList...)
		assert.EqualValues(t, testData.expectedCode, statusCode)
	}
}

func TestPodExists(t *testing.T) {
	fmt.Println("== Testing >", host.CheckCommandPodExists)

	kubectlClient := getKubernetesClient()
	testPodExists := func(dataConfig *DataConfig) {
		// This will create object & return icinga_host name
		// and number of pods under it
		name, count := GetTestData(kubectlClient.Client, dataConfig)
		time.Sleep(time.Second * 30)

		testDataList := []testData{
			testData{
				// To check for any pods
				data: map[string]interface{}{
					"host": name,
				},
				expectedCode: 0,
			},
			testData{
				// To check for specific number of pods
				data: map[string]interface{}{
					"host":  name,
					"count": count,
				},
				expectedCode: 0,
			},
			testData{
				// To check for critical when pod number mismatch
				data: map[string]interface{}{
					"host":  name,
					"count": count + 1,
				},
				expectedCode: 2,
				deleteObject: true,
			},
		}

		for _, testData := range testDataList {
			argList := []string{
				"check_pod_exists",
			}
			for key, val := range testData.data {
				argList = append(argList, fmt.Sprintf("--%s=%v", key, val))
			}
			statusCode := execCheckCommand("hyperalert", argList...)
			assert.EqualValues(t, testData.expectedCode, statusCode)
		}
	}

	ns := rand.WithUniqSuffix("e2e")
	dataConfig := &DataConfig{
		Namespace: ns,
	}

	fmt.Println(">> Creating namespace", ns)
	CreateNewNamespace(kubectlClient.Client, ns)
	fmt.Println()

	fmt.Println(">> Testing plugings for", host.TypeReplicationcontrollers)
	dataConfig.ObjectType = host.TypeReplicationcontrollers
	testPodExists(dataConfig)

	fmt.Println(">> Testing plugings for", host.TypeReplicasets)
	dataConfig.ObjectType = host.TypeReplicasets
	testPodExists(dataConfig)

	fmt.Println(">> Testing plugings for", host.TypeDaemonsets)
	dataConfig.ObjectType = host.TypeDaemonsets
	testPodExists(dataConfig)

	fmt.Println(">> Testing plugings for", host.TypeDeployments)
	dataConfig.ObjectType = host.TypeDeployments
	testPodExists(dataConfig)

	fmt.Println(">> Testing plugings for", host.TypeServices)
	dataConfig.ObjectType = host.TypeServices
	testPodExists(dataConfig)

	fmt.Println(">> Testing plugings for", host.TypeCluster)
	dataConfig.ObjectType = host.TypeCluster
	dataConfig.CheckCommand = host.CheckCommandPodExists
	testPodExists(dataConfig)

	fmt.Println(">> Deleting namespace", ns)
	deleteNewNamespace(kubectlClient.Client, ns)

	fmt.Println()
}

func TestPodStatus(t *testing.T) {
	fmt.Println("== Testing >", host.CheckCommandPodStatus)

	kubectlClient := getKubernetesClient()

	testPodStatus := func(dataConfig *DataConfig) {
		// This will create object & return icinga_host name
		name, _ := GetTestData(kubectlClient.Client, dataConfig)
		time.Sleep(time.Second * 30)

		// This will check pod status under specific object
		// and will return 2 (critical) if any pod is not running
		expectedCode := pod_status.GetStatusCodeForPodStatus(kubectlClient, name)

		testDataList := []testData{
			testData{
				data: map[string]interface{}{
					"host": name,
				},
				expectedCode: expectedCode,
			},
		}

		for _, testData := range testDataList {
			argList := []string{
				"check_pod_status",
			}
			for key, val := range testData.data {
				argList = append(argList, fmt.Sprintf("--%s=%v", key, val))
			}
			statusCode := execCheckCommand("hyperalert", argList...)
			assert.EqualValues(t, testData.expectedCode, statusCode)
		}
	}

	ns := rand.WithUniqSuffix("e2e")
	dataConfig := &DataConfig{
		Namespace: ns,
	}

	fmt.Println(">> Creating namespace", ns)
	CreateNewNamespace(kubectlClient.Client, ns)
	fmt.Println()

	fmt.Println(">> Testing plugings for", host.TypeReplicationcontrollers)
	dataConfig.ObjectType = host.TypeReplicationcontrollers
	testPodStatus(dataConfig)

	fmt.Println(">> Testing plugings for", host.TypeReplicasets)
	dataConfig.ObjectType = host.TypeReplicasets
	testPodStatus(dataConfig)

	fmt.Println(">> Testing plugings for", host.TypeDaemonsets)
	dataConfig.ObjectType = host.TypeDaemonsets
	testPodStatus(dataConfig)

	fmt.Println(">> Testing plugings for", host.TypeDeployments)
	dataConfig.ObjectType = host.TypeDeployments
	testPodStatus(dataConfig)

	fmt.Println(">> Testing plugings for", host.TypeServices)
	dataConfig.ObjectType = host.TypeServices
	testPodStatus(dataConfig)

	fmt.Println(">> Testing plugings for", host.TypePods)
	dataConfig.ObjectType = host.TypePods
	testPodStatus(dataConfig)

	fmt.Println(">> Testing plugings for", host.TypeCluster)
	dataConfig.ObjectType = host.TypeCluster
	dataConfig.CheckCommand = host.CheckCommandPodStatus
	testPodStatus(dataConfig)

	fmt.Println(">> Deleting namespace", ns)
	deleteNewNamespace(kubectlClient.Client, ns)

	fmt.Println()
}
