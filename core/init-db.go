package core

import (
	"bytes"
	"fmt"
	"github.com/untanky/git-charged/plumbing"
	"io"
	"os"
	"path"
)

const gitDirectoryName = ".git"

type InitDBParams struct {
	Name            string
	Directory       string
	CreateLicense   bool
	CreateReadme    bool
	CreateGitignore bool
}

func InitDB(params InitDBParams) error {
	gitDirectory := path.Join(params.Directory, gitDirectoryName)

	err := os.MkdirAll(gitDirectory, os.ModePerm)
	if err != nil {
		return fmt.Errorf("cannot create .git directory: %w", err)
	}

	err = os.Mkdir(path.Join(gitDirectory, "objects"), os.ModePerm)
	if err != nil {
		return fmt.Errorf("cannot create .git directory: %w", err)
	}

	err = os.Mkdir(path.Join(gitDirectory, "temp"), os.ModePerm)
	if err != nil {
		return fmt.Errorf("cannot create .git directory: %w", err)
	}

	plumbing.SetDirectory(params.Directory)

	tree := plumbing.NewTree()

	if params.CreateGitignore {
		filepath := path.Join(params.Directory, ".gitignore")
		hash, err := createFile(gitDirectory, filepath, 5, bytes.NewReader([]byte("test\n")))
		if err != nil {
			return fmt.Errorf("cannot create .gitignore: %w", err)
		}

		tree.AddObject(plumbing.ObjectTypeFile|0644, ".gitignore", hash)
	}

	if params.CreateReadme {
		filepath := path.Join(params.Directory, "README.md")
		hash, err := createFile(gitDirectory, filepath, 13, bytes.NewReader([]byte("# Hello World\n")))
		if err != nil {
			return fmt.Errorf("cannot create .gitignore: %w", err)
		}

		tree.AddObject(plumbing.ObjectTypeFile|0644, "README.md", hash)
	}
	_, err = plumbing.WriteObject(tree)

	return nil
}

func createFile(gitDirectory string, filepath string, size uint32, reader io.Reader) ([]byte, error) {
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
