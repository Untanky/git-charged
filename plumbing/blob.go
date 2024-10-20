package plumbing

import (
	"fmt"
	"io"
)

type Blob interface {
	Object
}

type blob struct {
	size   uint32
	reader io.Reader
}

func NewBlob(size uint32, reader io.Reader) Blob {
	return &blob{size: size, reader: reader}
}

func (b *blob) WriteTo(w io.Writer) (n int64, err error) {
	var m int
	m, err = fmt.Fprintf(w, "blob %d\000", b.size)
	n += int64(m)
	if err != nil {
		return n, err
	}

	var k int64
	k, err = io.Copy(w, b.reader)
	n += k
	if err != nil {
		return n, err
	}

	return n, nil
}
