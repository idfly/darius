package main

import (
	"errors"
	"os"
	"strings"

	"github.com/idfly/darius"
	"github.com/shagabutdinov/arguments"
)

func call(state *state, optionsArray []string) error {
	arguments, err := options.Parse(optionsArray)
	if err != nil {
		state.utils.err(err.Error(), true)
		os.Exit(1)
	}

	err = runTask(state, arguments)

	if err != nil {
		state.Log(darius.LogTaskFail, " ** task execution failed (check logs "+
			"for details) ** ")
	} else {
		state.Log(darius.LogTaskSuccess, "task completed")
	}

	return err
}

func runTask(state *state, arguments arguments.Values) error {
	tail, _, err := arguments.Strings("tail", []string{})
	if len(tail) > 0 && strings.HasPrefix(tail[0], "-") {
		return errors.New("unknown option " + tail[0])
	}

	state.argv = tail[1:]

	state.runLocally, _, err = arguments.Boolean("local", false)
	if err != nil {
		return err
	}

	help, _, err := arguments.Boolean("help", false)
	if err != nil {
		return err
	}

	if help {
		return errors.New("help is not available yet")
	}

	configLoader := darius.Config{
		ReadFile: state.utils.readFile,
		Glob:     state.utils.glob,
	}

	configFile, _, err := arguments.String("config", ".darius.yml")
	if err != nil {
		return err
	}

	config, err := configLoader.Load(configFile)
	if err != nil {
		return err
	}

	state.config = config

	tasksRaw, ok := config["tasks"]
	if !ok {
		return errors.New("tasks section must be set in config")
	}

	tasksExpanded, err := state.Expand(tasksRaw, false)
	if err != nil {
		return err
	}

	tasks, ok := tasksExpanded.(map[interface{}]interface{})
	if !ok {
		return errors.New("tasks section must be map")
	}

	if len(tail) == 0 {
		state.utils.err("task must be set in command line options; use "+
			"--help to receive help", true)
		return errors.New("no task provided")
	}

	task, ok := tasks[tail[0]]
	if !ok {
		return errors.New("task " + tail[0] + " not found in configuration " +
			"file")
	}

	err = state.Call("call", darius.CreateTask(task))
	if err != nil {
		return err
	}

	return nil
}
