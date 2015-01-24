package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type FStore struct {
	Path    string
	TmpPath string
}

func (store *FStore) AddFile(r io.Reader) (string, error) {
	file, err := ioutil.TempFile(store.TmpPath, "rabbity-tmp")
	log.Printf("uploading to temporary file %v", file.Name())
	if err != nil {
		return "", err
	}
	written, hashsum, err := Sha3HashCopy(file, r)
	if err != nil {
		return "", err
	}
	if err := file.Close(); err != nil {
		return "", err
	}
	log.Printf("Wrote %v bytes to %v (hashsum is %v)", written, file.Name(), hashsum)
	log.Printf("Moving %v to %v", file.Name(), store.getHashPath(hashsum))
	err = os.Rename(file.Name(), store.getHashPath(hashsum))
	return hashsum, nil
}

func (store *FStore) getHashPath(hashsum string) string {
	dirPath := filepath.Join(store.Path, string(hashsum[0]), string(hashsum[1]))
	err := os.MkdirAll(dirPath, 0700)
	if err != nil {
		panic(err)
	}
	return filepath.Join(dirPath, hashsum[2:])
}

func (store *FStore) getFilePath(hashsum string) string {
	return filepath.Join(store.Path, string(hashsum[0]),
		string(hashsum[1]), hashsum[2:])
}

func (store *FStore) HasFile(hashsum string) bool {
	filePath := store.getFilePath(hashsum)
	if _, err := os.Stat(filePath); err == nil {
		return true
	}
	return false
}

func (store *FStore) GetFile(hashsum string) (io.Reader, error) {
	return os.Open(store.getFilePath(hashsum))
}
