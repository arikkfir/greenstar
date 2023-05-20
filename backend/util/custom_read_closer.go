package util

import "io"

type ReaderCloser struct {
	Reader io.Reader
}

func (rc *ReaderCloser) Read(p []byte) (int, error) {
	return rc.Reader.Read(p)
}

func (rc *ReaderCloser) Close() error {
	if closer, ok := rc.Reader.(io.Closer); ok {
		return closer.Close()
	} else {
		return nil
	}
}
