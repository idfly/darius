package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunTaskSendsStdout(test *testing.T) {
	state, utils := newTestState(false)
	defer state.Destroy()
	tasks := `{task: "echo test"}`
	utils.On("readFile", ".darius.yml").Return("tasks: "+tasks, nil)
	utils.On("out", "\x1b[1;32m$ echo test\x1b[0m", true)
	utils.On("out", "  > test", true)
	utils.On("out", "\x1b[1;37;42mtask completed\x1b[0m", true)
	err := call(state, []string{"task"})
	assert.NoError(test, err)
	utils.AssertExpectations(test)
}

func TestRunTaskSendsStderr(test *testing.T) {
	state, utils := newTestState(false)
	defer state.Destroy()
	tasks := `{task: "echo test 1>&2"}`
	utils.On("readFile", ".darius.yml").Return("tasks: "+tasks, nil)
	utils.On("out", "\x1b[1;32m$ echo test 1>&2\x1b[0m", true)
	utils.On("out", "  \x1b[31m! test\x1b[0m", true)
	utils.On("out", "\x1b[1;37;42mtask completed\x1b[0m", true)
	err := call(state, []string{"task"})
	assert.NoError(test, err)
	utils.AssertExpectations(test)
}

func TestRunReportsName(test *testing.T) {
	state, utils := newTestState(false)
	defer state.Destroy()
	tasks := `{task: {name: "NAME", command: "echo test"}}`
	utils.On("readFile", ".darius.yml").Return("tasks: "+tasks, nil)
	utils.On("out", "\x1b[1;33m# NAME\x1b[0m", true)
	utils.On("out", "  \x1b[1;32m$ echo test\x1b[0m", true)
	utils.On("out", "    > test", true)
	utils.On("out", "\x1b[1;37;42mtask completed\x1b[0m", true)
	err := call(state, []string{"task"})
	assert.NoError(test, err)
	utils.AssertExpectations(test)
}

func TestRunRunsTwoTasks(test *testing.T) {
	state, utils := newTestState(false)
	defer state.Destroy()
	tasks := `{task: ["echo test1", "echo test2"]}`
	utils.On("readFile", ".darius.yml").Return("tasks: "+tasks, nil)
	utils.On("out", "\x1b[1;32m$ echo test1\x1b[0m", true)
	utils.On("out", "  > test1", true)
	utils.On("out", "\x1b[1;32m$ echo test2\x1b[0m", true)
	utils.On("out", "  > test2", true)
	utils.On("out", "\x1b[1;37;42mtask completed\x1b[0m", true)
	err := call(state, []string{"task"})
	assert.NoError(test, err)
	utils.AssertExpectations(test)
}

func TestRunRunsTwoTasksWithName(test *testing.T) {
	state, utils := newTestState(false)
	defer state.Destroy()
	tasks := `{task: {name: "NAME", command: ["echo test1", "echo test2"]}}`
	utils.On("readFile", ".darius.yml").Return("tasks: "+tasks, nil)
	utils.On("out", "\x1b[1;33m# NAME\x1b[0m", true)
	utils.On("out", "  \x1b[1;32m$ echo test1\x1b[0m", true)
	utils.On("out", "    > test1", true)
	utils.On("out", "  \x1b[1;32m$ echo test2\x1b[0m", true)
	utils.On("out", "    > test2", true)
	utils.On("out", "\x1b[1;37;42mtask completed\x1b[0m", true)
	err := call(state, []string{"task"})
	assert.NoError(test, err)
	utils.AssertExpectations(test)
}

func TestRunChecksContext(test *testing.T) {
	state, utils := newTestState(false)
	defer state.Destroy()
	tasks := `{task: {command: "echo test1", context: "/bin/false"}}`
	utils.On("readFile", ".darius.yml").Return("tasks: "+tasks, nil)
	utils.On("out", "\x1b[36m? /bin/false\x1b[0m", true)
	utils.On("out", "\x1b[1;37;42mtask completed\x1b[0m", true)
	err := call(state, []string{"task"})
	assert.NoError(test, err)
	utils.AssertExpectations(test)
}

