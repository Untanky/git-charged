package core

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/untanky/git-charged/plumbing"
	"io"
	"os"
	"path"
	"time"
)

const gitDirectoryName = ".git"

type InitDBParams struct {
	Name            string
	Directory       string
	CreateLicense   bool
	CreateReadme    bool
	CreateGitignore bool
	GitIgnoreReader *bytes.Reader
	ReadmeReader    *bytes.Reader
}

func InitDB(params InitDBParams) error {
	gitDirectory := path.Join(params.Directory, gitDirectoryName)

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

	plumbing.SetDirectory(params.Directory)

	tree := plumbing.NewTree()

	if params.CreateGitignore {
		filepath := path.Join(params.Directory, ".gitignore")
		hash, err := createFile(filepath, uint32(params.GitIgnoreReader.Size()), params.GitIgnoreReader)
		if err != nil {
			return fmt.Errorf("cannot create .gitignore: %w", err)
		}

		tree.AddObject(plumbing.ObjectTypeFile|0644, ".gitignore", hash)
	}

	if params.CreateReadme {
		filepath := path.Join(params.Directory, "README.md")
		hash, err := createFile(filepath, 13, bytes.NewReader([]byte("# Hello World\n")))
		if err != nil {
			return fmt.Errorf("cannot create .gitignore: %w", err)
		}

		tree.AddObject(plumbing.ObjectTypeFile|0644, "README.md", hash)
	}

	hash, err := plumbing.WriteObject(tree)

	me := plumbing.AuthorData{
		Name:      "Lukas Grimm",
		Email:     "lukaskingsmail@gmail.com",
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

func createFile(filepath string, size uint32, reader io.Reader) ([]byte, error) {
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("cannot create %s: %w", filepath, err)
	}
	defer file.Close()

	_, err = io.Copy(file, reader)
	if err != nil {
		return nil, fmt.Errorf("cannot write %s: %w", filepath, err)
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, fmt.Errorf("cannot seek to %s: %w", filepath, err)
	}

	blob := plumbing.NewBlob(size, file)

	hash, err := plumbing.WriteObject(blob)
	if err != nil {
		return nil, fmt.Errorf("cannot write %s: %w", filepath, err)
	}

	return hash, nil
}
