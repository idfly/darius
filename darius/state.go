package main

import (
	"errors"
	"io/ioutil"
	"os"
	"os/user"

	"github.com/idfly/darius"
	"github.com/idfly/darius/tasks"

	"github.com/shagabutdinov/arguments"
	"github.com/shagabutdinov/shell"

	"golang.org/x/crypto/ssh"
)

const (
	logLineLimit = 2048
)

var (
	tasksFuncs = map[string]func(
		darius.State,
		map[interface{}]interface{},
	) error{
		"call":    tasks.Call,
		"execute": tasks.Execute,
		"run":     tasks.Run,
	}
)

func newState() (*state, error) {
	shell, err := shell.NewLocal(shell.LocalConfig{LineLimit: 1024})
	if err != nil {
		return nil, err
	}

	state := &state{
		args:  map[interface{}]interface{}{},
		tasks: tasksFuncs,
		shell: shell,
		utils: &utils{
			stdout: os.Stdout,
			stderr: os.Stderr,
		},
	}

	state.expression = darius.NewExpression(state, state.expandExpression)

	return state, nil
}

type state struct {
	config     map[interface{}]interface{}
	argv       []string
	args       map[interface{}]interface{}
	runLocally bool

	shell      shell.Shell
	utils      utilsInterface
	expression *darius.Expression

	tasks map[string]func(darius.State, map[interface{}]interface{}) error
	task  map[interface{}]interface{}

	level  int
	parent *state
}

func (state *state) Config() map[interface{}]interface{} {
	return state.config
}

func (state *state) Task() map[interface{}]interface{} {
	return state.task
}

func (state *state) Args() map[interface{}]interface{} {
	return state.args
}

func (state *state) Parent() (darius.State, bool) {
	return state.parent, state.parent != nil
}

func (state *state) Execute(
	command string,
	handler func(shell.MessageType, string) error,
) (int, error) {
	var shell shell.Shell = nil
	current := state

	for {
		shell = current.shell
		if shell != nil {
			break
		}

		if current.parent == nil {
			return -1, errors.New("all shells closed")
		}

		current = current.parent
	}

	return shell.Run(command, handler)
}

func (state *state) Expand(
	value interface{},
	recursive bool,
) (interface{}, error) {
	return state.expression.Expand(value, recursive)
}

func (state *state) expandExpression(
	kind string,
	expr string,
) (interface{}, error) {
	return nil, errors.New("unknown expression: " + expr)
}

func (state *state) Log(level darius.LogLevel, message string) {
	format, ok := formatters[level]
	if !ok {
		panic("unknown log level")
	}

	state.utils.out(format(state.level, message), true)
}

func (state *state) Call(task string, value map[interface{}]interface{}) error {
	return state.tasks[task](state, value)
}

func (oldState *state) Spawn(
	task map[interface{}]interface{},
) (darius.State, error) {
	task = darius.Copy(task).(map[interface{}]interface{})

	result := &state{
		config:     oldState.config,
		runLocally: oldState.runLocally,
		argv:       oldState.argv,
		args:       darius.Copy(oldState.args).(map[interface{}]interface{}),
		parent:     oldState,
		tasks:      oldState.tasks,
		utils:      oldState.utils,
		level:      oldState.level,
		task:       task,
	}

	result.expression = darius.NewExpression(result, result.expandExpression)

	err := result.createArgs(task)
	if err != nil {
		return nil, err
	}

	err = darius.ExpandTask(result, result.task)
	if err != nil {
		return nil, err
	}

	err = result.reportName(task)
	if err != nil {
		return nil, err
	}

	err = result.createShell(task)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (state *state) reportName(task map[interface{}]interface{}) error {
	_, ok := task["name"]
	if !ok {
		return nil
	}

	name, ok := task["name"].(string)
	if !ok {
		return errors.New("task name should be string")
	}

	state.Log(darius.LogName, name)
	state.level += 1
	return nil
}

func (state *state) createShell(task map[interface{}]interface{}) error {
	if state.runLocally {
		return nil
	}

	raw, ok := task["host"]
	if !ok {
		return nil
	}

	host := ""
	hostMapping, ok := raw.(map[interface{}]interface{})
	if !ok {
		hostString, ok := raw.(string)
		if !ok {
			return errors.New("host should be string or map")
		}

		host = hostString
	} else {
		hostRaw, ok := hostMapping["host"]
		if !ok {
			return errors.New("host must be set in host section")
		}

		host, ok = hostRaw.(string)
		if !ok {
			return errors.New("host must be string")
		}
	}

	keyFile := ""
	keyFileRaw, ok := hostMapping["key"]
	if ok {
		keyFile, ok = keyFileRaw.(string)
		if !ok {
			return errors.New("keyfile should be string")
		}
	} else {
		usr, err := user.Current()
		if err != nil {
			return err
		}

		keyFile = usr.HomeDir + "/.ssh/id_rsa"
	}

	keyContents, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return err
	}

	key, err := ssh.ParsePrivateKey(keyContents)
	if err != nil {
		return err
	}

	state.Log(darius.LogSystem, "connecting to "+host+"...")
	state.shell, err = shell.NewRemote(shell.RemoteConfig{
		Address:   host,
		Auth:      []ssh.AuthMethod{ssh.PublicKeys(key)},
		LineLimit: logLineLimit,
	})

	if err != nil {
		return errors.New("failed to connect to " + host + ": " + err.Error())
	} else {
		state.Log(darius.LogSystem, "connection established")
	}

	return nil
}

func (state *state) createArgs(task map[interface{}]interface{}) error {
	_, ok := task["params"]
	if ok {
		var err error
		task["params"], err = state.Expand(task["params"], true)
		if err != nil {
			return err
		}

		mapping, ok := task["params"].(map[interface{}]interface{})
		if !ok {
			return errors.New("params should be map")
		}

		for key, value := range mapping {
			state.args[key] = value
		}
	}

	_, ok = task["args"]
	if !ok {
		return nil
	}

	var err error
	task["args"], err = state.Expand(task["args"], true)
	if err != nil {
		return err
	}

	mapping, ok := task["args"].(map[interface{}]interface{})
	if !ok {
		return errors.New("args must be map")
	}

	targetArgs := map[interface{}]interface{}{}
	for key, value := range mapping {
		_, ok := state.args[key]
		if ok {
			continue
		}

		targetArgs[key] = value
	}

	argsInfo, err := arguments.Create(targetArgs)
	if err != nil {
		return err
	}

	args, err := argsInfo.Parse(state.argv)
	if err != nil {
		return err
	}

	for key, value := range args {
		state.args[key] = value
	}

	return nil
}

func (state *state) Destroy() error {
	if state.shell != nil {
		err := state.shell.Close()
		if err != nil {
			return err
		}
	}

	return nil
}
