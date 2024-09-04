package cmd

import (
	"fmt"
	"html-converter/converters"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "htmlConverter <filename>",
	Short: "Markdown to HTML converter",
	Long:  `A command-line tool to convert Markdown format to HTML.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := converters.ProcessFile(args[0]); err != nil {
			fmt.Printf("Error processing file: %v\n", err)
			os.Exit(1)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
