package commands

import "github.com/spf13/cobra"

func New(version string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "gorar",
		Version: version,
		Short:   "Cross-platform file compression tool",
	}
	rootCmd.AddCommand(NewExtractCmd())
	rootCmd.AddCommand(NewCompressCmd())
	return rootCmd
}
