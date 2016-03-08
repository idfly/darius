package jobs

import (
	"errors"
	"testing"

	"github.com/idfly/darius"

	"github.com/stretchr/testify/assert"
)

func TestExecuteReturnsErrOnUndefinedCommand(test *testing.T) {
	state := newState()
	state.On("Log", darius.LogCommandFail, "command should be defined "+
		"in task").Return(nil)
	err := Execute(state, map[interface{}]interface{}{})
	assert.Error(test, err)
}

func TestExecuteRunsArray(test *testing.T) {
	state := newState()
	commands := []interface{}{"COMMAND1", "COMMAND2"}
	command1 := map[interface{}]interface{}{"command": "COMMAND1"}
	state.On("Call", "call", command1).Return(nil)
	command2 := map[interface{}]interface{}{"command": "COMMAND2"}
	state.On("Call", "call", command2).Return(nil)
	err := Execute(state, map[interface{}]interface{}{"command": commands})
	assert.NoError(test, err)
	state.AssertExpectations(test)
}

func TestExecuteRunsMap(test *testing.T) {
	state := newState()
	command := map[interface{}]interface{}{"command": "COMMAND"}
	state.On("Call", "call", command).Return(nil)
	err := Execute(state, map[interface{}]interface{}{"command": command})
	assert.NoError(test, err)
	state.AssertExpectations(test)
}

func TestExecuteCallsShell(test *testing.T) {
	state := newState()
	state.On("Execute", "COMMAND").Return(0, nil)
	state.On("Log", darius.LogCommand, "COMMAND").Return(nil)
	state.On("Log", darius.LogStdOut, "OUT").Return(nil)
	state.On("Log", darius.LogStdErr, "ERR").Return(nil)
	err := Execute(state, map[interface{}]interface{}{"command": "COMMAND"})
	assert.NoError(test, err)
	state.AssertExpectations(test)
}

func TestExecuteReturnsErrorOnError(test *testing.T) {
	state := newState()
	state.On("Execute", "COMMAND").Return(0, errors.New("ERROR"))
	state.On("Log", darius.LogCommand, "COMMAND").Return(nil)
	state.On("Log", darius.LogStdOut, "OUT").Return(nil)
	state.On("Log", darius.LogStdErr, "ERR").Return(nil)
	state.On("Log", darius.LogCommandFail, "ERROR").Return(nil)
	err := Execute(state, map[interface{}]interface{}{"command": "COMMAND"})
	assert.Error(test, err)
}

func TestExecuteReturnsErrorOnNonZeroExitStatus(test *testing.T) {
	state := newState()
	state.On("Execute", "COMMAND").Return(1, nil)
	state.On("Log", darius.LogCommand, "COMMAND").Return(nil)
	state.On("Log", darius.LogStdOut, "OUT").Return(nil)
	state.On("Log", darius.LogStdErr, "ERR").Return(nil)
	state.On("Log", darius.LogCommandFail, "command execution failed: "+
		"non-zero exit status 1 received").Return(nil)
	err := Execute(state, map[interface{}]interface{}{"command": "COMMAND"})
	assert.Error(test, err)
}
