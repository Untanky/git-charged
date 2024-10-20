package plumbing

import (
	"bytes"
	"fmt"
	"io"
)

const (
	ObjectTypeFile         = 0b1000_0000_0000_0000
	ObjectTypeDirectory    = 0b0100_0000_0000_0000
	ObjectTypeSymbolicLink = 0b1010_0000_0000_0000
	ObjectTypeGitLink      = 0b1110_0000_0000_0000
)

type Tree interface {
	Object
	AddObject(mode uint16, name string, hash []byte)
}

type treeEntry struct {
	mode uint16
	name string
	hash []byte
}

func (entry *treeEntry) WriteTo(writer io.Writer) (n int64, err error) {
	m, err := fmt.Fprintf(writer, "%o %s\000", entry.mode, entry.name)
	n += int64(m)
	if err != nil {
		return n, err
	}

	m, err = writer.Write(entry.hash)
	n += int64(m)
	if err != nil {
		return n, err
	}

	return n, nil
}

type tree struct {
	entries []treeEntry
}

func NewTree() Tree {
	return &tree{
		entries: make([]treeEntry, 0),
	}
}

func (t *tree) AddObject(mode uint16, name string, hash []byte) {
	t.entries = append(t.entries, treeEntry{mode, name, hash})
}

func (t *tree) WriteTo(writer io.Writer) (n int64, err error) {
	buffer := bytes.NewBuffer(make([]byte, 0, 1024))

	for _, entry := range t.entries {
		_, err = entry.WriteTo(buffer)
		if err != nil {
			return 0, err
		}
	}

	header := fmt.Sprintf("tree %d\000", len(buffer.Bytes()))
	m, err := writer.Write([]byte(header))
	n += int64(m)
	if err != nil {
		return
	}

	k, err := io.Copy(writer, buffer)
	n += k
	if err != nil {
		return
	}

	return n, nil
}
