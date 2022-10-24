package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/MichaelDiBernardo/monkey/lexer"
	"github.com/MichaelDiBernardo/monkey/parser"
)

type command struct {
	cmdarg string // arg to invoke this command at the command line
	short  string // Terse usage description (77 chars or fewer)
	run    func() // function to run for this command
}

var commands = []command{
	{"run", "[filename.monkey] will run the monkey program in the given file.", run},
	{"repl", "will start a monkey read-evaluate-print loop.", repl},
}

func main() {
	if len(os.Args) <= 1 {
		printUsage()
		return
	}

	cmdarg := os.Args[1]
	for _, cmd := range commands {
		if cmdarg == cmd.cmdarg {
			cmd.run()
			break
		}
	}
}

func printUsage() {
	fmt.Print(`monkey is a toy programming language. Enjoy playing with it 🐒

Usage:

 monkey <command> [arguments]

`)

	for _, cmd := range commands {
		fmt.Printf(" 🐵 monkey %s %s\n", cmd.cmdarg, cmd.short)
	}
}

func run() {
	args := os.Args[2:]

	rfatal := func(msg string) {
		fatal("run", msg)
	}

	if len(args) != 1 {
		rfatal(fmt.Sprintf("expected [filename.monkey], got %q\n", strings.Join(args, " ")))
	}

	srcpath := args[0]

	if strings.ToLower(filepath.Ext(srcpath)) != ".monkey" {
		rfatal(fmt.Sprintf("expected [filename.monkey], got %s\n", srcpath))
	}

	abspath, err := filepath.Abs(srcpath)

	if err != nil {
		rfatal(fmt.Sprintf("error getting abspath for %s: %v\n", srcpath, err))
	}

	if _, err := os.Stat(abspath); err != nil {
		rfatal(fmt.Sprintf("file %s does not exist\n", abspath))
	}

	lex, err := lexer.NewFromPath(abspath)

	if err != nil {
		rfatal(fmt.Sprintf("could not run %s: %v\n", abspath, err))
	}

	parse := parser.New(lex)

	program := parse.ParseProgram()

	if perrs := parse.Errors(); len(perrs) > 0 {
		var out bytes.Buffer
		out.WriteString("found parse errors\n\n")
		for i, perr := range perrs {
			loc := perr.Location
			out.WriteString(fmt.Sprintf("[%d] In %s (line %d, col %d): %s\n", i+1, loc.Path, loc.LineN, loc.CharN, perr.Message))
		}
		rfatal(out.String())
	}

	fmt.Print(program.String(), "\n")
}

func fatal(cmd string, msg string) {
	fmt.Fprintf(os.Stderr, "🙈 monkey %s: %s", cmd, msg)
	os.Exit(1)
}

func repl() {

}
