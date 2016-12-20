package check_prometheus_metric

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/appscode/searchlight/util"
	"github.com/prometheus/client_golang/api/prometheus"
	"github.com/prometheus/common/model"
	"github.com/spf13/cobra"
)

type request struct {
	host        string
	query       string
	metric_name string
	method      string
	warning     int64
	critical    int64
	accept_nan  bool
}

func checkResult(method string, valueToCheck, result int64) bool {
	switch method {
	case "gt":
		if result > valueToCheck {
			return true
		}
	case "ge":
		if result >= valueToCheck {
			return true
		}
	case "lt":
		if result < valueToCheck {
			return true
		}
	case "le":
		if result <= valueToCheck {
			return true
		}
	case "eq":
		if result == valueToCheck {
			return true
		}
	case "ne":
		if result != valueToCheck {
			return true
		}
	}
	return false
}

func checkPrometheusMetric(req *request) {
	config := prometheus.Config{
		Address: req.host,
	}
	client, err := prometheus.New(config)
	if err != nil {
		fmt.Fprintln(os.Stdout, util.State[3], err)
		os.Exit(3)
	}
	queryApi := prometheus.NewQueryAPI(client)

	result, err := queryApi.Query(context.Background(), req.query, time.Now())
	if err != nil {
		fmt.Fprintln(os.Stdout, util.State[3], err)
		os.Exit(3)
	}

	vector := result.(model.Vector)

	if len(vector) == 0 {
		if req.accept_nan {
			fmt.Fprintln(os.Stdout, util.State[0], errors.New("NaN"))
			os.Exit(0)
		}
		fmt.Fprintln(os.Stdout, util.State[3], errors.New("No data found"))
		os.Exit(3)
	} else if len(vector) > 1 {
		fmt.Fprintln(os.Stdout, util.State[3], errors.New("Invalid query.\nQuery should return single float or int"))
		os.Exit(3)
	}

	value := vector[0].Value.String()
	val, err := strconv.ParseFloat(value, 64)
	if err != nil {
		fmt.Fprintln(os.Stdout, util.State[3], err)
		os.Exit(3)
	}

	valInt64 := int64(val)
	outputStr := fmt.Sprintf("%v", valInt64)
	if req.metric_name != "" {
		outputStr = fmt.Sprintf("%v is %v", req.metric_name, valInt64)
	}

	if checkResult(req.method, req.critical, valInt64) {
		fmt.Fprintln(os.Stdout, util.State[2], outputStr)
		os.Exit(2)
	}
	if checkResult(req.method, req.warning, valInt64) {
		fmt.Fprintln(os.Stdout, util.State[1], outputStr)
		os.Exit(1)
	}

	fmt.Fprintln(os.Stdout, util.State[0], outputStr)
	os.Exit(0)
	return
}

func NewCmd() *cobra.Command {
	var req request

	c := &cobra.Command{
		Use:     "check_prometheus_metric",
		Short:   "Check prometheus metric",
		Example: "",

		Run: func(cmd *cobra.Command, args []string) {
			if req.host == "" {
				fmt.Fprintln(os.Stdout, util.State[3], errors.New("No prometheus host found"))
				os.Exit(3)
			}
			util.EnsureFlagsSet(cmd, "query", "warning", "critical")
			checkPrometheusMetric(&req)
		},
	}

	c.Flags().StringVarP(&req.host, "prom_host", "H", "", "URL of Prometheus host to query")
	c.Flags().StringVarP(&req.query, "query", "q", "", "Prometheus query that returns a float or int")
	c.Flags().StringVarP(&req.metric_name, "metric_name", "n", "", "A name for the metric being checked")
	c.Flags().StringVarP(&req.method, "method", "m", "ge", `Comparison method, one of gt, ge, lt, le, eq, ne
	(defaults to ge unless otherwise specified)`)
	c.Flags().BoolVarP(&req.accept_nan, "accept_nan", "O", false, `Accept NaN as an "OK" result`)
	c.Flags().Int64VarP(&req.warning, "warning", "w", 0, "Warning level value (must be zero or positive)")
	c.Flags().Int64VarP(&req.critical, "critical", "c", 0, "Critical level value (must be zero or positive)")
	return c
}
