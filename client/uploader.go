package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"nurio.at/testi"
	"os"
	"path/filepath"
)

// Creates a new file upload http request with optional extra params
func newfileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, chan error, chan string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, nil, err
	}
	bodyReader, bodyWriter := io.Pipe()
	multiWriter := multipart.NewWriter(bodyWriter)
	errChan := make(chan error, 1)
	hashChan := make(chan string, 1)
	go func() {
		defer bodyWriter.Close()
		defer file.Close()
		part, err := multiWriter.CreateFormFile(paramName, filepath.Base(path))
		if err != nil {
			errChan <- err
			return
		}
		_, sha3sum, err := testi.Sha3HashCopy(part, file)
		if err != nil {
			errChan <- err
			return
		}
		for k, v := range params {
			if err := multiWriter.WriteField(k, v); err != nil {
				errChan <- err
				return
			}
		}
		errChan <- multiWriter.Close()
		hashChan <- sha3sum
	}()
	req, err := http.NewRequest("POST", uri, bodyReader)
	req.Header.Set("Content-Type", "multipart/form-data; boundary="+multiWriter.Boundary())
	return req, errChan, hashChan, err
}

func main() {
	path := os.Args[1]
	extraParams := map[string]string{
		"title":       "My Document",
		"author":      "Matt Aimonetti",
		"description": "A document with all the Go programming language secrets",
	}
	request, _, hashChan, err := newfileUploadRequest("http://localhost:8080/upload", extraParams, "payload", path)
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	} else {
		body := &bytes.Buffer{}
		_, err := body.ReadFrom(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		resp.Body.Close()
		fmt.Println(resp.StatusCode)
		fmt.Println(resp.Header)
		fmt.Println(body)
		hash := <-hashChan
		fmt.Println("calculated hash: " + hash)
	}
}