func TestRunRunsRescue(test *testing.T) {
	state, utils := newTestState(false)
	defer state.Destroy()
	tasks := `{task: {command: "/bin/false", rescue: "echo RESCUE"}}`
	utils.On("readFile", ".darius.yml").Return("tasks: "+tasks, nil)
	utils.On("out", "\x1b[1;32m$ /bin/false\x1b[0m", true)
	utils.On("out", "\x1b[1;37;41m ** command execution failed: non-zero "+
		"exit status 1 received\x1b[0m", true)
	utils.On("out", "\x1b[1;33m[rescue]\x1b[0m", true)
	utils.On("out", "\x1b[1;32m$ echo RESCUE\x1b[0m", true)
	utils.On("out", "  > RESCUE", true)
	utils.On("out", "\x1b[1;37;42mtask completed\x1b[0m", true)
	err := call(state, []string{"task"})
	assert.NoError(test, err)
	utils.AssertExpectations(test)
}

func TestRunRunsEnsure(test *testing.T) {
	state, utils := newTestState(false)
	defer state.Destroy()
	tasks := `{task: {command: "/bin/false", ensure: "echo ENSURE"}}`
	utils.On("readFile", ".darius.yml").Return("tasks: "+tasks, nil)
	utils.On("out", "\x1b[1;32m$ /bin/false\x1b[0m", true)
	utils.On("out", "\x1b[1;37;41m ** command execution failed: non-zero "+
		"exit status 1 received\x1b[0m", true)
	utils.On("out", "\x1b[1;33m[ensure]\x1b[0m", true)
	utils.On("out", "\x1b[1;32m$ echo ENSURE\x1b[0m", true)
	utils.On("out", "  > ENSURE", true)
	utils.On("out", "\x1b[1;37;41m ** task execution failed (check logs for "+
		"details) ** \x1b[0m", true)
	err := call(state, []string{"task"})
	assert.Error(test, err)
	utils.AssertExpectations(test)
}

func TestRunRunsOnHost(test *testing.T) {
	state, utils := newTestState(false)
	defer state.Destroy()
	tasks := `{task: {host: "ssh.darius.local", command: "cat /etc/hostname"}}`
	utils.On("readFile", ".darius.yml").Return("tasks: "+tasks, nil)
	utils.On("out", "\x1b[35m% connecting to ssh.darius.local...\x1b[0m",
		true)
	utils.On("out", "\x1b[35m% connection established\x1b[0m", true)
	utils.On("out", "\x1b[1;32m$ cat /etc/hostname\x1b[0m", true)
	utils.On("out", "  > ssh.darius.local", true)
	utils.On("out", "\x1b[1;37;42mtask completed\x1b[0m", true)
	err := call(state, []string{"task"})
	assert.NoError(test, err)
	utils.AssertExpectations(test)
}

func TestRunRunsOnHostWithKeyFile(test *testing.T) {
	state, utils := newTestState(false)
	defer state.Destroy()
	tasks := `{task: {host: {host: "user@ssh.darius.local",
		key: "/root/.ssh/user"}, command: "whoami"}}`
	utils.On("readFile", ".darius.yml").Return("tasks: "+tasks, nil)
	utils.On("out", "\x1b[35m% connecting to user@ssh.darius.local"+
		"...\x1b[0m", true)
	utils.On("out", "\x1b[35m% connection established\x1b[0m", true)
	utils.On("out", "\x1b[1;32m$ whoami\x1b[0m", true)
	utils.On("out", "  > user", true)
	utils.On("out", "\x1b[1;37;42mtask completed\x1b[0m", true)
	err := call(state, []string{"task"})
	assert.NoError(test, err)
	utils.AssertExpectations(test)
}

func TestRunIncludesConfiguration(test *testing.T) {
	state, utils := newTestState(false)
	defer state.Destroy()
	utils.On("readFile", "FILE").Return(`{task: "echo test"}`, nil)
	utils.On("readFile", ".darius.yml").Return(`tasks: "${include
		FILE}"`, nil)
	utils.On("out", "\x1b[1;32m$ echo test\x1b[0m", true)
	utils.On("out", "  > test", true)
	utils.On("out", "\x1b[1;37;42mtask completed\x1b[0m", true)
	err := call(state, []string{"task"})
	assert.NoError(test, err)
	utils.AssertExpectations(test)
}

