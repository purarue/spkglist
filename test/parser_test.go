package parser_test

import (
	"io/ioutil"
	"strings"

	"testing"

	"github.com/purarue/spkglist/src/parser"
)

func runParser(contents string) []string {
	parsed, err := parser.ParseBuffer(ioutil.NopCloser(strings.NewReader(contents)))
	if err != nil {
		panic(err)
	}
	return parsed
}

func assertSlice(s1 []string, s2 []string, t *testing.T) {
	if len(s1) != len(s2) {
		t.Errorf("Slices are not the same length\n")
	}
	i := 0
	for {
		if s1[i] != s2[i] {
			t.Errorf("Item at index %d differs: '%s' '%s'\n", i, s1[i], s2[i])
		}
		i += 1
		if i == len(s1) {
			break
		}
	}
}

func TestParserBasic(t *testing.T) {
	results := runParser(`
golang
java
something 

`)
	expected := []string{"golang", "java", "something"}
	assertSlice(results, expected, t)
}

func TestParserComments(t *testing.T) {
	results := runParser(`
# some other comment

  # ignored
java#no spaces
something   # same line

`)
	expected := []string{"java", "something"}
	assertSlice(results, expected, t)
}

// perhaps this should be changed. for now this just reports the error to stderr
// and tries to parse the rest of the file
func TestParserIgnoresErrors(t *testing.T) {
	results := runParser(`
# this should  fail
🎵  # ignored
succeeds
`)
	expected := []string{"succeeds"}
	assertSlice(results, expected, t)
}

func TestParserQuotedPackage(t *testing.T) {
	results := runParser("`quoted#master` # comment\n\nunquoted\n#comment line")
	expected := []string{"quoted#master", "unquoted"}
	assertSlice(results, expected, t)
}
