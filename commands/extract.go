package commands

import (
	"github.com/spf13/cobra"
	"github.com/weblfe/gorar/pkg/extract"
)

func NewExtractCmd() *cobra.Command {
	var outputDir string
	cmd := &cobra.Command{
		Use:   "extract [archive]",
		Short: "Extract compressed files",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if outputDir == "" {
				outputDir = "."
			}
			return extract.Extract(args[0], outputDir) // 自动处理文件名冲突
		},
	}
	cmd.Flags().StringVarP(&outputDir, "output", "o", "", "Output directory")
	return cmd
}
