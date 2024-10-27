package plumbing

import (
	"compress/zlib"
	"crypto"
	_ "crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path"
)

type Object interface {
	io.WriterTo
}

const (
	objectsDirectory = "objects"
)

var (
	hashFactory = crypto.SHA1

	gitDirectory = ".git"
)

func SetDirectory(directory string) {
	gitDirectory = directory
}

func WriteObject(object Object) ([]byte, error) {
	hashWriter := hashFactory.New()
	temporaryFilename := path.Join(gitDirectory, "temp", "0")
	file, err := os.Create(temporaryFilename)
	if err != nil {
		return nil, err
	}
	zlibWriter := zlib.NewWriter(file)
	writer := io.MultiWriter(hashWriter, zlibWriter)

	_, err = object.WriteTo(writer)
	if err != nil {
		zlibWriter.Close()
		file.Close()
		return nil, err
	}

	zlibWriter.Close()
	file.Close()

	hash := hashWriter.Sum(nil)
	hexa := hex.EncodeToString(hash)

	fmt.Println(hexa)

	err = os.Mkdir(path.Join(gitDirectory, objectsDirectory, hexa[:2]), 0755)
	if err != nil && !os.IsExist(err) {
		return nil, err
	}

	filename := path.Join(gitDirectory, objectsDirectory, hexa[:2], hexa[2:])

	err = os.Rename(temporaryFilename, filename)
	if err != nil {
		return nil, err
	}

	return hash, nil
}
