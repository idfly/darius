package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatTwoLinesReturnsIndentedMessage(test *testing.T) {
	colorize := func(values ...interface{}) string {
		return values[0].(string)
	}

	result := format(0, "  ", "$ ", "line 1\nline 2", colorize)
	assert.Equal(test, "  $ line 1\n    line 2", result)
}
