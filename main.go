package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/weblfe/gorar/commands"
)

var (
	version = "0.0.1"
)

func main() {
	rootCmd := &cobra.Command{
		Use:     "gorar",
		Version: version,
		Short:   "Cross-platform file compression tool",
	}
	rootCmd.AddCommand(commands.NewExtractCmd())
	rootCmd.AddCommand(commands.NewCompressCmd())
	if err := rootCmd.Execute(); err != nil {
		log.WithField("error", err).Errorln()
	}
}
