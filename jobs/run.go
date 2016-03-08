package jobs

import (
	"errors"

	"github.com/idfly/darius"
)

func Run(state darius.State, task map[interface{}]interface{}) error {
	err, report := run(state, task)
	if err != nil {
		if report {
			state.Log(darius.LogCommandFail, err.Error())
		}

		return err
	}

	return nil
}

func run(state darius.State, task map[interface{}]interface{}) (error, bool) {
	_, ok := task["task"]
	if !ok {
		return errors.New("\"task\" should be defined in task"), true
	}

	var err error
	task["task"], err = state.Expand(task["task"], false)
	if err != nil {
		return err, true
	}

	str, ok := task["task"].(string)
	if !ok {
		return errors.New("task name should be string"), true
	}

	userTask, ok := state.Config()["tasks"].(map[interface{}]interface{})[str]
	if !ok {
		return errors.New("task " + str + " not found"), true
	}

	return state.Call("call", darius.CreateTask(userTask)), false
}
