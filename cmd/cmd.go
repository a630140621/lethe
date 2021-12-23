package cmd

import (
	"github.com/a630140621/lethe/pkg/generate"
	"github.com/spf13/cobra"
)

func Execute() error {
	rootCmd := &cobra.Command{
		Use:   "lethe",
		Short: "Lethe is a tools for kubernetes",
		Long:  `Lethe now can generate kubeconfig from ServiceAccount`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
		},
	}

	rootCmd.AddCommand(generate.NewCmd())
	return rootCmd.Execute()
}
