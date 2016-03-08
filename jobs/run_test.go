package jobs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunCallsCall(test *testing.T) {
	state := newState()
	state.On("Config").Return(map[interface{}]interface{}{
		"tasks": map[interface{}]interface{}{"TASK": "TASK"},
	})

	task := map[interface{}]interface{}{"command": "TASK"}
	state.On("Call", "call", task).Return(nil)
	err := Run(state, map[interface{}]interface{}{"task": "TASK"})
	assert.NoError(test, err)
	state.AssertExpectations(test)
}
