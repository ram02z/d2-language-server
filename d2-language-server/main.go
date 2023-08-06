package main

import (
	"log"
	"os"

	"github.com/ram02z/d2-language-server/internal/lsp"
	"github.com/spf13/cobra"
	serverpkg "github.com/tliron/glsp/server"
)

var debug bool

func main() {

	var cmd = &cobra.Command{
		Use: lsp.Name,
		PersistentPreRun: func(_ *cobra.Command, _ []string) {
			if debug {
				// TODO: use logger
			}
		},
		Run: func(_ *cobra.Command, _ []string) {
			err := run()
			if err != nil {
				log.Fatal(err.Error())
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
