/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"context"
	"fmt"
	"github.com/google/go-github/v66/github"
	"github.com/spf13/cobra"
	"github.com/untanky/git-charged/core"
	"github.com/untanky/git-charged/ui"
	"log"
	"os"
	"os/exec"
	"path"
)

var client *github.Client

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

		noGitignore, err := cmd.Flags().GetBool("no-gitignore")
		if err != nil {
			noGitignore = false
		}

		var gitIgnoreReader *bytes.Reader
		if !noGitignore {
			gitIgnoreReader, err = selectGitignore()
			if err != nil {
				log.Fatalf("failed to init git: %s", err)
			}
		}

		noReadme, err := cmd.Flags().GetBool("no-readme")

		var readmeReader *bytes.Reader
		if !noReadme {
			readmeReader, err = selectReadme(directory)
			if err != nil {
				log.Fatalf("failed to init git: %s", err)
			}
		}

		err = core.InitDB(core.InitDBParams{
			Directory:       directory,
			CreateGitignore: !noGitignore,
			GitIgnoreReader: gitIgnoreReader,
			CreateReadme:    true,
			ReadmeReader:    readmeReader,
		})
		if err != nil {
			log.Fatalf("failed to init git: %s", err)
		}
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
	initCmd.Flags().Bool("no-readme", false, "Do not generate a README file")
}

func selectGitignore() (*bytes.Reader, error) {
	gitignoreOption, _, err := client.Gitignores.List(context.TODO())
	if err != nil {
		return nil, err
	}

	selectedOptions, err := ui.NewMultiSelect(gitignoreOption).Run()
	if err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(make([]byte, 0))

	for _, option := range selectedOptions {
		gitignore, _, err := client.Gitignores.Get(context.TODO(), option)
		if err != nil {
			return nil, err
		}
		buffer.WriteString(fmt.Sprintf("# Gitignore for %s\n\n", option))
		buffer.WriteString(gitignore.GetSource())
		buffer.WriteString("\n\n")
	}

	return bytes.NewReader(buffer.Bytes()), nil
}

func selectReadme(directory string) (*bytes.Reader, error) {
	file, err := os.Create(path.Join(directory, "README.md"))
	if err != nil {
		return nil, err
	}
	_, err = fmt.Fprintf(file, "# %s\n\n[//]: # (Write something about your new project)", path.Base(directory))
	if err != nil {
		return nil, err
	}
	err = file.Close()
	if err != nil {
		return nil, err
	}

	cmd := exec.Command("vim", path.Join(directory, "README.md"))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		return nil, err
	}

	return nil, nil
}
