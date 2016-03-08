package jobs

import (
	"errors"
	"testing"

	"github.com/idfly/darius"

	"github.com/stretchr/testify/assert"
)

func TestCallRunsString(test *testing.T) {
	state := newState()
	command := map[interface{}]interface{}{"command": "CMD"}
	state.On("Call", "execute", command).Return(nil)
	err := Call(state, command)
	assert.NoError(test, err)
	state.AssertExpectations(test)
}

func TestCallReturnsStringError(test *testing.T) {
	state := newState()
	command := map[interface{}]interface{}{"command": "CMD"}
	state.On("Call", "execute", command).Return(errors.New("ERROR"))
	err := Call(state, command)
	assert.Error(test, err)
	state.AssertExpectations(test)
}

func TestCallReturnsErrorOnUnknownTaskType(test *testing.T) {
	state := newState()
	state.On("Log", darius.LogCommandFail, "\"job\" or \"command\" should be "+
		"defined in task")
	err := Call(state, map[interface{}]interface{}{})
	assert.Error(test, err)
}

func TestCallChecksContext(test *testing.T) {
	state := newState()
	command := map[interface{}]interface{}{"command": "CMD", "context": "CTX"}
	state.On("Log", darius.LogContext, "CTX").Return(nil)
	state.On("Execute", "CTX").Return(0, nil)
	state.On("Call", "execute", command).Return(nil)
	err := Call(state, command)
	assert.NoError(test, err)
	state.AssertExpectations(test)
}

func TestCallNotExecutesCommandIfContextCheckFailed(test *testing.T) {
	state := newState()
	command := map[interface{}]interface{}{"command": "CMD", "context": "CTX"}
	state.On("Log", darius.LogContext, "CTX").Return(nil)
	state.On("Execute", "CTX").Return(1, nil)
	err := Call(state, command)
	assert.NoError(test, err)
	state.AssertExpectations(test)
}

func TestCallRunsMapWithTaskAndCommand(test *testing.T) {
	state := newState()
	task := map[interface{}]interface{}{"job": "execute", "command": "CMD"}
	state.On("Call", "execute", task).Return(nil)
	err := Call(state, task)
	assert.NoError(test, err)
	state.AssertExpectations(test)
}

func TestCallReturnsMapError(test *testing.T) {
	state := newState()
	command := map[interface{}]interface{}{"command": "CMD"}
	state.On("Call", "execute", command).Return(errors.New("ERROR"))
	err := Call(state, command)
	assert.Error(test, err)
	state.AssertExpectations(test)
}

func TestCallReturnsErrorOnStatError(test *testing.T) {
	state := newState()
	task := map[interface{}]interface{}{"job": "UNKNOWN"}
	state.On("Call", "UNKNOWN", task).Return(errors.New("ERROR"))
	err := Call(state, task)
	assert.Error(test, err)
	state.AssertExpectations(test)
}

func TestCallReturnsErrorIfTaskIsNotString(test *testing.T) {
	state := newState()
	state.On("Log", darius.LogCommandFail, "job should be string").Return(nil)
	err := Call(state, map[interface{}]interface{}{"job": 0})
	assert.Error(test, err)
	state.AssertExpectations(test)
}

func TestCallCallsRescueIfErrorOccured(test *testing.T) {
	state := newState()
	command := map[interface{}]interface{}{"command": "CMD", "rescue": "RESCUE"}
	state.On("Call", "execute", command).Return(errors.New("ERROR"))
	state.On("Call", "call", map[interface{}]interface{}{"command": "RESCUE"}).
		Return(nil)
	state.On("Log", darius.LogRescue, "[rescue]").Return(nil)
	err := Call(state, command)
	assert.NoError(test, err)
	state.AssertExpectations(test)
}

func TestCallReturnsRescueError(test *testing.T) {
	state := newState()
	command := map[interface{}]interface{}{"command": "CMD", "rescue": "RESCUE"}
	state.On("Call", "execute", command).Return(errors.New("ERROR1"))
	state.On("Call", "call", map[interface{}]interface{}{"command": "RESCUE"}).
		Return(errors.New("ERROR2"))
	state.On("Log", darius.LogRescue, "[rescue]").Return(nil)
	err := Call(state, command)
	assert.Error(test, err, "ERROR2")
	state.AssertExpectations(test)
}

func TestCallCallsEnsure(test *testing.T) {
	state := newState()
	command := map[interface{}]interface{}{"command": "CMD", "ensure": "ENSURE"}
	state.On("Call", "execute", command).Return(nil)
	state.On("Call", "call", map[interface{}]interface{}{"command": "ENSURE"}).
		Return(nil)
	state.On("Log", darius.LogEnsure, "[ensure]").Return(nil)
	err := Call(state, command)
	assert.NoError(test, err)
	state.AssertExpectations(test)
}

func TestCallCallsEnsureWithError(test *testing.T) {
	state := newState()
	command := map[interface{}]interface{}{"command": "CMD", "ensure": "ENSURE"}
	state.On("Call", "execute", command).Return(errors.New("ERROR"))
	state.On("Call", "call", map[interface{}]interface{}{"command": "ENSURE"}).
		Return(nil)
	state.On("Log", darius.LogEnsure, "[ensure]").Return(nil)
	err := Call(state, command)
	assert.Error(test, err)
	state.AssertExpectations(test)
}

func TestCallCallsEnsureAfterError(test *testing.T) {
	state := newState()
	command := map[interface{}]interface{}{"command": "CMD", "ensure": "ENSURE"}
	state.On("Call", "execute", command).Return(errors.New("ERROR1"))
	state.On("Call", "call", map[interface{}]interface{}{"command": "ENSURE"}).
		Return(errors.New("ERROR2"))
	state.On("Log", darius.LogEnsure, "[ensure]").Return(nil)
	err := Call(state, command)
	assert.Error(test, err, "ERROR2")
	state.AssertExpectations(test)
}
