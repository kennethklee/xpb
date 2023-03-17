package cmd

import (
	"fmt"
	"os/exec"

	"github.com/kennethklee/xpb"
	"github.com/spf13/cobra"
)

func NewBuildCommand() *cobra.Command {
	var withPlugins []string
	cmd := &cobra.Command{
		Use:   "build <version>",
		Short: "Build a custom pocketbase",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := exec.LookPath("go")
			if err != nil {
				return fmt.Errorf("go toolchain not found: %w", err)
			}

			fmt.Println("args", args)

			builder, err := xpb.NewBuilder(args[0], withPlugins...)
			if err != nil {
				return err
			}
			defer builder.Close()

			if err := builder.Compile(args[1:]...); err != nil {
				fmt.Println("[ERROR]", "Failed to compile:", err)
				fmt.Println("[INFO]", "You can find the project directory at:", builder.TempProjectDir)
				fmt.Println("Press any key to continue...")
				fmt.Scanln()

				return err
			}

			return nil
		},
	}
	cmd.Flags().StringArrayVar(&withPlugins, "with", []string{}, "include plugin  (format: module[@version][=replacement])")
	return cmd
}
