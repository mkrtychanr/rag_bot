// nolint
package commands

import (
	"fmt"

	"github.com/mkrtychanr/rag_bot/internal/common/version"
	"github.com/spf13/cobra"
)

var (
	svs        = "RAG Bot"
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Display application version",
		Run: func(c *cobra.Command, args []string) {
			fmt.Println(version.BuildVersionString(svs))
		},
	}
)
