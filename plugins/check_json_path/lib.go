package check_json_path

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"bytes"
	"os/exec"
	"reflect"

	"github.com/appscode/go-httpclient"
	"github.com/appscode/searchlight/pkg/config"
	"github.com/appscode/searchlight/pkg/util"
	"github.com/spf13/cobra"
)

type request struct {
	url      string
	query    string
	secret   string
	warning  string
	critical string
}

type AuthInfo struct {
	CertificateAuthorityData []byte `json:"certificate-authority-data,omitempty"`
	ClientCertificateData    []byte `json:"client-certificate-data,omitempty"`
	Token                    string `json:"token,omitempty"`
	Username                 string `json:"username,omitempty"`
	Password                 string `json:"password,omitempty"`
}

type JQ struct {
	J string `json:"j"`
	Q string `json:"q"`
}

const (
	auth = "auth"
)

func getData(secretName, urlPath string) string {
	httpClient := httpclient.Default().WithBaseURL(urlPath)
	if secretName != "" {
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

		secretData := new(AuthInfo)
		if data, found := secret.Data[auth]; found {
			if err := json.Unmarshal(data, secretData); err != nil {
				fmt.Fprintln(os.Stdout, util.State[3], err)
				os.Exit(3)
			}
		}
		httpClient.WithBearerToken(secretData.Token)
		httpClient.WithBasicAuth(secretData.Username, secretData.Password)
		if secretData.CertificateAuthorityData != nil {
			httpClient.WithTLSConfig(secretData.ClientCertificateData, secretData.CertificateAuthorityData)
		}
	}

	var respJson interface{}
	resp, err := httpClient.Call("GET", "", nil, &respJson, true)
	if err != nil {
		fmt.Fprintln(os.Stdout, util.State[3], err)
		os.Exit(3)
	}
	if resp.StatusCode != 200 {
		fmt.Fprintln(os.Stdout, util.State[3], fmt.Sprintf("Invalid status_code %v", resp.StatusCode))
		os.Exit(3)
	}

	data, err := json.Marshal(respJson)
	if err != nil {
		fmt.Fprintln(os.Stdout, util.State[3], err)
		os.Exit(3)
	}

	return string(data)
}

func (j *JQ) eval() (res interface{}) {
	cmd := exec.Command("jq", j.Q)
	cmd.Stdin = bytes.NewBufferString(j.J)
	cmdOut, err := cmd.Output()
	if err != nil {
		fmt.Fprintln(os.Stdout, util.State[3], err)
		os.Exit(3)

	}
	if err = json.Unmarshal(cmdOut, &res); err != nil {
		fmt.Fprintln(os.Stdout, util.State[3], err)
		os.Exit(3)
	}
	return
}

func checkResult(evalDataString, checkQuery string) bool {
	jqData := &JQ{
		J: string(evalDataString),
		Q: checkQuery,
	}
	result := jqData.eval()
	if reflect.TypeOf(result).Kind() != reflect.Bool {
		fmt.Fprintln(os.Stdout, util.State[3], fmt.Sprintf("Invalid check query: %v", checkQuery))
		os.Exit(3)
	}
	return result.(bool)
}

func checkJsonPath(cmd *cobra.Command, req *request) {
	jsonData := getData(req.secret, req.url)
	jqData := &JQ{
		J: jsonData,
		Q: req.query,
	}

	evalData := jqData.eval()
	if evalData == nil {
		fmt.Fprintln(os.Stdout, util.State[3], "Invalid query. No data found")
		os.Exit(3)
	}

	evalDataByte, err := json.Marshal(evalData)
	if err != nil {
		fmt.Fprintln(os.Stdout, util.State[3], err)
		os.Exit(3)

	}

	evalDataString := string(evalDataByte)
	if req.critical != "" {
		if checkResult(evalDataString, req.critical) {
			fmt.Fprintln(os.Stdout, util.State[2], fmt.Sprintf("%v", req.critical))
			os.Exit(2)
		}
	}
	if req.warning != "" {
		if checkResult(evalDataString, req.warning) {
			fmt.Fprintln(os.Stdout, util.State[1], fmt.Sprintf("%v", req.warning))
			os.Exit(1)
		}
	}

	fmt.Fprintln(os.Stdout, util.State[0])
	os.Exit(0)
}

func NewCmd() *cobra.Command {
	var req request

	c := &cobra.Command{
		Use:     "json_path",
		Short:   "Check Json Object",
		Example: "",

		Run: func(cmd *cobra.Command, args []string) {
			util.EnsureFlagsSet(cmd, "url", "query")
			util.EnsureAlterableFlagsSet(cmd, "warning", "critical")
			checkJsonPath(cmd, &req)
		},
	}

	c.Flags().StringVarP(&req.url, "url", "u", "", "URL to get data")
	c.Flags().StringVarP(&req.query, "query", "q", "", `JQ query`)
	c.Flags().StringVarP(&req.secret, "secret", "s", "", `Kubernetes secret name`)
	c.Flags().StringVarP(&req.warning, "warning", "w", "", `Warning JQ query which returns [true/false]`)
	c.Flags().StringVarP(&req.critical, "critical", "c", "", `Critical JQ query which returns [true/false]`)
	return c
}
