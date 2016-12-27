package main

import (
	_ "net/http/pprof"

	"github.com/appscode/errors"
	err_logger "github.com/appscode/errors/h/log"
	"github.com/appscode/go/flags"
	"github.com/appscode/go/hold"
	v "github.com/appscode/go/version"
	"github.com/appscode/log"
	logs "github.com/appscode/log/golog"
	"github.com/appscode/searchlight/cmd/searchlight/app"
	"github.com/appscode/searchlight/cmd/searchlight/app/options"
	"github.com/spf13/pflag"
)

var (
	Version         string
	VersionStrategy string
	Os              string
	Arch            string
	CommitHash      string
	GitBranch       string
	GitTag          string
	CommitTimestamp string
	BuildTimestamp  string
	BuildHost       string
	BuildHostOs     string
	BuildHostArch   string
)

func init() {
	v.Version.Version = Version
	v.Version.VersionStrategy = VersionStrategy
	v.Version.Os = Os
	v.Version.Arch = Arch
	v.Version.CommitHash = CommitHash
	v.Version.GitBranch = GitBranch
	v.Version.GitTag = GitTag
	v.Version.CommitTimestamp = CommitTimestamp
	v.Version.BuildTimestamp = BuildTimestamp
	v.Version.BuildHost = BuildHost
	v.Version.BuildHostOs = BuildHostOs
	v.Version.BuildHostArch = BuildHostArch
}

func main() {
	config := options.NewConfig()
	config.AddFlags(pflag.CommandLine)

	flags.InitFlags()
	logs.InitLogs()
	defer logs.FlushLogs()
	errors.Handlers.Add(err_logger.LogHandler{})
	flags.DumpAll()

	log.Infoln("Starting Searchlight Controller...")
	go app.Run(config)

	hold.Hold()
}
