package console

import (
	"bufio"
	"os"
)

// LineReader could read buffered line
type LineReader struct {
	file   *os.File
	reader *bufio.Reader
	lf     byte
}

// NewLineReader construct a reader object
func NewLineReader(path string, delim byte) (*LineReader, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	buffered := bufio.NewReader(f)

	return &LineReader{
		file:   f,
		reader: buffered,
		lf:     delim,
	}, nil
}

// Close close a file of LineReader
func (r *LineReader) Close() {
	r.file.Close()
}

// Line read buffered line
func (r *LineReader) Line() (string, error) {
	return r.reader.ReadString(r.lf)
}
