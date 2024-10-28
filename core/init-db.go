package core

import (
	"encoding/hex"
	"fmt"
	"github.com/untanky/git-charged/config"
	"github.com/untanky/git-charged/plumbing"
	"io"
	"os"
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

	err = os.WriteFile(path.Join(gitDirectory, "HEAD"), []byte("ref: refs/heads/main"), 0644)

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
