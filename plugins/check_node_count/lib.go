package check_node_count

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
	count int
}

func checkNodeStatus(req *request) {
	kubeClient, err := config.GetKubeClient()
	if err != nil {
		fmt.Fprintln(os.Stdout, util.State[3], err)
		os.Exit(3)
	}

	nodeList, err := kubeClient.Nodes().List(kApi.ListOptions{
		LabelSelector: labels.Everything(),
	})
	if err != nil {
		fmt.Fprintln(os.Stdout, util.State[3], err)
		os.Exit(3)
	}

	if len(nodeList.Items) == req.count {
		fmt.Fprintln(os.Stdout, util.State[0], "Found all nodes")
		os.Exit(0)
	} else {
		fmt.Fprintln(os.Stdout, util.State[2], fmt.Sprintf("Found %d node(s) instead of %d", len(nodeList.Items), req.count))
		os.Exit(2)
	}
}

func NewCmd() *cobra.Command {
	var req request

	c := &cobra.Command{
		Use:     "node_count",
		Short:   "Count Kubernetes Nodes",
		Example: "",

		Run: func(cmd *cobra.Command, args []string) {
			util.EnsureFlagsSet(cmd, "count")
			checkNodeStatus(&req)
		},
	}

	c.Flags().IntVarP(&req.count, "count", "c", 0, "Number of Kubernetes Node")
	return c
}
