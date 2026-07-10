package wordninja

import (
	"reflect"
	"testing"
)

func TestSplit(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want []string
	}{
		{name: "simple", in: "derekanderson", want: []string{"derek", "anderson"}},
		{name: "caps", in: "DEREKANDERSON", want: []string{"DEREK", "ANDERSON"}},
		{name: "digits", in: "win32intel", want: []string{"win", "32", "intel"}},
		{name: "apostrophes", in: "that'sthesheriff'sbadge", want: []string{"that's", "the", "sheriff's", "badge"}},
		{name: "readme teapot", in: "imateapot", want: []string{"im", "a", "teapot"}},
		{name: "readme phrase", in: "heshotwhointhewhatnow", want: []string{"he", "shot", "who", "in", "the", "what", "now"}},
		{name: "whitespace preserved", in: "derek anderson", want: []string{"derek", " ", "anderson"}},
		{name: "whitespace run exploded", in: "a  b", want: []string{"a", " ", " ", "b"}},
		{name: "mixed whitespace exploded", in: "a \t b", want: []string{"a", " ", "\t", " ", "b"}},
		{name: "leading and repeated whitespace exploded", in: "  multiple   spaces  here", want: []string{" ", " ", "multiple", " ", " ", " ", "spaces", " ", " ", "here"}},
		{name: "hyphen", in: "derek-anderson", want: []string{"derek", "-", "a", "n", "d", "e", "r", "s", "o", "n"}},
		{name: "underscore", in: "derek_anderson", want: []string{"derek", "_", "a", "n", "d", "e", "r", "s", "o", "n"}},
		{name: "slash", in: "derek/anderson", want: []string{"derek", "/", "a", "n", "d", "e", "r", "s", "o", "n"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Split(tt.in); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("Split(%q) = %#v, want %#v", tt.in, got, tt.want)
			}
		})
	}
}

func TestCustomModel(t *testing.T) {
	lm, err := NewLanguageModelFile("test_lang.txt.gz")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := lm.Split("derek"), []string{"der", "ek"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Split custom model = %#v, want %#v", got, want)
	}
}
