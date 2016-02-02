package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStateExpandReturnsNonString(test *testing.T) {
	state, _ := newTestState(false)
	defer state.Destroy()
	result, err := state.Expand(true, false)
	assert.NoError(test, err)
	assert.Equal(test, true, result)
}

func TestStateExpandReturnsEscapedExpression(test *testing.T) {
	state, _ := newTestState(false)
	defer state.Destroy()
	result, err := state.Expand("$${test}", false)
	assert.NoError(test, err)
	assert.Equal(test, "${test}", result)
}

func TestStateExpandExpandsVariable(test *testing.T) {
	state, _ := newTestState(false)
	defer state.Destroy()
	state.config = map[interface{}]interface{}{
		"vars": map[interface{}]interface{}{"var": "VARIABLE"},
	}

	result, err := state.Expand("${vars.var}", false)
	assert.NoError(test, err)
	assert.Equal(test, "VARIABLE", result)
}

func TestStateExpandPreventsRecursion(test *testing.T) {
	state, _ := newTestState(false)
	defer state.Destroy()
	state.config = map[interface{}]interface{}{
		"vars": map[interface{}]interface{}{"var": "${vars.var}"},
	}

	_, err := state.Expand("${vars.var}", false)
	assert.Error(test, err)
}

func TestStateExpandExpandsSubVariable(test *testing.T) {
	state, _ := newTestState(false)
	defer state.Destroy()
	state.config = map[interface{}]interface{}{
		"vars": map[interface{}]interface{}{
			"VAR1": map[interface{}]interface{}{"VAR2": "SUBVAR"},
		},
	}

	result, err := state.Expand("${vars.VAR1.VAR2}", false)
	assert.NoError(test, err)
	assert.Equal(test, "SUBVAR", result)
}

func TestStateExpandExpandsVariableFromTask(test *testing.T) {
	state, _ := newTestState(false)
	defer state.Destroy()
	state.config = map[interface{}]interface{}{
		"vars": map[interface{}]interface{}{"var": "WRONG"},
	}

	state.task = map[interface{}]interface{}{
		"vars": map[interface{}]interface{}{
			"var": "VARIABLE",
		},
	}

	result, err := state.Expand("${vars.var}", false)
	assert.NoError(test, err)
	assert.Equal(test, "VARIABLE", result)
}

func TestStateExpandExpandsVariableFromParentTask(test *testing.T) {
	testState, _ := newTestState(false)
	defer testState.Destroy()
	testState.config = map[interface{}]interface{}{
		"vars": map[interface{}]interface{}{"var": "WRONG"},
	}

	testState.parent = &state{
		task: map[interface{}]interface{}{
			"vars": map[interface{}]interface{}{
				"var": "VARIABLE",
			},
		},
	}

	testState.task = map[interface{}]interface{}{
		"vars": map[interface{}]interface{}{
			"WRONG": "WRONG",
		},
	}

	result, err := testState.Expand("${vars.var}", false)
	assert.NoError(test, err)
	assert.Equal(test, "VARIABLE", result)
}

func TestStateExpandReportsUnknownVariable(test *testing.T) {
	state, _ := newTestState(false)
	defer state.Destroy()
	state.config = map[interface{}]interface{}{
		"vars": map[interface{}]interface{}{"var": "WRONG"},
	}

	_, err := state.Expand("${vars.UNKNOWN}", false)
	assert.Error(test, err)
}

func TestStateExpandExpandsArguments(test *testing.T) {
	oldState, _ := newTestState(false)
	defer oldState.Destroy()
	oldState.argv = []string{"-a", "VALUE"}
	task := map[interface{}]interface{}{
		"args": map[interface{}]interface{}{
			"arg": map[interface{}]interface{}{
				"name":     "arg",
				"type":     "string",
				"shortcut": "a",
			},
		},
	}

	newState, err := oldState.Spawn(task)
	assert.NoError(test, err)
	result, err := newState.Expand("${args.arg}", false)
	assert.NoError(test, err)
	assert.Equal(test, "VALUE", result)
}

func TestStateExpandReturnsErrorIfRequiredArgumentIsNotSet(test *testing.T) {
	oldState, _ := newTestState(false)
	defer oldState.Destroy()
	oldState.argv = []string{}
	task := map[interface{}]interface{}{
		"args": map[interface{}]interface{}{
			"arg": map[interface{}]interface{}{
				"type":     "string",
				"required": true,
			},
		},
	}

	_, err := oldState.Spawn(task)
	assert.Error(test, err)
}

func TestStateDoesNotSpawnConnectionWhenRunningLocally(test *testing.T) {
	oldState, _ := newTestState(false)
	defer oldState.Destroy()
	oldState.runLocally = true
	task := map[interface{}]interface{}{"host": "user@example.com"}
	newState, err := oldState.Spawn(task)
	assert.NoError(test, err)
	assert.Equal(test, nil, newState.(*state).shell)
}
