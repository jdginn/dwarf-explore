/*
Copyright © 2022 Justin Ginn <Justindavid.ginn@gmail.com>

*/
package cmd

import (
	"github.com/spf13/cobra"

	"github.com/jdginn/dwarf-explore/explorer/interactive"
)

// interactiveCmd represents the interactive command
var interactiveCmd = &cobra.Command{
	Use:   "interactive",
	Short: "Launch an interactive session",
	Long:  `Launch an interactive dwarf-explore session`,
	Run: func(cmd *cobra.Command, args []string) {
		file, err := rootCmd.Flags().GetString("debug-file")
		if err != nil {
			panic(err)
		}
		interactive.Start(file)
	},
}

func init() {
	rootCmd.AddCommand(interactiveCmd)
}
