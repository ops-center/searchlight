package check_node_status

import (
	"fmt"
	"os"
	"strings"

	"github.com/appscode/searchlight/pkg/config"
	"github.com/appscode/searchlight/pkg/util"
	"github.com/spf13/cobra"
	kApi "k8s.io/kubernetes/pkg/api"
)

type request struct {
	name string
}

func checkNodeStatus(req *request) {
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
	var host string
	c := &cobra.Command{
		Use:     "node_status",
		Short:   "Check Kubernetes Node",
		Example: "",

		Run: func(cmd *cobra.Command, args []string) {
			util.EnsureFlagsSet(cmd, "host")
			parts := strings.Split(host, "@")
			if len(parts) != 2 {
				fmt.Fprintln(os.Stdout, util.State[3], "Invalid icinga host.name")
				os.Exit(3)
			}

			req.name = parts[0]
			checkNodeStatus(&req)
		},
	}

	c.Flags().StringVarP(&host, "host", "H", "", "Icinga host name")
	return c
}
