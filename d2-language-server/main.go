package main

import (
	"os"

	"github.com/ram02z/d2-language-server/internal/lsp"
	"github.com/spf13/cobra"
	"github.com/tliron/commonlog"
	_ "github.com/tliron/commonlog/simple"
	serverpkg "github.com/tliron/glsp/server"
)

var debug bool

func main() {
	commonlog.Configure(1, nil)

	var cmd = &cobra.Command{
		Use: lsp.Name,
		PersistentPreRun: func(_ *cobra.Command, _ []string) {
			if debug {
				commonlog.SetMaxLevel(nil, commonlog.Debug)
			}
		},
		Run: func(_ *cobra.Command, _ []string) {
			err := run()
			if err != nil {
				log.Error(err.Error())
				os.Exit(1)
			}
		},
		Version: lsp.Version,
	}

	cmd.Flags().BoolVarP(&debug, "debug", "d", false, "increase verbosity of log messages")

	if err := cmd.Execute(); err != nil {
		os.Exit(0)
	}
}

func run() error {
	server := serverpkg.NewServer(&lsp.Handler, lsp.Name, debug)
	return server.RunStdio()
}
