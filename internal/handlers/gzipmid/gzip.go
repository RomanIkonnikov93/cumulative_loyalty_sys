package gzipmid

import (
	"bytes"
	"compress/gzip"
	"io"
	"log"
	"net/http"
)

type GzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w GzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func DecompressGZIP(data []byte) ([]byte, error) {

	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		log.Printf("%v", err)
	}
	defer r.Close()

	b, err := io.ReadAll(r)
	if err != nil {
		log.Printf("%v", err)
	}

	return b, nil
}
