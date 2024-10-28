package config

import (
	"os"
	"path"
)

func init() {
	ReloadConfig()
}

var loadedFiles []File

func ReloadConfig() {
	loadedFiles = []File{}

	paths := []string{
		path.Join(".git", "config"),
		path.Join(os.Getenv("HOME"), ".gitconfig"),
		path.Join("/etc", "gitconfig"),
	}

	for _, p := range paths {
		file, err := LoadFile(p)
		if err == nil {
			loadedFiles = append(loadedFiles, file)
		}
	}
}

func Has(key string) bool {
	for _, file := range loadedFiles {
		if file.Has(key) {
			return true
		}
	}

	return false
}

func Get(key string) (string, bool) {
	for _, file := range loadedFiles {
		if value, ok := file.Get(key); ok {
			return value, true
		}
	}

	return "", false
}
