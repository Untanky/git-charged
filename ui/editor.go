package ui

import (
	"github.com/untanky/git-charged/config"
	"os"
	"os/exec"
	"strings"
)

func OpenEditor(filepath string) error {
	editor, ok := config.Get("core.editor")
	if !ok {
		editor = "vim"
	}

	split := strings.Split(editor, " ")
	cmd := exec.Command(split[0], append(split[1:], filepath)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
