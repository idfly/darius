package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunReportsErrorIfNoTaskProvided(test *testing.T) {
	state, utils := newTestState(true)
	defer state.Destroy()
	utils.On("readFile", ".darius.yml").Return("tasks: {}", nil)
	utils.On("err", "task must be set in command line options; use --help "+
		"to receive help", true)
	utils.On("out", "\x1b[1;37;41m ** task execution failed (check logs for "+
		"details) ** \x1b[0m", true)
	err := call(state, []string{})
	assert.Error(test, err, "task must be set")
}

func TestRunRunsTask(test *testing.T) {
	state, utils := newTestState(true)
	defer state.Destroy()
	utils.On("readFile", ".darius.yml").Return("tasks: {T: TASK}", nil)
	utils.On("call", map[interface{}]interface{}{"command": "TASK"}).Return(nil)
	utils.On("out", "\x1b[1;37;42mtask completed\x1b[0m", true)
	err := call(state, []string{"T"})
	assert.NoError(test, err)
	utils.AssertExpectations(test)
}
