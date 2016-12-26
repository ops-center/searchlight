package icinga

import (
	"fmt"
	"github.com/appscode/go/io"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var configDataPath = os.ExpandEnv("$GOPATH") + "/src/github.com/appscode/searchlight/pkg/client/icinga/config.ini"

func TestNewClient(t *testing.T) {
	config, err := io.ReadINIConfig(configDataPath)
	assert.Nil(t, err)

	icingaConfig := &IcingaConfig{
		Endpoint: fmt.Sprintf("https://%v:5665/v1", config[IcingaURL]),
		CaCert:   nil,
	}
	icingaConfig.BasicAuth.Username = config[IcingaAPIUser]
	icingaConfig.BasicAuth.Username = config[IcingaAPIPass]

	icingaClient := NewClient(icingaConfig)
	resp := icingaClient.Check().Get([]string{}).Do()
	assert.Nil(t, resp.Err)
}
