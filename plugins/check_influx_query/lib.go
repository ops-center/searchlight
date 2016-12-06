package check_influx_query

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/Knetic/govaluate"
	"github.com/appscode/searchlight/pkg/config"
	"github.com/appscode/searchlight/pkg/util"
	influxdb "github.com/influxdata/influxdb/client"
	"github.com/spf13/cobra"
	ini "github.com/vaughan0/go-ini"
)

const (
	admin string = ".admin"

	influxDBHost     string = "INFLUX_HOST"
	influxDBDatabase string = "INFLUX_DB"
	influxDBReadUser string = "INFLUX_READ_USER"
	influxDBReadPass string = "INFLUX_READ_PASSWORD"

	influxDBDefaultDatabase string = "k8s"
	influxDBHostPort               = 8086
)

type request struct {
	host          string
	a, b, c, d, e string
	r             string
	warning       string
	critical      string
	secret        string
}

type authInfo struct {
	host     string
	database string
	username string
	password string
}

func trunc(val float64) interface{} {
	intData := int64(val * 1000)
	return float64(intData) / 1000.0
}

func getInfluxDBSecretData(secretName string) *authInfo {
	kubeClient, err := config.GetKubeClient()
	if err != nil {
		fmt.Fprintln(os.Stdout, util.State[3], err)
		os.Exit(3)
	}

	parts := strings.Split(secretName, ".")
	name := parts[0]
	namespace := "default"
	if len(parts) > 1 {
		namespace = parts[1]
	}

	secret, err := kubeClient.Secrets(namespace).Get(name)
	if err != nil {
		fmt.Fprintln(os.Stdout, util.State[3], err)
		os.Exit(3)
	}

	authData := new(authInfo)
	if data, found := secret.Data[admin]; found {
		dataReader := strings.NewReader(string(data))
		secretData, err := ini.Load(dataReader)
		if err != nil {
			fmt.Fprintln(os.Stdout, util.State[3], err)
			os.Exit(3)
		}

		if host, found := secretData.Get("", influxDBHost); found {
			authData.host = fmt.Sprintf("%s:%d", host, influxDBHostPort)
		}

		if authData.database, found = secretData.Get("", influxDBDatabase); !found {
			authData.database = influxDBDefaultDatabase
		}
		if authData.username, found = secretData.Get("", influxDBReadUser); !found {
			fmt.Fprintln(os.Stdout, util.State[3], errors.New("No InfluxDB read user found"))
			os.Exit(3)
		}
		if authData.password, found = secretData.Get("", influxDBReadPass); !found {
			fmt.Fprintln(os.Stdout, util.State[3], errors.New("No InfluxDB read password found"))
			os.Exit(3)
		}
		return authData
	}
	fmt.Fprintln(os.Stdout, util.State[3], errors.New("Invalid InfluxDB secret"))
	os.Exit(3)
	return nil
}

func getInfluxDBClient(authData *authInfo) *influxdb.Client {
	config := &influxdb.Config{
		URL: url.URL{
			Scheme: "http",
			Host:   authData.host,
		},
		Username: authData.username,
		Password: authData.password,
	}
	client, err := influxdb.NewClient(*config)
	if err != nil {
		fmt.Fprintln(os.Stdout, util.State[3], err)
		os.Exit(3)
	}
	return client
}

func getInfluxdbData(con *influxdb.Client, command, db, queryName string) float64 {
	q := influxdb.Query{
		Command:  command,
		Database: db,
	}
	res, err := con.Query(q)
	if err != nil {
		fmt.Fprintln(os.Stdout, util.State[3], err)
		os.Exit(3)
	}

	if len(res.Results[0].Series) == 0 {
		fmt.Fprintln(os.Stdout, util.State[3], errors.New(fmt.Sprint("Value not found for query: ", queryName)))
		os.Exit(3)
	}
	data, err := strconv.ParseFloat((string(res.Results[0].Series[0].Values[0][1].(json.Number))), 64)
	if err != nil {
		fmt.Fprintln(os.Stdout, util.State[3], err)
		os.Exit(3)
	}
	return data
}

