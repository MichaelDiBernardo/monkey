package main

import (
	"fmt"
	"os"
)

type command struct {
	cmdarg string // arg to invoke this command at the command line
	short  string // Terse usage description (77 chars or fewer)
	run    func() // function to run for this command
}

var commands = []command{
	{"run", "<filename.monkey> will run the monkey program in the given file.", run},
	{"repl", "will start a monkey read-evaluate-print loop.", repl},
}

func main() {
	if len(os.Args) <= 1 {
		printUsage()
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

}

func repl() {

}
