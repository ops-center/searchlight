package check_component_status

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/appscode/searchlight/pkg/client/k8s"
	"github.com/appscode/searchlight/util"
	"github.com/spf13/cobra"
	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/labels"
)

type objectInfo struct {
	Name   string `json:"name,omitempty"`
	Status string `json:"status,omitempty"`
}

type serviceOutput struct {
	Objects []*objectInfo `json:"objects,omitempty"`
	Message string        `json:"message,omitempty"`
}

func checkComponentStatus() {
	kubeClient, err := k8s.NewClient()
	if err != nil {
		fmt.Fprintln(os.Stdout, util.State[3], err)
		os.Exit(3)
	}

	components, err := kubeClient.Client.Core().
		ComponentStatuses().List(
		kapi.ListOptions{
			LabelSelector: labels.Everything(),
		},
	)
	if err != nil {
		fmt.Fprintln(os.Stdout, util.State[3], err)
		os.Exit(3)
	}

	objectInfoList := make([]*objectInfo, 0)
	for _, component := range components.Items {
		for _, condition := range component.Conditions {
			if condition.Type == kapi.ComponentHealthy && condition.Status == kapi.ConditionFalse {
				objectInfoList = append(objectInfoList,
					&objectInfo{
						Name:   component.Name,
						Status: "Unhealthy",
					},
				)
			}
		}
	}

	if len(objectInfoList) == 0 {
		fmt.Fprintln(os.Stdout, util.State[0], "All components are healthy")
		os.Exit(0)
	} else {
		output := &serviceOutput{
			Objects: objectInfoList,
			Message: fmt.Sprintf("%d unhealthy component(s)", len(objectInfoList)),
		}
		outputByte, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			fmt.Fprintln(os.Stdout, util.State[3], err)
			os.Exit(3)
		}
		fmt.Fprintln(os.Stdout, util.State[2], string(outputByte))
		os.Exit(2)
	}
}

func NewCmd() *cobra.Command {
	c := &cobra.Command{
		Use:     "check_component_status",
		Short:   "Check Kubernetes Component Status",
		Example: "",

		Run: func(cmd *cobra.Command, args []string) {
			checkComponentStatus()
		},
	}
	return c
}
