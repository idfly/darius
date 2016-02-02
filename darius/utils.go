package main

import (
	"io"
	"io/ioutil"
	"path/filepath"
)

type utilsInterface interface {
	out(string, bool)
	err(string, bool)
	readFile(string) ([]byte, error)
	glob(string) ([]string, error)
}

type utils struct {
	stdout io.Writer
	stderr io.Writer
}

func (utils utils) out(message string, newLine bool) {
	utils.stdout.Write([]byte(message))
	if newLine {
		utils.stdout.Write([]byte("\n"))
	}
}

func (utils utils) err(message string, newLine bool) {
	utils.stderr.Write([]byte(message))
	if newLine {
		utils.stdout.Write([]byte("\n"))
	}
}

func (utils utils) readFile(file string) ([]byte, error) {
	return ioutil.ReadFile(file)
}

func (utils utils) glob(file string) ([]string, error) {
	return filepath.Glob(file)
}
