package jobs

import (
	"errors"

	"github.com/idfly/darius"
)

func RunUserTask(state darius.State, task map[interface{}]interface{}) error {
	err, report := runUserTask(state, task)
	if err != nil {
		if report {
			state.Log(darius.LogCommandFail, err.Error())
		}

		return err
	}

	return nil
}

func runUserTask(
	state darius.State,
	task map[interface{}]interface{},
) (error, bool) {
	_, ok := task["task-name"]
	if !ok {
		return errors.New("\"task-name\" should be defined in task"), true
	}

	var err error
	task["task-name"], err = state.Expand(task["task-name"], false)
	if err != nil {
		return err, true
	}

	str, ok := task["task-name"].(string)
	if !ok {
		return errors.New("task name should be string"), true
	}

	userTask, ok := state.Config()["tasks"].(map[interface{}]interface{})[str]
	if !ok {
		return errors.New("task " + str + " not found"), true
	}

	return state.Call("call", darius.CreateTask(userTask)), false
}