func getValue(con *influxdb.Client, db string, req *request) map[string]interface{} {

	valMap := make(map[string]interface{})

	defer func() {
		if e := recover(); e != nil {
			fmt.Fprintln(os.Stdout, util.State[3], e)
			os.Exit(3)
		}
	}()

	if req.a != "" {
		data := getInfluxdbData(con, req.a, db, "A")
		valMap["A"] = data
	}

	if req.b != "" {
		data := getInfluxdbData(con, req.b, db, "B")
		valMap["B"] = data
	}

	if req.c != "" {
		data := getInfluxdbData(con, req.c, db, "C")
		valMap["C"] = data
	}

	if req.d != "" {
		data := getInfluxdbData(con, req.d, db, "D")
		valMap["D"] = data
	}

	if req.e != "" {
		data := getInfluxdbData(con, req.e, db, "E")
		valMap["E"] = data
	}

	return valMap
}

func checkResult(checkQuery string, valueMap map[string]interface{}) bool {
	expr, err := govaluate.NewEvaluableExpression(checkQuery)
	if err != nil {
		fmt.Fprintln(os.Stdout, util.State[3], err)
		os.Exit(3)
	}

	res, err := expr.Evaluate(valueMap)
	if err != nil {
		fmt.Fprintln(os.Stdout, util.State[3], err)
		os.Exit(3)
	}

	return res.(bool)
}

func checkInfluxQuery(req *request) {
	authData := getInfluxDBSecretData(req.secret)
	if req.host != "" {
		authData.host = req.host
	}
	if authData.host == "" {
		fmt.Fprintln(os.Stdout, util.State[3], errors.New("No InfluxDB host found"))
		os.Exit(3)
	}
	client := getInfluxDBClient(authData)

	valMap := getValue(client, authData.database, req)

	expression, err := govaluate.NewEvaluableExpression(req.r)
	if err != nil {
		fmt.Fprintln(os.Stdout, util.State[3], err)
		os.Exit(3)
	}

	if valMap["R"], err = expression.Evaluate(valMap); err != nil {
		fmt.Fprintln(os.Stdout, util.State[3], err)
		os.Exit(3)
	}
	valMap["R"] = trunc(valMap["R"].(float64))

	if req.critical != "" {
		checkResult(req.critical, valMap)
	}

	if req.warning != "" {
		checkResult(req.warning, valMap)
	}

	fmt.Fprintln(os.Stdout, util.State[0], "Fine")
	os.Exit(0)
}

func NewCmd() *cobra.Command {
	var req request

	c := &cobra.Command{
		Use:     "influx_query",
		Short:   "Check InfluxDB Query Data",
		Example: "",

		Run: func(cmd *cobra.Command, args []string) {
			util.EnsureFlagsSet(cmd, "secret", "R")
			util.EnsureAlterableFlagsSet(cmd, "A", "B", "C", "D", "E")
			util.EnsureAlterableFlagsSet(cmd, "warning", "critical")
			checkInfluxQuery(&req)
		},
	}

	c.Flags().StringVarP(&req.host, "influx_host", "H", "", "URL of InfluxDB host to query")
	c.Flags().StringVarP(&req.secret, "secret", "s", "", `Kubernetes secret name`)
	c.Flags().StringVar(&req.a, "A", "", "InfluxDB query A")
	c.Flags().StringVar(&req.b, "B", "", "InfluxDB query B")
	c.Flags().StringVar(&req.c, "C", "", "InfluxDB query C")
	c.Flags().StringVar(&req.d, "D", "", "InfluxDB query D")
	c.Flags().StringVar(&req.e, "E", "", "InfluxDB query E")
	c.Flags().StringVar(&req.r, "R", "", `Equation to evaluate result`)
	c.Flags().StringVarP(&req.warning, "warning", "w", "", `Warning query which returns [true/false]`)
	c.Flags().StringVarP(&req.critical, "critical", "c", "", `Critical query which returns [true/false]`)
	return c
}
