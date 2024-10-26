/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"github.com/google/go-github/v66/github"
	"github.com/spf13/cobra"
	"github.com/untanky/git-charged/ui"
	"log"
)

var client *github.Client

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a new git project ⚡️super-charged⚡️",
	Long:  ``, // TODO: add long description
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		//var directory string
		//if len(args) == 1 {
		//	directory = args[0]
		//} else {
		//	var err error
		//	directory, err = os.Getwd()
		//	if err != nil {
		//		log.Fatalf("failed to get current directory: %s", err)
		//	}
		//}

		noGitignore, err := cmd.Flags().GetBool("no-gitignore")
		if err != nil {
			noGitignore = false
		}

		if !noGitignore {
			err = selectGitignore()
			if err != nil {
				log.Fatalf("failed to init git: %s", err)
			}
		}

		//err = core.InitDB(core.InitDBParams{
		//	Directory:       directory,
		//	CreateGitignore: !noGitignore,
		//	CreateReadme:    true,
		//})
		//if err != nil {
		//	log.Fatalf("failed to init git: %s", err)
		//}
	},
}

func init() {
	client = github.NewClient(nil)

	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	initCmd.Flags().Bool("no-gitignore", false, "Do not generate a .gitignore file")
}

func selectGitignore() error {
	gitignoreOption, _, err := client.Gitignores.List(context.TODO())
	if err != nil {
		return err
	}

	fmt.Println(ui.NewMultiSelect(gitignoreOption).Run())
	return nil
}
