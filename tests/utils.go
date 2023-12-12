package tests

import "errors"

// Helper to handle bad read errors
type ErrorReader int

func (ErrorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("ErrorReader error")
}
