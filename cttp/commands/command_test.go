package commands

import (
    "fmt"
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestUsage(t *testing.T) {
    testArgs := map[string]string{
        "-testArg1": "this is a test argument 1",
        "-testArg2": "this is a test argument 2",
    }
    RegisterCommand("test", "test usage", testArgs, func(args string) error { return nil })

    anotherCommandArgs := map[string]string{
        "-testArg1": "this is a test argument 1",
        "-testArg2": "this is a test argument 2",
    }
    RegisterCommand("another", "another command", anotherCommandArgs, func(args string) error { return nil })

    usage := Usage()
    expected := fmt.Sprintf("test - test usage\n\t\t-testArg1 - this is a test argument 1\n\t\t-testArg2 - this is a test argument 2\n")
    expected += fmt.Sprintf("another - another command\n\t\t-testArg1 - this is a test argument 1\n\t\t-testArg2 - this is a test argument 2")
    assert.Equal(t, expected, usage)
}
