/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/untanky/git-charged/core"
	"log"
	"os"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a new git project ⚡️super-charged⚡️",
	Long:  ``, // TODO: add long description
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var directory string
		if len(args) == 1 {
			directory = args[0]
		} else {
			var err error
			directory, err = os.Getwd()
			if err != nil {
				log.Fatalf("failed to get current directory: %s", err)
			}
		}

		err := core.InitDB(core.InitDBParams{
			Directory:       directory,
			CreateGitignore: true,
			CreateReadme:    true,
		})
		if err != nil {
			log.Fatalf("failed to init git: %s", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
