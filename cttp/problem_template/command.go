package problem_template

import (
	"errors"
	"github.com/merisho/comprog/cttp/commands"
	"github.com/merisho/comprog/cttp/flags"
	"strings"
)

var (
	fileNameIsNotProvidedError = errors.New("file name is not provided")
)

func init() {
	args := map[string]string{
		"-s": "boolean, specifies a simple template that has no multiple test cases",
	}
	commands.RegisterCommand("tp", "tp <file name> [arguments], generates problem solution template", args, exec)
}

type Args struct {
	Simple bool `flag:"s"`
}

func exec(argsStr string) error {
	argsStr = strings.TrimSpace(argsStr)
	if argsStr == "" {
		return fileNameIsNotProvidedError
	}

	args := &Args{}
	if err := flags.Parse(argsStr, args); err != nil {
		return err
	}

	firstSpace := strings.Index(argsStr, " ")
	fileName := argsStr
	if firstSpace != -1 {
		fileName = argsStr[:firstSpace]
	}

	var err error
	if args.Simple {
		err = ProblemTemplate(fileName)
	} else {
		err = ProblemTemplateWithMultipleTests(fileName)
	}

	return err
}
