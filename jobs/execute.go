package jobs

import (
	"errors"
	"strconv"

	"github.com/idfly/darius"
)

import (
	"github.com/shagabutdinov/shell"
)

func Execute(state darius.State, task map[interface{}]interface{}) error {
	report, err := execute(state, task)
	if err != nil {
		if report {
			state.Log(darius.LogCommandFail, err.Error())
		}

		return err
	}

	return nil
}

func execute(
	state darius.State,
	task map[interface{}]interface{},
) (bool, error) {
	_, ok := task["command"]
	if !ok {
		return true, errors.New("command should be defined in task")
	}

	var err error
	task["command"], err = state.Expand(task["command"], false)
	if err != nil {
		return true, err
	}

	array, ok := task["command"].([]interface{})
	if ok {
		for index, element := range array {
			array[index], err = state.Expand(element, false)
			if err != nil {
				return true, err
			}

			err := state.Call("call", darius.CreateTask(array[index]))
			if err != nil {
				return false, err
			}
		}

		return false, nil
	}

	mapping, ok := task["command"].(map[interface{}]interface{})
	if ok {
		return false, state.Call("call", mapping)
	}

	str, ok := task["command"].(string)
	if !ok {
		return true, errors.New("command should be string, array or map")
	}

	state.Log(darius.LogCommand, str)
	status, err := state.Execute(
		str,
		func(kind shell.MessageType, message string) error {
			if kind == shell.StdOut {
				state.Log(darius.LogStdOut, message)
			} else if kind == shell.StdErr {
				state.Log(darius.LogStdErr, message)
			} else {
				return errors.New("Unknown message type %d")
			}

			return nil
		},
	)

	if err != nil {
		return true, err
	}

	if status != 0 {
		return true, errors.New("command execution failed: non-zero exit " +
			"status " + strconv.Itoa(status) + " received")
	}

	return false, nil
}
