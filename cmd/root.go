package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "liquigo",
	Short: "liquigo is a flexible DB migration tool",
	Long: `liquigo is a flexible DB migration tool
                That doesn't require you to have sequential DB versions.
                Deliver your new DB version when it's ready, and fasttack another meanwhile!`,
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(initUpdateCmd())
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
