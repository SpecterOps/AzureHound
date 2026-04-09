package bloodhound

import (
	"bytes"
	"io"
)

type rewindableByteReader struct {
	data *bytes.Reader
}

func (s *rewindableByteReader) Read(p []byte) (int, error) {
	return s.data.Read(p)
}

func (s *rewindableByteReader) Close() error {
	return nil
}

func (s *rewindableByteReader) Rewind() (int64, error) {
	return s.data.Seek(0, io.SeekStart)
}
