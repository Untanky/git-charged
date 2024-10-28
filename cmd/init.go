/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/go-github/v66/github"
	"github.com/spf13/cobra"
	"github.com/untanky/git-charged/core"
	"github.com/untanky/git-charged/ui"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

var client *github.Client

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a new git project ⚡️super-charged⚡️",
	Long:  ``, // TODO: add long description
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			err := os.Mkdir(args[0], os.ModePerm)
			if err != nil && !errors.Is(err, os.ErrNotExist) {
				log.Fatalf("failed to create directory: %s", err)
			}

			err = os.Chdir(args[0])
			if err != nil {
				log.Fatalf("failed to change directory: %s", err)
			}
		}

		noGitignore, err := cmd.Flags().GetBool("no-gitignore")
		if err != nil {
			noGitignore = false
		}

		var gitignoreFile *os.File
		if !noGitignore {
			gitignoreFile, err = os.OpenFile(".gitignore", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
			if err != nil {
				log.Fatalf("failed to create .gitignore: %s", err)
			}

			err = selectGitignore(gitignoreFile)
			if err != nil {
				log.Fatalf("failed to init git: %s", err)
			}
		}

		noReadme, err := cmd.Flags().GetBool("no-readme")

		var readmeReader *os.File
		if !noReadme {
			readmeReader, err = os.OpenFile("README.md", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
			err = selectReadme(readmeReader, "foo") // TODO: set correctly
			if err != nil {
				log.Fatalf("failed to init git: %s", err)
			}
		}

		noLicense, err := cmd.Flags().GetBool("no-license")

		var licenseReader *os.File
		if !noLicense {
			licenseReader, err = os.OpenFile("LICENSE", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
			err = selectLicense(licenseReader)
			if err != nil {
				log.Fatalf("failed to init git: %s", err)
			}
		}

		err = core.InitDB(core.InitDBParams{
			GitIgnoreFile: gitignoreFile,
			ReadmeFile:    readmeReader,
			LicenseFile:   licenseReader,
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
	initCmd.Flags().Bool("no-license", false, "Do not generate a LICENSE file")
}

func selectGitignore(writer io.Writer) error {
	gitignoreOption, _, err := client.Gitignores.List(context.TODO())
	if err != nil {
		return err
	}

	selectedOptions, err := ui.NewMultiSelect(gitignoreOption).Run()
	if err != nil {
		return err
	}

	for _, option := range selectedOptions {
		gitignore, _, err := client.Gitignores.Get(context.TODO(), option)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(writer, "# Gitignore for %s\n\n%s\n\n", option, gitignore.GetSource())
		if err != nil {
			return err
		}
	}

	return nil
}

func selectReadme(readmeFile *os.File, projectName string) error {
	_, err := fmt.Fprintf(readmeFile, "# %s\n\n[//]: # (Write something about your new project)", projectName)
	if err != nil {
		return err
	}

	err = ui.OpenEditor(readmeFile.Name())
	if err != nil {
		return err
	}

	return nil
}

func selectLicense(readmeFile *os.File) error {
	const owner = "github"
	const repo = "choosealicense.com"
	const licenseDirectory = "_licenses"
	_, directoryContent, _, err := client.Repositories.GetContents(context.TODO(), owner, repo, licenseDirectory, nil)
	if err != nil {
		return err
	}
	if directoryContent == nil {
		return errors.New("no licenses found")
	}

	options := make([]string, 0)

	for _, content := range directoryContent {
		name := content.GetName()
		options = append(options, strings.TrimSuffix(name, ".txt"))
	}

	selectedOption, err := ui.NewSelect(options).Run()
	if err != nil {
		return err
	}

	fileContent, _, _, err := client.Repositories.GetContents(context.TODO(), owner, repo, fmt.Sprintf("%s/%s.txt", licenseDirectory, selectedOption), nil)
	if err != nil {
		return err
	}
	if fileContent == nil {
		return errors.New("no licenses found")
	}

	content, err := fileContent.GetContent()
	if err != nil {
		return err
	}

	_, err = readmeFile.WriteString(removeFrontmatter(content))
	if err != nil {
		return err
	}

	err = ui.OpenEditor(readmeFile.Name())
	if err != nil {
		return err
	}

	return nil
}

func removeFrontmatter(content string) string {
	// Regular expression to match frontmatter
	re := regexp.MustCompile(`(?s)^---\s*.*?\n---\s*\n`)

	// Remove frontmatter
	contentWithoutFrontmatter := re.ReplaceAllString(content, "")

	// Trim any leading whitespace
	return strings.TrimLeft(contentWithoutFrontmatter, "\n\r\t ")
}
