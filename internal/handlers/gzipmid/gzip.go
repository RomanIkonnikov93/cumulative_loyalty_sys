package gzipmid

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"

	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/logging"
)

type GzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w GzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func DecompressGZIP(data []byte) ([]byte, error) {

	logger := logging.GetLogger()

	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		logger.Error(err)
	}
	defer r.Close()

	b, err := io.ReadAll(r)
	if err != nil {
		logger.Error(err)
	}

	return b, nil
}
