package tracer

import "github.com/jpillora/go-ogle-analytics"

var GAKey string

func sendTrace(trace *trace) {
	if GAKey != "" {
		client, _ := ga.NewClient(GAKey)
		client.Send(ga.NewEvent(trace.namespace, trace.command).
			Label(trace.statusCode).
			Value(trace.elapsed.Nanoseconds()))
	}
}
