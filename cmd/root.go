package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd is the entry
var (
	rootCmd = &cobra.Command{
		Use:   "token-price-sp",
		Short: "service-providers-go-go",
	}
)

func main() {
	cobra.EnableCommandSorting = false

	rootCmd.AddCommand(startCmd())
	rootCmd.AddCommand(keysCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
