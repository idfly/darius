package darius

import "github.com/shagabutdinov/shell"

type LogLevel int

const (
	LogName LogLevel = iota
	LogSystem
	LogStdOut
	LogStdErr
	LogCommand
	LogCommandFail
	LogContext
	LogRescue
	LogEnsure
	LogTaskFail
	LogTaskSuccess
)

type State interface {
	Call(string, map[interface{}]interface{}) error
	Log(LogLevel, string)
	Execute(string, func(shell.MessageType, string) error) (int, error)
	Expand(interface{}, bool) (interface{}, error)

	Config() map[interface{}]interface{}
	Task() map[interface{}]interface{}
	Args() map[interface{}]interface{}
	Parent() (State, bool)

	Spawn(map[interface{}]interface{}) (State, error)
	Destroy() error
}
