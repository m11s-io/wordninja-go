# wordninja-go

Go port of [wordninja](https://github.com/keredson/wordninja), probabilistic text segmentation for splitting concatenated words.

It bundles the upstream English word list via `go:embed`, so default splitting works without external files.

## Install

```sh
go get github.com/m11s-io/wordninja-go
```

## Usage

```go
package main

import (
	"fmt"

	wordninja "github.com/m11s-io/wordninja-go"
)

func main() {
	fmt.Println(wordninja.Split("derekanderson"))
	fmt.Println(wordninja.Split("imateapot"))
	fmt.Println(wordninja.Split("thequickbrownfoxjumpsoverthelazydog"))
}
```

Load a custom gzip-compressed language model with one word per line in decreasing probability:

```go
lm, err := wordninja.NewLanguageModelFile("my_lang.txt.gz")
if err != nil {
	panic(err)
}
fmt.Println(lm.Split("derek"))
```

## Compatibility

This package ports the checked-in Python implementation, including its handling of whitespace and unknown punctuation. Whitespace separators are preserved in the output, repeated whitespace runs are returned as individual one-character tokens, and unknown punctuation can force following unknown spans to split into individual characters.

## License

MIT - see [LICENSE](./LICENSE).
