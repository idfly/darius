package jobs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunUserTaskCallsCall(test *testing.T) {
	state := newState()
	state.On("Config").Return(map[interface{}]interface{}{
		"tasks": map[interface{}]interface{}{"TASK": "TASK"},
	})

	task := map[interface{}]interface{}{"command": "TASK"}
	state.On("Call", "call", task).Return(nil)
	err := RunUserTask(state, map[interface{}]interface{}{"task-name": "TASK"})
	assert.NoError(test, err)
	state.AssertExpectations(test)
}
