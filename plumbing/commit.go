package plumbing

import (
	"fmt"
	"io"
	"time"
)

type AuthorData struct {
	Email     string
	Name      string
	Timestamp time.Time
}

type Commit struct {
	Tree      []byte
	Parent    []byte
	Author    AuthorData
	Committer AuthorData
	Message   string
}

func (c *Commit) WriteTo(w io.Writer) (int64, error) {
	var data string
	if c.Parent == nil {
		data = fmt.Sprintf(`tree %x
author %s <%s> %d +0200
committer %s <%s> %d +0200

%s`, c.Tree, c.Author.Name, c.Author.Email, c.Author.Timestamp.Unix(), c.Committer.Name, c.Committer.Email, c.Committer.Timestamp.Unix(), c.Message)
	} else {
		data = fmt.Sprintf(`tree %x
parent %x
author %s <%s> %d +0200
committer %s %s %d +0200

%s`, c.Tree, c.Parent, c.Author.Name, c.Author.Email, c.Author.Timestamp.Unix(), c.Committer.Name, c.Committer.Email, c.Committer.Timestamp.Unix(), c.Message)
	}

	m, err := fmt.Fprintf(w, "commit %d\000%s", len(data), data)
	if err != nil {
		return int64(m), err
	}

	return int64(m), err
}
