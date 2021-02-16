package commands

import "fmt"

type ExecFunc func(args string) error

var commands = make(map[string]command)

func Usage() string {
	u := ""
	for _, c := range commands {
		u += c.usage() + "\n"
	}

	if len(u) == 0 {
		return u
	}

	return u[:len(u) - 1]
}

func RegisterCommand(alias, description string, args map[string]string, exec ExecFunc) {
	if _, ok := commands[alias]; ok {
		panic(fmt.Sprintf("command with alias %s already registered", alias))
	}

	commands[alias] = command{
		exec: exec,
		alias: alias,
		descr: description,
		args: args,
	}
}

func Command(alias string) ExecFunc {
	return commands[alias].exec
}

type command struct {
	exec ExecFunc
	alias string
	descr string
	args map[string]string
}

func (c command) usage() string {
	u := fmt.Sprintf("%s - %s", c.alias, c.descr)
	for arg, descr := range c.args {
		u += fmt.Sprintf("\n\t\t%s - %s", arg, descr)
	}

	return u
}
