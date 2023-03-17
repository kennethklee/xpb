package cmd

import "github.com/spf13/cobra"

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "xpocketbase",
		Short: "xpocketbase is a tool for building customized pocketbase",
	}
	rootCmd.AddCommand(NewBuildCommand())
	return rootCmd
}
