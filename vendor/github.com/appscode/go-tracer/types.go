package tracer

import (
	"time"

	"github.com/spf13/cobra"
)

type trace struct {
	namespace  string
	enable     bool
	command    string
	statusCode string

	startTime time.Time
	endTime   time.Time
	elapsed   time.Duration
}

func New() *trace {
	return &trace{
		statusCode: "0",
		startTime:  time.Now(),
		enable:     false,
	}
}

func (t *trace) WithNamespace(ns string) *trace {
	t.namespace = ns
	return t
}

func (t *trace) WithStatus(code string) *trace {
	t.statusCode = code
	return t
}

func (t *trace) Enable(e bool) *trace {
	t.enable = e
	return t
}

func (t *trace) Done(cmd *cobra.Command) {
	t.elapsed = time.Since(t.startTime)
	t.endTime = time.Now()
	t.command = cmd.CommandPath()
	if t.enable {
		sendTrace(t)
	}
}
