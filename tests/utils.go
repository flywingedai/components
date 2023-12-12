package tests

import "errors"

// Helper to handle bad read errors
type ErrorIO struct {
	ReadErr  error
	CloseErr error
}

func NewDefaultErrorIO() *ErrorIO {
	return &ErrorIO{
		ReadErr:  errors.New("ErrorIO: Read error"),
		CloseErr: errors.New("ErrorIO: Close error"),
	}
}

func (e *ErrorIO) Read(p []byte) (n int, err error) {
	return 0, e.ReadErr
}

func (e *ErrorIO) Close() error {
	return e.CloseErr
}
