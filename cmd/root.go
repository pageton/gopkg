package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gopkg",
	Short: "Gopkg is a dependency manager for Go modules",
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
