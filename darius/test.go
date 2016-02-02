package main

import (
	"github.com/idfly/darius"

	"github.com/shagabutdinov/shell"
)

func newTestState(mockCall bool) (*state, *utilsMock) {
	shell, err := shell.NewLocal(shell.LocalConfig{LineLimit: 1024})
	if err != nil {
		panic(err)
	}

	utils := &utilsMock{}
	state := &state{
		tasks: map[string]func(darius.State, map[interface{}]interface{}) error{
			"call": utils.call,
		},
		shell:  shell,
		utils:  utils,
		parent: nil,
	}

	if !mockCall {
		state.tasks = tasksFuncs
	}

	state.expression = darius.NewExpression(state, state.expandExpression)

	return state, utils
}
