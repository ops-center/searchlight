package check_node_status

import (
	"fmt"
	"os"

	"github.com/appscode/searchlight/pkg/config"
	"github.com/appscode/searchlight/pkg/util"
	"github.com/spf13/cobra"
	kApi "k8s.io/kubernetes/pkg/api"
)

type request struct {
	name string
}

func checkNodeStatus(cmd *cobra.Command, req *request) {
	kubeClient, err := config.GetKubeClient()
	if err != nil {
		fmt.Fprintln(os.Stdout, util.State[3], err)
		os.Exit(3)
	}

	node, err := kubeClient.Nodes().Get(req.name)
	if err != nil {
		fmt.Fprintln(os.Stdout, util.State[3], err)
		os.Exit(3)
	}

	if node == nil {
		fmt.Fprintln(os.Stdout, util.State[2], "Node not found")
		os.Exit(2)
	}

	for _, condition := range node.Status.Conditions {
		if condition.Type == kApi.NodeReady {
			if condition.Status == kApi.ConditionFalse {
				fmt.Fprintln(os.Stdout, util.State[2], "Node is not Ready")
				os.Exit(2)
			}
		}
	}

	fmt.Fprintln(os.Stdout, util.State[0], "Node is Ready")
	os.Exit(0)
}

func NewCmd() *cobra.Command {
	var req request

	c := &cobra.Command{
		Use:     "node_status",
		Short:   "Check Kubernetes Node",
		Example: "",

		Run: func(cmd *cobra.Command, args []string) {
			util.EnsureFlagsSet(cmd, "name")
			checkNodeStatus(cmd, &req)
		},
	}

	c.Flags().StringVarP(&req.name, "name", "n", "", "Kubernetes node name")
	return c
}
