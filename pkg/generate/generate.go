package generate

import (
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	var generateCmd = &cobra.Command{
		Use:   "generate",
		Short: "generate some thing",
		Args:  cobra.ExactValidArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
		},
	}

	generateCmd.AddCommand(newCmdKubeconfig())
	return generateCmd
}
