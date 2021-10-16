package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
  Use:   "version",
  Short: "Print the version number of liquigo",
  Long:  `All software has versions. This is liquigo's`,
  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println("liquigo DB migration v0.1.0")
  },
}
