package wordninja

import (
	"bytes"
	"compress/gzip"
	_ "embed"
	"io"
)

//go:embed data/wordninja_words.txt.gz
var defaultWordsGz []byte

// DefaultWords returns a reader over the bundled English word list.
func DefaultWords() (io.Reader, error) {
	return gzip.NewReader(bytes.NewReader(defaultWordsGz))
}
