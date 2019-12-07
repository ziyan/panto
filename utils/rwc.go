package utils

import (
	"io"
)

func NewReadWriteCloser(rc io.ReadCloser, wc io.WriteCloser) io.ReadWriteCloser {
	return &readWriteCloser{ReadCloser: rc, WriteCloser: wc}
}

type readWriteCloser struct {
	io.ReadCloser
	io.WriteCloser
}

func (rwc *readWriteCloser) Close() error {
	err := rwc.WriteCloser.Close()
	err2 := rwc.ReadCloser.Close()
	if err == nil {
		err = err2
	}
	return err
}
