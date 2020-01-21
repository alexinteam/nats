package cmd

import (
	"github.com/spf13/cobra"
	"log"
)

var rootCmd = &cobra.Command{
	Use:   "bus",
	Short: "Bus cmd",
	Long:  "Bus cmd",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err.Error())
	}
}

func init() {
	cobra.OnInitialize()
}
