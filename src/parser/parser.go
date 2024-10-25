package parser

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/purarue/spkglist/src/lexer"
)

func ParseBuffer(r io.ReadCloser) (result []string, pErr error) {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		pErr = err
		return
	}
	// create lexer
	lex := lexer.NewSpkglist(antlr.NewInputStream(string(buf)))
	// debug error listener
	// TODO: attach a better error listener here, see NextToken definition
	// lex.AddErrorListener(antlr.NewDiagnosticErrorListener(true))

	for {
		t := lex.NextToken()
		tokenType := t.GetTokenType()
		if tokenType == antlr.TokenEOF {
			break
		}
		if tokenType == lexer.SpkglistQUOTED_PACKAGE {
			result = append(result, strings.Trim(t.GetText(), "`"))
		} else if tokenType == lexer.SpkglistBARE_PACKAGE {
			result = append(result, strings.TrimSpace(t.GetText()))
		}
	}
	return
}

func ParseFile(filepath string) ([]string, error) {
	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return nil, errors.New(fmt.Sprintf("File '%s' does not exist", filepath))
	}
	f, err := os.Open(filepath)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not open file at '%s'", filepath))
	}
	defer f.Close()
	return ParseBuffer(f)
}
