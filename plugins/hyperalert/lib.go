package hyperalert

import (
	"flag"

	v "github.com/appscode/go/version"
	"github.com/appscode/searchlight/client/clientset/versioned/scheme"
	"github.com/appscode/searchlight/plugins"
	"github.com/appscode/searchlight/plugins/analytics_id"
	"github.com/appscode/searchlight/plugins/check_ca_cert"
	"github.com/appscode/searchlight/plugins/check_cert"
	"github.com/appscode/searchlight/plugins/check_component_status"
	"github.com/appscode/searchlight/plugins/check_env"
	"github.com/appscode/searchlight/plugins/check_event"
	"github.com/appscode/searchlight/plugins/check_json_path"
	"github.com/appscode/searchlight/plugins/check_node_exists"
	"github.com/appscode/searchlight/plugins/check_node_status"
	"github.com/appscode/searchlight/plugins/check_pod_exec"
	"github.com/appscode/searchlight/plugins/check_pod_exists"
	"github.com/appscode/searchlight/plugins/check_pod_status"
	"github.com/appscode/searchlight/plugins/check_volume"
	"github.com/appscode/searchlight/plugins/check_webhook"
	"github.com/appscode/searchlight/plugins/notifier"
	"github.com/spf13/cobra"
	clientsetscheme "k8s.io/client-go/kubernetes/scheme"
	"kmodules.xyz/client-go/logs"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hyperalert",
		Short: "AppsCode Icinga2 plugin",
		PersistentPreRun: func(c *cobra.Command, args []string) {
			scheme.AddToScheme(clientsetscheme.Scheme)
		},
		Run: func(c *cobra.Command, args []string) {
			c.Help()
		},
	}
	cmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)
	logs.ParseFlags()
	cmd.PersistentFlags().String(plugins.FlagKubeConfig, "", "Path to kubeconfig file with authorization information (the master location is set by the master flag).")
	cmd.PersistentFlags().String(plugins.FlagKubeConfigContext, "", "Use the context in kubeconfig")
	cmd.PersistentFlags().Int(plugins.FlagCheckInterval, 30, "Icinga check_interval in second. [Format: 30, 300]")

	// CheckCluster
	cmd.AddCommand(check_component_status.NewCmd())
	cmd.AddCommand(check_json_path.NewCmd())
	cmd.AddCommand(check_node_exists.NewCmd())
	cmd.AddCommand(check_pod_exists.NewCmd())
	cmd.AddCommand(check_event.NewCmd())
	cmd.AddCommand(check_ca_cert.NewCmd())
	cmd.AddCommand(check_cert.NewCmd())
	cmd.AddCommand(check_env.NewCmd())
	cmd.AddCommand(check_webhook.NewCmd())

	// CheckNode
	cmd.AddCommand(check_node_status.NewCmd())

	// CheckPod
	cmd.AddCommand(check_pod_status.NewCmd())
	cmd.AddCommand(check_pod_exec.NewCmd())

	// Combined
	cmd.AddCommand(check_volume.NewCmd())

	// Notifier
	cmd.AddCommand(notifier.NewCmd())

	cmd.AddCommand(analytics_id.NewCmd())
	cmd.AddCommand(v.NewCmdVersion())

	return cmd
}