func TestRunIncludesConfigurationGlob(test *testing.T) {
	state, utils := newTestState(false)
	defer state.Destroy()
	utils.On("readFile", "FILE").Return(`{task: "echo test"}`, nil)
	utils.On("readFile", ".darius.yml").Return(`tasks: "${include P/*}"`, nil)
	utils.On("glob", "P/*").Return([]string{"FILE"}, nil)
	utils.On("out", "\x1b[1;32m$ echo test\x1b[0m", true)
	utils.On("out", "  > test", true)
	utils.On("out", "\x1b[1;37;42mtask completed\x1b[0m", true)
	err := call(state, []string{"task"})
	assert.NoError(test, err)
	utils.AssertExpectations(test)
}

func TestRunRunsNestedWithNameAndContext(test *testing.T) {
	state, utils := newTestState(false)
	defer state.Destroy()
	utils.On("readFile", ".darius.yml").Return(`
tasks:
  task:
    name: "NAME"
    command:
      - context: '! /bin/false'
        command: [echo 1, echo 2]
`, nil)
	utils.On("out", "\x1b[1;33m# NAME\x1b[0m", true)
	utils.On("out", "  \x1b[36m? ! /bin/false\x1b[0m", true)
	utils.On("out", "  \x1b[1;32m$ echo 1\x1b[0m", true)
	utils.On("out", "    > 1", true)
	utils.On("out", "  \x1b[1;32m$ echo 2\x1b[0m", true)
	utils.On("out", "    > 2", true)
	utils.On("out", "\x1b[1;37;42mtask completed\x1b[0m", true)
	err := call(state, []string{"task"})
	assert.NoError(test, err)
	utils.AssertExpectations(test)
}

func TestRunExpandsArgument(test *testing.T) {
	state, utils := newTestState(false)
	defer state.Destroy()
	state.argv = []string{"--arg", "VALUE"}
	utils.On("readFile", ".darius.yml").Return(`
tasks:
  task:
    args: {arg: {type: string}}
    command: echo ${args.arg}
`, nil)
	utils.On("out", "\x1b[1;32m$ echo VALUE\x1b[0m", true)
	utils.On("out", "  > VALUE", true)
	utils.On("out", "\x1b[1;37;42mtask completed\x1b[0m", true)
	err := call(state, []string{"task"})
	assert.NoError(test, err)
	utils.AssertExpectations(test)
}

func TestRunExpandsArgumentByParam(test *testing.T) {
	state, utils := newTestState(false)
	defer state.Destroy()
	state.argv = []string{}
	utils.On("readFile", ".darius.yml").Return(`
tasks:
  task:
    params: {arg: VALUE}
    command:
      args: {arg: {type: string}}
      command: echo ${args.arg}
`, nil)
	utils.On("out", "\x1b[1;32m$ echo VALUE\x1b[0m", true)
	utils.On("out", "  > VALUE", true)
	utils.On("out", "\x1b[1;37;42mtask completed\x1b[0m", true)
	err := call(state, []string{"task"})
	assert.NoError(test, err)
	utils.AssertExpectations(test)
}

func TestRunRunsUserTask(test *testing.T) {
	state, utils := newTestState(false)
	defer state.Destroy()
	state.argv = []string{}
	utils.On("readFile", ".darius.yml").Return(`
tasks:
  task:
    task: run-user-task
    task-name: echo
  echo: echo VALUE
`, nil)
	utils.On("out", "\x1b[35m% \"run-user-task\" is deprecated\x1b[0m", true)
	utils.On("out", "\x1b[1;32m$ echo VALUE\x1b[0m", true)
	utils.On("out", "  > VALUE", true)
	utils.On("out", "\x1b[1;37;42mtask completed\x1b[0m", true)
	err := call(state, []string{"task"})
	assert.NoError(test, err)
	utils.AssertExpectations(test)
}
