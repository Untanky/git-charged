package core

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/google/go-github/v66/github"
	"github.com/untanky/git-charged/config"
	"github.com/untanky/git-charged/plumbing"
	"io"
	"os"
	"os/exec"
	"path"
	"time"
)

const gitDirectoryName = ".git"

type InitDBParams struct {
	Name          string
	CreateLicense bool
	GitIgnoreFile *os.File
	ReadmeFile    *os.File
	LicenseFile   *os.File
}

func InitDB(params InitDBParams) error {
	gitDirectory := gitDirectoryName

	err := os.MkdirAll(gitDirectory, os.ModePerm)
	if err != nil {
		return fmt.Errorf("cannot create .git directory: %w", err)
	}

	err = createDirs(gitDirectory, "objects",
		"refs",
		path.Join("refs", "heads"),
		path.Join("refs", "tags"),
		"hooks",
		"info",
		"logs",
		"temp",
	)
	if err != nil {
		return fmt.Errorf("cannot create .git directory: %w", err)
	}

	plumbing.SetDirectory(gitDirectory)

	tree := plumbing.NewTree()

	if params.GitIgnoreFile != nil {
		hash, err := addFile(".gitignore", params.GitIgnoreFile)
		if err != nil {
			return fmt.Errorf("cannot create .gitignore: %w", err)
		}

		tree.AddObject(plumbing.ObjectTypeFile|0644, ".gitignore", hash)
	}

	if params.ReadmeFile != nil {
		hash, err := addFile("README.md", params.ReadmeFile)
		if err != nil {
			return fmt.Errorf("cannot create README.md: %w", err)
		}

		tree.AddObject(plumbing.ObjectTypeFile|0644, "README.md", hash)
	}

	if params.LicenseFile != nil {
		hash, err := addFile("LICENSE", params.LicenseFile)
		if err != nil {
			return fmt.Errorf("cannot create LICENSE: %w", err)
		}

		tree.AddObject(plumbing.ObjectTypeFile|0644, "LICENSE", hash)
	}

	hash, err := plumbing.WriteObject(tree)

	name, ok := config.Get("user.name")
	if !ok {
		return fmt.Errorf("no user name found")
	}

	email, ok := config.Get("user.email")
	if !ok {
		return fmt.Errorf("no user email found")
	}

	me := plumbing.AuthorData{
		Name:      name,
		Email:     email,
		Timestamp: time.Now(),
	}
	commit := plumbing.Commit{
		Tree:      hash,
		Author:    me,
		Committer: me,
		Message:   "Initial commit\n",
	}

	hash, err = plumbing.WriteObject(&commit)

	err = os.WriteFile(path.Join(gitDirectory, "refs", "heads", "main"), []byte(hex.EncodeToString(hash)), 0644)
	if err != nil {
		return err
	}

	client := github.NewClient(nil)
	client = client.WithAuthToken(os.Getenv("GITHUB_TOKEN"))

	repository := &github.Repository{
		Name: github.String("test"),
	}

	repository, _, err = client.Repositories.Create(context.TODO(), "", repository)
	if err != nil {
		return fmt.Errorf("cannot create repository: %w", err)
	}

	err = setGitConfig(repository.GetSSHURL())
	if err != nil {
		return fmt.Errorf("cannot set git config: %w", err)
	}

	err = os.WriteFile(path.Join(gitDirectory, "HEAD"), []byte("ref: refs/heads/main"), 0644)
	if err != nil {
		return fmt.Errorf("cannot create HEAD: %w", err)
	}

	err = exec.Command("git", "push").Run()
	if err != nil {
		return fmt.Errorf("cannot push HEAD: %w", err)
	}

	return nil
}

func createDirs(gitDirectory string, directories ...string) error {
	for _, directory := range directories {
		err := os.Mkdir(path.Join(gitDirectory, directory), os.ModePerm)
		if err != nil {
			return fmt.Errorf("cannot create .git directory: %w", err)
		}
	}

	return nil
}

func addFile(filepath string, file *os.File) ([]byte, error) {
	stat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("cannot create .git directory: %w", err)
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, fmt.Errorf("cannot seek to %s: %w", filepath, err)
	}

	blob := plumbing.NewBlob(uint32(stat.Size()), file)

	hash, err := plumbing.WriteObject(blob)
	if err != nil {
		return nil, fmt.Errorf("cannot write %s: %w", filepath, err)
	}

	return hash, nil
}

func setGitConfig(remoteUrl string) error {
	file, err := os.OpenFile(".git/config", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	_, err = file.WriteString(fmt.Sprintf(`[core]
    repositoryformatversion = 0
    filemode = true
    bare = false
    ignorecase = true
    precomposeunicode = true
[remote "origin"]
    url = %s
    fetch = +refs/heads/*:refs/remotes/origin/*
[branch "main"]
    remote = origin
    merge = refs/heads/main
`, remoteUrl))
	if err != nil {
		return err
	}

	return nil
}
