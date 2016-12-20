package tracer

import (
	"time"

	"github.com/spf13/cobra"
)

// DefaultTrace is the default trace used by tracer.
var DefaultTrace *trace = &trace{
	statusCode: "0",
	enable:     false,
}

func Init() {
	DefaultTrace.startTime = time.Now()
}

func SetNamespace(ns string) {
	DefaultTrace.namespace = ns
}

func SetStatus(code string) {
	DefaultTrace.statusCode = code
}

func Enable(e bool) {
	DefaultTrace.enable = e
}

func Done(cmd *cobra.Command) {
	DefaultTrace.Done(cmd)
}
