package http

import (
	"io"

	"github.com/goccy/go-json"
)

func ReadBody(body io.ReadCloser, receiver any) error {
	err := json.NewDecoder(body).Decode(receiver)
	if err != nil {
		return err
	}

	return nil
}
