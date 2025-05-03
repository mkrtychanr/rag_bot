package commands

import (
	"os"

	"github.com/mkrtychanr/rag_bot/internal/logger"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	SilenceUsage: true,
	RunE:         startCommandRun,
}

// Execute tries to execute commands and command line parameters.
func Execute() {
	rootCmd.HelpFunc()
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.PersistentFlags().String("config", "", "config file")
	rootCmd.PersistentFlags().String("cpuprofile", "", "write cpu profile to `file`")
	rootCmd.PersistentFlags().String("memprofile", "", "write memory profile to `file`")

	if err := rootCmd.Execute(); err != nil {
		logger.GetLogger().Err(err).Msg("cou")

		os.Exit(1)
	}
}
