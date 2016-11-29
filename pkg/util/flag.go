package util

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func EnsureFlagsSet(cmd *cobra.Command, name ...string) {
	for _, n := range name {
		flag := cmd.Flag(n)
		if flag == nil {
			continue
		}
		if !flag.Changed {
			fmt.Fprintln(os.Stdout, State[3], fmt.Sprintf("flag [%v] is required but not provided.\n", flag.Name))
			os.Exit(3)
		}
	}
}

func EnsureAlterableFlagsSet(cmd *cobra.Command, name ...string) {
	provided := false
	for _, n := range name {
		flag := cmd.Flag(n)
		if flag.Changed == true {
			provided = true
			break
		}
	}
	if provided == false {
		fmt.Fprintln(os.Stdout, State[3], fmt.Sprintf("one of the flags %v needs to be set.\n", name))
		os.Exit(3)
	}
}
