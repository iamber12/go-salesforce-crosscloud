package main

import (
	"github.com/spf13/cobra"
	"go-salesforce-crosscloud/cmd/serve"
	"log"
)

func main() {
	rootCmd := &cobra.Command{
		Use:  "webrtc-and-go",
		Long: "webrtc-and-go",
	}

	serveCmd := serve.NewServeCommand()
	rootCmd.AddCommand(serveCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("error running command: %v", err)
	}
}
