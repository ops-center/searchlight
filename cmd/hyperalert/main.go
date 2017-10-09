//go:generate stringer -type=State ../../pkg/icinga/types.go
package main

import (
	"os"

	logs "github.com/appscode/go/log/golog"
	_ "github.com/appscode/searchlight/client/scheme"
	"github.com/appscode/searchlight/plugins/hyperalert"
	_ "k8s.io/client-go/kubernetes/fake"
	_ "k8s.io/client-go/pkg/api/v1"
)

func main() {
	logs.InitLogs()
	defer logs.FlushLogs()
	if err := hyperalert.NewCmd().Execute(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
