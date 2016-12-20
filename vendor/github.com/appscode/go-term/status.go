package term

import (
	"os"

	"github.com/appscode/api"
	"github.com/appscode/api/dtypes"
	tracer "github.com/appscode/go-tracer"
	"github.com/appscode/log"
)

func PrintStatus(resp api.Response) {
	s := resp.GetStatus()
	tracer.SetStatus(s.Code)
	if s.IsOK() {
		Successln(s.Message)
	} else {
		Errorln(s.Message)
	}
	parseHelp(s)
	logDetails(s)
	if !s.IsOK() {
		os.Exit(1)
	}
}

func parseHelp(status *dtypes.Status) {
	if status.Help != nil {
		if len(status.Help.Description) >= 1 {
			Infoln(status.Help.Description)
		}
		if len(status.Help.Url) >= 1 {
			Infoln("See here for details:", status.Help.Url)
		}
	}
}

func logDetails(status *dtypes.Status) {
	for _, d := range status.Details {
		log.Infoln(d.TypeUrl, "status:", status.Code, d.String())
	}
}
