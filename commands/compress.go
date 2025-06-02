package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/weblfe/gorar/pkg/compress"
	"os"
	"path/filepath"
)

func NewCompressCmd() *cobra.Command {
	var format string
	cmd := &cobra.Command{
		Use:   "compress [path]",
		Short: "Compress files/folders",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			src := args[0]
			if _, err := os.Stat(src); os.IsNotExist(err) {
				return fmt.Errorf("path not exists: %s", src)
			}
			dest := fmt.Sprintf("%s.%s", filepath.Base(src), format)
			return compress.Compress(dest, src)
		},
	}
	cmd.Flags().StringVarP(&format, "format", "f", "zip", "Compression format (zip/tar/gz/xz/7z)")
	return cmd
}
