package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/purarue/spkglist/src/parser"
)

type Packages []string

func fatal(params ...interface{}) {
	fmt.Fprintln(os.Stderr, params...)
	os.Exit(1)
}

func fatalf(m string, params ...interface{}) {
	fmt.Fprintf(os.Stderr, m, params...)
	os.Exit(1)
}

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, `A simple package list parser
Pass one or more files to read from as arguments
If none provided, reads from STDIN.`)
		flag.PrintDefaults()
	}
	split := flag.Bool("split", false, "split on all whitespace, instead of newlines")
	printJson := flag.Bool("json", false, "print results as a JSON array")
	delimiter := flag.String("delimiter", "\n", "delimiter to print between results")
	// this is false by default because including a new line as the last delimiter
	// works better for processing lines in the shell
	skiplast := flag.Bool("skip-last-delim", false, "dont print the delimiter after the last item")
	null_char := flag.Bool("print0", false, "use the null character as the delimiter")
	flag.Parse()
	positionals := flag.Args()

	// handle delimiter flag
	delim := *delimiter
	if *null_char {
		delim = "\000"
	}

	// parse, either reading from STDIN of multiple files
	var parsedLines []string
	if len(positionals) == 0 {
		// read all of STDIN
		buf, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			fatal(err)
		}
		// parse Buffer
		lines, err := parser.ParseBuffer(ioutil.NopCloser(bytes.NewReader(buf)))
		if err != nil {
			fatal(err)
		}
		parsedLines = lines
	} else {
		for _, file := range positionals {
			lines, err := parser.ParseFile(file)
			if err != nil {
				fatal(err)
			}
			parsedLines = append(parsedLines, lines...)
		}
	}

	// handle 'split' flag
	var packages Packages
	if *split {
		for _, line := range parsedLines {
			for _, pkg := range strings.Fields(line) {
				packages = append(packages, pkg)
			}
		}
	} else {
		packages = parsedLines
	}

	// print to STDOUT
	if *printJson {
		jsonBytes, err := json.Marshal(packages)
		if err != nil {
			fatal(err)
		}
		fmt.Print(string(jsonBytes))
	} else {
		for i, p := range packages {
			fmt.Print(p)
			if i != len(packages)-1 {
				fmt.Print(delim)
			} else {
				if !*skiplast {
					fmt.Print(delim)
				}
			}
		}
	}
}
