package cli

import (
	"github.com/gizo-network/gizo/helpers"
	"github.com/gizo-network/gizo/p2p"
	"github.com/spf13/cobra"
)

func init() {
	workerCmd.Flags().IntVarP(&port, "port", "p", 9998, "port to run worker on")
}

var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Spin up a worker node",
	Run: func(cmd *cobra.Command, args []string) {
		helpers.Banner()
		w := p2p.NewWorker()
		w.Start()
	},
}
