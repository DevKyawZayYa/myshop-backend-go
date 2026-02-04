package cmd

import (
	"errors"
	"fmt"
	"os"
	"product-service/cmd/server"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "server",
	Short:        "Setting Server",
	SilenceUsage: true,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires at least one arg")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to Product Service")
	},
}

func init() {
	rootCmd.AddCommand(server.StartLocalCmd)
	rootCmd.AddCommand(server.StartDevelopmentCmd)
	rootCmd.AddCommand(server.StartProductionCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
