package tests

import "errors"

/////////////////
// READ CLOSER //
/////////////////

type MockReadCloser struct {
	ReadErr  error
	CloseErr error
}

/*
Creates a tests.MockReadCloser object with default read and close errors. These
errors can be modified by accessing the ReadErr and CloseErr values.
*/
func NewDefaultMockReadCloser() *MockReadCloser {
	return &MockReadCloser{
		ReadErr:  errors.New("MockReadCloser: Read error"),
		CloseErr: errors.New("MockReadCloser: Close error"),
	}
}

func (e *MockReadCloser) Read(p []byte) (n int, err error) {
	return 0, e.ReadErr
}

func (e *MockReadCloser) Close() error {
	return e.CloseErr
}
