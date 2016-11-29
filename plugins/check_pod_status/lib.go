package check_pod_status

import (
	"fmt"
	"os"

	"encoding/json"

	"github.com/appscode/searchlight/pkg/config"
	"github.com/appscode/searchlight/pkg/util"
	"github.com/spf13/cobra"
	kApi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/labels"
)

type request struct {
	namespace  string
	objectType string
	objectName string
}

type objectInfo struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Status    string `json:"status,omitempty"`
}

type serviceOutput struct {
	Objects []*objectInfo `json:"objects,omitempty"`
	Message string        `json:"message,omitempty"`
}

func checkPodStatus(cmd *cobra.Command, req *request) {
	kubeClient, err := config.GetKubeClient()
	if err != nil {
		fmt.Fprintln(os.Stdout, util.State[3], err)
		os.Exit(3)
	}

	objectInfoList := make([]*objectInfo, 0)
	if req.objectType == config.TypePods {
		pod, err := kubeClient.Pods(req.namespace).Get(req.objectName)
		if err != nil {
			fmt.Fprintln(os.Stdout, util.State[3], err)
			os.Exit(3)
		}

		if !(pod.Status.Phase == kApi.PodSucceeded || pod.Status.Phase == kApi.PodRunning) {
			objectInfoList = append(objectInfoList, &objectInfo{Name: pod.Name, Status: string(pod.Status.Phase), Namespace: pod.Namespace})
		}
	} else {
		var labelSelector labels.Selector
		if req.objectType == "" {
			labelSelector = labels.Everything()
		} else {
			if labelSelector, err = kubeClient.GetLabels(req.namespace, req.objectType, req.objectName); err != nil {
				fmt.Fprintln(os.Stdout, util.State[3], err)
				os.Exit(3)
			}
		}

		podList, err := kubeClient.Pods(req.namespace).List(kApi.ListOptions{
			LabelSelector: labelSelector,
		})
		if err != nil {
			fmt.Fprintln(os.Stdout, util.State[3], err)
			os.Exit(3)
		}

		for _, pod := range podList.Items {
			if !(pod.Status.Phase == kApi.PodSucceeded || pod.Status.Phase == kApi.PodRunning) {
				objectInfoList = append(objectInfoList, &objectInfo{Name: pod.Name, Status: string(pod.Status.Phase), Namespace: pod.Namespace})
			}
		}
	}

	if len(objectInfoList) == 0 {
		fmt.Fprintln(os.Stdout, util.State[0], "All pods are Ready")
		os.Exit(0)
	} else {
		output := &serviceOutput{
			Objects: objectInfoList,
			Message: fmt.Sprintf("Found %d not running pods(s)", len(objectInfoList)),
		}
		outputByte, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			fmt.Fprintln(os.Stdout, util.State[3], err)
			os.Exit(3)
		}
		fmt.Fprintln(os.Stdout, util.State[0], string(outputByte))
		os.Exit(0)
	}
}

func NewCmd() *cobra.Command {
	var req request
	c := &cobra.Command{
		Use:     "pod_status",
		Short:   "Check Kubernetes Pod(s)",
		Example: "",

		Run: func(cmd *cobra.Command, args []string) {
			if req.objectType != "" {
				util.EnsureFlagsSet(cmd, "namespace", "object_name")
			}
			checkPodStatus(cmd, &req)

		},
	}
	c.Flags().StringVarP(&req.namespace, "namespace", "n", "", "Kubernetes namespace")
	c.Flags().StringVarP(&req.objectType, "object_type", "t", "", "Kubernetes Object Type")
	c.Flags().StringVarP(&req.objectName, "object_name", "N", "", "Kubernetes Object Name")
	return c
}
