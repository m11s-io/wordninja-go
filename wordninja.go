package wordninja

import (
	"bufio"
	"compress/gzip"
	"errors"
	"io"
	"math"
	"os"
	"regexp"
	"strings"
	"unicode"
)

var splitRE = regexp.MustCompile(`\s+`)

// LanguageModel splits concatenated words using word frequencies ranked by
// decreasing probability.
type LanguageModel struct {
	wordCost map[string]float64
	maxWord  int
}

// NewLanguageModel builds a language model from a newline- or whitespace-
// separated word list ordered from most likely to least likely.
func NewLanguageModel(r io.Reader) (*LanguageModel, error) {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)
	scanner.Buffer(make([]byte, 1024), 1024*1024)

	var words []string
	maxWord := 0
	for scanner.Scan() {
		word := scanner.Text()
		words = append(words, word)
		if l := len([]rune(word)); l > maxWord {
			maxWord = l
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if len(words) == 0 {
		return nil, errors.New("wordninja: language model has no words")
	}

	wordCost := make(map[string]float64, len(words))
	for i, word := range words {
		wordCost[word] = math.Log(float64(i+1) * math.Log(float64(len(words))))
	}

	return &LanguageModel{wordCost: wordCost, maxWord: maxWord}, nil
}

// NewLanguageModelFile builds a language model from a gzip-compressed word
// list file ordered from most likely to least likely.
func NewLanguageModelFile(path string) (*LanguageModel, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return nil, err
	}
	defer gz.Close()

	return NewLanguageModel(gz)
}

func must(model *LanguageModel, err error) *LanguageModel {
	if err != nil {
		panic(err)
	}
	return model
}

// DefaultLanguageModel is the bundled English model.
var DefaultLanguageModel = must(NewLanguageModel(mustReader(DefaultWords())))

func mustReader(r io.Reader, err error) io.Reader {
	if err != nil {
		panic(err)
	}
	return r
}

// Split uses the bundled English model to infer spaces in s.
func Split(s string) []string {
	return DefaultLanguageModel.Split(s)
}

// Split infers spaces in s using dynamic programming.
func (lm *LanguageModel) Split(s string) []string {
	parts := make([]string, 0)
	last := 0
	for _, loc := range splitRE.FindAllStringIndex(s, -1) {
		parts = append(parts, lm.splitText(s[last:loc[0]])...)
		for _, r := range s[loc[0]:loc[1]] {
			parts = append(parts, string(r))
		}
		last = loc[1]
	}
	parts = append(parts, lm.splitText(s[last:])...)
	return parts
}

func (lm *LanguageModel) splitText(s string) []string {
	runes := []rune(s)
	cost := make([]float64, len(runes)+1)

	for i := 1; i <= len(runes); i++ {
		c, _ := lm.bestMatch(runes, cost, i)
		cost[i] = c
	}

	out := make([]string, 0)
	for i := len(runes); i > 0; {
		_, k := lm.bestMatch(runes, cost, i)
		token := string(runes[i-k : i])

		newToken := true
		if token != "'" && len(out) > 0 {
			if out[len(out)-1] == "'s" || (unicode.IsDigit(runes[i-1]) && startsWithDigit(out[len(out)-1])) {
				out[len(out)-1] = token + out[len(out)-1]
				newToken = false
			}
		}
		if newToken {
			out = append(out, token)
		}

		i -= k
	}

	for i, j := 0, len(out)-1; i < j; i, j = i+1, j-1 {
		out[i], out[j] = out[j], out[i]
	}
	return out
}

func (lm *LanguageModel) bestMatch(runes []rune, cost []float64, i int) (float64, int) {
	limit := lm.maxWord
	if i < limit {
		limit = i
	}

	bestCost := math.Inf(1)
	bestLength := 1
	for k := 1; k <= limit; k++ {
		word := strings.ToLower(string(runes[i-k : i]))
		candidate := cost[i-k] + lm.wordCostWithDefault(word)
		if candidate < bestCost || (candidate == bestCost && k < bestLength) {
			bestCost = candidate
			bestLength = k
		}
	}
	return bestCost, bestLength
}

func (lm *LanguageModel) wordCostWithDefault(word string) float64 {
	if cost, ok := lm.wordCost[word]; ok {
		return cost
	}
	return math.Inf(1)
}

func startsWithDigit(s string) bool {
	for _, r := range s {
		return unicode.IsDigit(r)
	}
	return false
}
