package cmd

import (
	"fmt"
	"os"

	"github.com/0xcb9ff9/goleveldb-analyze/v2/cmd/kvsize"
	"github.com/0xcb9ff9/goleveldb-analyze/v2/cmd/stats"

	"github.com/spf13/cobra"
)

type RootCommand struct {
	baseCmd *cobra.Command
}

func (root *RootCommand) registerSubCommands() {
	root.baseCmd.AddCommand(
		kvsize.GetCommand(),
		stats.GetCommand(),
	)
}

func (root *RootCommand) Execute() {
	if err := root.baseCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}
}

func NewRootCommand() *RootCommand {
	rootCommand := &RootCommand{
		baseCmd: &cobra.Command{
			Use:   "goleveldb-analyze",
			Short: "goleveldb-analyze is a simple tool to analyze leveldb.",
		},
	}

	RegisterLeveldbPathFlag(rootCommand.baseCmd)

	rootCommand.registerSubCommands()

	return rootCommand
}
