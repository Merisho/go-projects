package main

import (
	"errors"
	"fmt"
	"github.com/merisho/comprog/cttp/commands"
	"os"
	"strings"
	_ "github.com/merisho/comprog/cttp/problem_template"
)

var (
	notSpecifiedError = errors.New("command is not specified")
	unknownError = errors.New("unknown command")
)

func main() {
	err := Run(os.Args[1:])
	if err != nil {
		handleError(err)
		os.Exit(1)
	}
}

func handleError(err error) {
	_, e := os.Stderr.WriteString(err.Error())
	if e != nil {
		panic(e)
	}

	_, e = os.Stderr.WriteString(fmt.Sprintf("\nUsage: cttp <command> [arguments]\n\t%s", commands.Usage()))
	if e != nil {
		panic(e)
	}
}

func Run(args []string) error {
	if len(args) == 0 {
		return notSpecifiedError
	}

	cmd := commands.Command(args[0])
	if cmd == nil {
		return unknownError
	}

	cmdArgs := strings.Join(args[1:], " ")
	return cmd(cmdArgs)
}
