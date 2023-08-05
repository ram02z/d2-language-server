package main

import (
	"fmt"
	"os"

	"github.com/ram02z/d2-language-server/internal/lsp"
	"github.com/spf13/cobra"
)

const lsName = "d2-language-server"

var lsVersion string = "0.0.1"

func main() {
	var cmd = &cobra.Command{
		Use: lsName,
		Run: func(_ *cobra.Command, _ []string) {
			err := run()
			if err != nil {
				fmt.Fprint(os.Stderr, err.Error())
				os.Exit(1)
			}
		},
	}

	cmd.Version = lsVersion

	if err := cmd.Execute(); err != nil {
		os.Exit(0)
	}
}

func run() error {
	server := lsp.NewServer(lsp.ServerOpts{
		Name:    lsName,
		Version: lsVersion,
	})

	return server.Run()
}
