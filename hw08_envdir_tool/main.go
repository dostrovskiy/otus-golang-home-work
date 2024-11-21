package main

import (
	"errors"
	"log"
	"os"
)

// ErrNotEnoughArgs is returned when not enough arguments.
// The util expects at least 2 arguments: directory with environment variables files and command for execution.
var ErrNotEnoughArgs = errors.New("not enough arguments")

func main() {
	if len(os.Args) < 3 {
		log.Fatalf("execution error: %+v\n", ErrNotEnoughArgs)
	}
	envdir := os.Args[1]
	cmd := os.Args[2:]
	env, err := ReadDir(envdir)
	if err != nil {
		log.Fatalf("execution error: %+v\n", err)
	}

	ret := RunCmd(cmd, env)
	os.Exit(ret)
}
