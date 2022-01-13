package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	openEditorFlag  bool
	printConfigFlag bool
	profileNameFlag string
)

func extractCliVars() ([]string, map[string]string) {
	cliArgs := make([]string, 0)
	cliVars := make(map[string]string)

	for _, arg := range flag.Args() {

		if !strings.Contains(arg, "=") {
			cliArgs = append(cliArgs, arg)
			continue
		}

		pair := strings.SplitN(arg, "=", 2)
		key, val := pair[0], pair[1]

		if strings.HasPrefix(val, "\"") && strings.HasSuffix(val, "\"") {
			val = val[1 : len(val)-1]
		} else if strings.HasPrefix(val, "'") && strings.HasSuffix(val, "'") {
			val = val[1 : len(val)-1]
		}

		cliVars[key] = val
	}
	return cliArgs, cliVars
}

func parseFlags() {
	flag.BoolVar(
		&printConfigFlag,
		"c",
		false,
		"print the default configuration to STDOUT",
	)
	flag.BoolVar(
		&openEditorFlag,
		"e",
		false,
		"open the editor after writing the file",
	)
	flag.StringVar(
		&profileNameFlag,
		"p",
		"default",
		"name of the profile to use",
	)

	flag.Usage = func() {
		fmt.Fprintf(
			os.Stderr,
			`USAGE: %s [-c] [-e] [-p PROFILE] INPUT

INPUT must be a list of arguments or - to get the input from the STDIN.
The input can also contains arguments using key=value format to pass variables for the template.

OPTIONS:
`,
			APP_NAME,
		)
		flag.PrintDefaults()
	}

	flag.Parse()
}
