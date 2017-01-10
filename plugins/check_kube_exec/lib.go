package check_kube_exec

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/appscode/go/flags"
	"github.com/appscode/searchlight/pkg/client/k8s"
	"github.com/appscode/searchlight/util"
	"github.com/spf13/cobra"
	kapi "k8s.io/kubernetes/pkg/api"
	clientset "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
	"k8s.io/kubernetes/pkg/client/unversioned/remotecommand"
	remotecommandserver "k8s.io/kubernetes/pkg/kubelet/server/remotecommand"
	utilexec "k8s.io/kubernetes/pkg/util/exec"
)

type Writer struct {
	Str []string
}

func (w *Writer) Write(p []byte) (n int, err error) {
	str := string(p)
	if len(str) > 0 {
		w.Str = append(w.Str, str)
	}
	return len(str), nil
}

func newStringReader(ss []string) io.Reader {
	formattedString := strings.Join(ss, "\n")
	reader := strings.NewReader(formattedString)
	return reader
}

func checkKubeExec(req *request) {
	kubeConfig, err := k8s.GetKubeConfig()
	if err != nil {
		fmt.Fprintln(os.Stdout, util.State[3], err)
		os.Exit(3)
	}

	kubeClient, err := clientset.NewForConfig(kubeConfig)
	if err != nil {
		fmt.Fprintln(os.Stdout, util.State[3], err)
		os.Exit(3)
	}

	pod, err := kubeClient.Core().Pods(req.namespace).Get(req.pod)
	if err != nil {
		fmt.Fprintln(os.Stdout, util.State[3], err)
		os.Exit(3)
	}

	foundContainer := false
	if req.container != "" {
		for _, container := range pod.Spec.Containers {
			if container.Name == req.container {
				foundContainer = true
				break
			}
		}
	} else {
		foundContainer = true
	}

	if !foundContainer {
		fmt.Fprintln(os.Stdout, util.State[3], fmt.Sprintf(`Container "%v" not found`, req.container))
		os.Exit(3)
	}

	execRequest := kubeClient.Core().RESTClient().Post().
		Resource("pods").
		Name(req.pod).
		Namespace(req.namespace).
		SubResource("exec").
		Param("container", req.container)

	execRequest.VersionedParams(&kapi.PodExecOptions{
		Container: req.container,
		Command:   []string{req.command},
		Stdin:     true,
		Stdout:    false,
		Stderr:    false,
		TTY:       false,
	}, kapi.ParameterCodec)

	exec, err := remotecommand.NewExecutor(kubeConfig, "POST", execRequest.URL())
	if err != nil {
		fmt.Fprintln(os.Stdout, util.State[3], err)
		os.Exit(3)
	}

	stdIn := newStringReader([]string{"-c", req.arg})
	stdOut := new(Writer)
	stdErr := new(Writer)

	err = exec.Stream(remotecommand.StreamOptions{
		SupportedProtocols: remotecommandserver.SupportedStreamingProtocols,
		Stdin:              stdIn,
		Stdout:             stdOut,
		Stderr:             stdErr,
		Tty:                false,
	})

	var exitCode int
	if err == nil {
		exitCode = 0
	} else {
		if exitErr, ok := err.(utilexec.ExitError); ok && exitErr.Exited() {
			exitCode = exitErr.ExitStatus()
		} else {
			fmt.Fprintln(os.Stdout, util.State[3], "Failed to find exit code.")
			os.Exit(3)
		}
	}

	output := fmt.Sprintf("Exit Code: %v", exitCode)
	if exitCode != 0 {
		exitCode = 2
	}

	fmt.Fprintln(os.Stdout, util.State[exitCode], output)
	os.Exit(exitCode)
}

type request struct {
	pod       string
	container string
	namespace string
	command   string
	arg       string
}

func NewCmd() *cobra.Command {
	var req request
	var host string
	c := &cobra.Command{
		Use:     "check_kube_exec",
		Short:   "Check exit code of exec command on kubernetes container",
		Example: "",

		Run: func(cmd *cobra.Command, args []string) {
			flags.EnsureRequiredFlags(cmd, "host", "arg")

			parts := strings.Split(host, "@")
			if len(parts) != 2 {
				fmt.Fprintln(os.Stdout, util.State[3], "Invalid icinga host.name")
				os.Exit(3)
			}
			req.pod = parts[0]
			req.namespace = parts[1]
			checkKubeExec(&req)
		},
	}

	c.Flags().StringVarP(&host, "host", "H", "", "Icinga host name")
	c.Flags().StringVarP(&req.container, "container", "C", "", "Container name in specified pod")
	c.Flags().StringVarP(&req.command, "cmd", "c", "/bin/sh", "Exec command. [Default: /bin/sh]")
	c.Flags().StringVarP(&req.arg, "argv", "a", "", "Arguments for exec command. [Format: 'arg; arg; arg']")
	return c
}
