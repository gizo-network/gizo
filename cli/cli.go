package cli

import (
	"github.com/kpango/glg"
	"github.com/spf13/cobra"
)

var gizoCmd = &cobra.Command{
	Use:     "gizo [command]",
	Short:   "Job scheduling system build using blockchain",
	Long:    `Decentralized distributed system built using blockchain technology to provide a marketplace for users to trade their processing power in reward for ethereum`,
	Args:    cobra.MinimumNArgs(1),
	Version: "1.0.0",
}

func Execute() {
	gizoCmd.AddCommand(workerCmd, dispatcherCmd)
	if err := gizoCmd.Execute(); err != nil {
		glg.Fatal(err)
	}
}
