package check_pod_exists

import (
	"fmt"
	"os"

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
	count      int
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

func checkPodExists(req *request, checkCount bool) {
	kubeClient, err := config.GetKubeClient()
	if err != nil {
		fmt.Fprintln(os.Stdout, util.State[3], err)
		os.Exit(3)
	}

	total_pod := 0
	if req.objectType == config.TypePods {
		pod, err := kubeClient.Pods(req.namespace).Get(req.objectName)
		if err != nil {
			fmt.Fprintln(os.Stdout, util.State[3], err)
			os.Exit(3)
		}
		if pod != nil {
			total_pod = 1
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

		total_pod = len(podList.Items)
	}

	if checkCount {
		if req.count != total_pod {
			fmt.Fprintln(os.Stdout, util.State[2], fmt.Sprintf("Found %d pod(s) instead of %d", total_pod, req.count))
			os.Exit(2)
		} else {
			fmt.Fprintln(os.Stdout, util.State[0], "Found all pods")
			os.Exit(0)
		}
	} else {
		if total_pod == 0 {
			fmt.Fprintln(os.Stdout, util.State[2], "No pod found")
			os.Exit(2)
		} else {
			fmt.Fprintln(os.Stdout, util.State[0], fmt.Sprintf("Found %d pods(s)", total_pod))
			os.Exit(0)
		}
	}
}

func NewCmd() *cobra.Command {
	var req request
	c := &cobra.Command{
		Use:     "pod_exists",
		Short:   "Check Kubernetes Pod(s)",
		Example: "",

		Run: func(cmd *cobra.Command, args []string) {
			if req.objectType != "" {
				util.EnsureFlagsSet(cmd, "namespace", "object_name")
			}
			checkCount := cmd.Flag("count").Changed
			checkPodExists(&req, checkCount)
		},
	}
	c.Flags().StringVarP(&req.namespace, "namespace", "n", "", "Kubernetes namespace")
	c.Flags().StringVarP(&req.objectType, "object_type", "t", "", "Kubernetes Object Type")
	c.Flags().StringVarP(&req.objectName, "object_name", "N", "", "Kubernetes Object Name")
	c.Flags().IntVarP(&req.count, "count", "c", 0, "Number of Kubernetes Node")
	return c
}
