package tasks

import (
	"errors"
	"log"

	"github.com/idfly/darius"

	"github.com/shagabutdinov/shell"
)

func Call(state darius.State, task map[interface{}]interface{}) error {
	newState, err := state.Spawn(task)
	if err != nil {
		state.Log(darius.LogCommandFail, err.Error())
		return err
	}

	task = newState.Task()

	defer func() {
		err := newState.Destroy()
		if err != nil {
			log.Println(err)
		}
	}()

	ok, err := runCheckContext(newState, task)
	if err == nil {
		if !ok {
			return nil
		}

		var report bool
		report, err = runMapping(newState, task)
		if err != nil && report {
			newState.Log(darius.LogCommandFail, err.Error())
		}
	}

	err = runTail(newState, task, err)
	if err != nil {
		return err
	}

	return nil
}

func runCheckContext(
	state darius.State,
	task map[interface{}]interface{},
) (bool, error) {
	raw, ok := task["context"]
	if !ok {
		return true, nil
	}

	context, ok := raw.(string)
	if !ok {
		err := errors.New("context should be string")
		state.Log(darius.LogCommandFail, err.Error())
		return false, err
	}

	state.Log(darius.LogContext, context)
	status, err := state.Execute(
		context,
		func(kind shell.MessageType, message string) error {
			return errors.New("context must not print to stdout or " +
				"stderr; any output counted as error; received output: " +
				message)
		},
	)

	if err != nil {
		state.Log(darius.LogCommandFail, err.Error())
		return false, err
	}

	return status == 0, nil
}

func runMapping(
	state darius.State,
	task map[interface{}]interface{},
) (bool, error) {
	raw, ok := task["job"]
	if !ok {
		_, ok := task["command"]
		if !ok {
			return true, errors.New("\"job\" or \"command\" should be " +
				"defined in task")
		}

		return false, state.Call("execute", task)
	}

	str, ok := raw.(string)
	if !ok {
		return true, errors.New("job should be string")
	}

	if raw == "call" {
		return true, errors.New("call is system name and can not be used")
	}

	return false, state.Call(str, task)
}

func runTail(
	state darius.State,
	task map[interface{}]interface{},
	err error,
) error {
	if err != nil {
		rescue, ok := task["rescue"]
		if ok {
			state.Log(darius.LogRescue, "[rescue]")
			err = state.Call("call", darius.CreateTask(rescue))
		}
	}

	ensure, ok := task["ensure"]
	if ok {
		state.Log(darius.LogEnsure, "[ensure]")
		ensureErr := state.Call("call", darius.CreateTask(ensure))
		if ensureErr != nil {
			err = ensureErr
		}

	}

	return err
}
