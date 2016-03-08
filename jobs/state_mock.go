package jobs

import (
	"github.com/idfly/darius"

	"github.com/shagabutdinov/shell"
	"github.com/stretchr/testify/mock"
)

func newState() *state {
	state := &state{}
	return state
}

type state struct {
	mock.Mock
	task map[interface{}]interface{}
}

func (mock *state) Args() map[interface{}]interface{} {
	args := mock.Called()
	return args.Get(0).(map[interface{}]interface{})
}

func (mock *state) Task() map[interface{}]interface{} {
	return mock.task
}

func (mock *state) Config() map[interface{}]interface{} {
	args := mock.Called()
	return args.Get(0).(map[interface{}]interface{})
}

func (mock *state) Parent() (darius.State, bool) {
	args := mock.Called()
	return args.Get(0).(darius.State), true
}

func (mock *state) Execute(
	command string,
	handler func(shell.MessageType, string) error,
) (int, error) {
	handler(shell.StdOut, "OUT")
	handler(shell.StdErr, "ERR")
	args := mock.Called(command)
	return args.Int(0), args.Error(1)
}

func (mock *state) Expand(
	value interface{},
	recursive bool,
) (interface{}, error) {
	return value, nil
}

func (mock *state) Log(level darius.LogLevel, message string) {
	mock.Called(level, message)
}

func (mock *state) Spawn(
	task map[interface{}]interface{},
) (darius.State, error) {
	mock.task = task
	return mock, nil
}

func (mock *state) Call(task string, value map[interface{}]interface{}) error {
	args := mock.Called(task, value)
	return args.Error(0)
}

func (mock *state) Destroy() error {
	return nil
}
