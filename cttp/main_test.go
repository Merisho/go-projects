package main

import (
    "github.com/merisho/comprog/cttp/commands"
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestNoCommandError(t *testing.T) {
    err := Run([]string{})
    assert.EqualError(t, err, "command is not specified")
}

func TestUnknownCommand(t *testing.T) {
    err := Run([]string{"unknown"})
    assert.EqualError(t, err, "unknown command")
}


func TestCommandRun(t *testing.T) {
   var (
       runCommand bool
       commandArgs string
   )
   commands.RegisterCommand("test", "", nil, func(args string) error {
       runCommand = true
       commandArgs = args
       return nil
   })

   args := []string{"test", "-commandArgs", "arg", "value", "-arg2", "111"}
   err := Run(args)

   assert.NoError(t, err)
   assert.True(t, runCommand)
   assert.Equal(t, "-commandArgs arg value -arg2 111", commandArgs)
}
