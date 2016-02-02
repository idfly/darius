package main

import (
	"github.com/idfly/darius"

	"github.com/stretchr/testify/mock"
)

type utilsMock struct {
	mock.Mock
}

func (mock *utilsMock) out(message string, newLine bool) {
	mock.Called(message, newLine)
}

func (mock *utilsMock) err(message string, newLine bool) {
	mock.Called(message, newLine)
}

func (mock *utilsMock) readFile(file string) ([]byte, error) {
	args := mock.Called(file)
	return []byte(args.String(0)), args.Error(1)
}

func (mock *utilsMock) glob(pattern string) ([]string, error) {
	args := mock.Called(pattern)
	return args.Get(0).([]string), args.Error(1)
}

func (mock *utilsMock) call(
	state darius.State,
	task map[interface{}]interface{},
) error {
	args := mock.Called(task)
	return args.Error(0)
}
