package cmd

import (
	"github.com/spf13/cobra"

	"github.com/0xcb9ff9/goleveldb-analyze/v2/flags"
)

func RegisterLeveldbPathFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().String(
		flags.LeveldbPathFlag,
		"",
		"open leveldb path",
	)
}
