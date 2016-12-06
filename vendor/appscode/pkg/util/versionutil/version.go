package versionutil

import (
	"github.com/appscode/api/version"
	"github.com/spf13/cobra"
)

var Version version.Version

func NewCmdVersion() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Prints binary version number.",
		Long:  `Prints binary version number.`,

		Run: func(cmd *cobra.Command, args []string) {
			Version.Print()
		},
	}
}
