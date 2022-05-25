/*
Copyright © 2022 Justin Ginn <Justindavid.ginn@gmail.com>

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dwarf-explore",
	Short: "Interactive utility for exploring types and variables exposed via DWARF",
	Long: `dwarf-explore is an interractive utility for exploring types and
variables exposed by DWARF in any binary compiled with debug symbols.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	TraverseChildren: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dwarf-explore.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().StringP("debug-file", "d", "", "Path to DWARF debug file")
	rootCmd.MarkFlagRequired("debug-file")
}
