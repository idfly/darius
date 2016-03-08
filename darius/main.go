package main

import (
	"os"

	"github.com/shagabutdinov/arguments"
)

var (
	options = arguments.Arguments{
		"config": arguments.Argument{
			"config",
			"configuration file",
			arguments.String,
			"c",
			false,
			nil,
			false,
		},

		"help": arguments.Argument{
			"help",
			"displays help",
			arguments.Flag,
			"h",
			false,
			nil,
			false,
		},

		"local": arguments.Argument{
			"local",
			"call all tasks locally",
			arguments.Flag,
			"l",
			false,
			nil,
			false,
		},

		"tail": arguments.Argument{
			"command",
			"command and its options to execute",
			arguments.Tail,
			"",
			false,
			nil,
			false,
		},
	}
)

func main() {
	state, err := newState()
	if err != nil {
		panic(err)
	}

	err = call(state, os.Args[1:])
	if err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
