package check_node_count

import (
	"fmt"
	"os"

	flags "github.com/appscode/go-flags"
	"github.com/appscode/searchlight/pkg/config"
	"github.com/appscode/searchlight/util"
	"github.com/spf13/cobra"
	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/labels"
)

type request struct {
	count int
}

func checkNodeStatus(req *request) {
	kubeClient, err := config.NewKubeClient()
	if err != nil {
		fmt.Fprintln(os.Stdout, util.State[3], err)
		os.Exit(3)
	}

	nodeList, err := kubeClient.Client.Core().
		Nodes().List(
		kapi.ListOptions{
			LabelSelector: labels.Everything(),
		},
	)
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
		Use:     "check_node_count",
		Short:   "Count Kubernetes Nodes",
		Example: "",

		Run: func(cmd *cobra.Command, args []string) {
			flags.EnsureRequiredFlags(cmd, "count")
			checkNodeStatus(&req)
		},
	}

	c.Flags().IntVarP(&req.count, "count", "c", 0, "Number of expected Kubernetes Node")
	return c
}
