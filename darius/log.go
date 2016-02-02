package main

import (
	"os"
	"strings"

	"github.com/idfly/darius"

	"github.com/fatih/color"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	colorLogName        func(...interface{}) string
	colorLogSystem      func(...interface{}) string
	colorLogStdOut      func(...interface{}) string
	colorLogStdErr      func(...interface{}) string
	colorLogCommand     func(...interface{}) string
	colorLogCommandFail func(...interface{}) string
	colorLogContext     func(...interface{}) string
	colorLogRescue      func(...interface{}) string
	colorLogEnsure      func(...interface{}) string
	colorLogTaskFail    func(...interface{}) string
	colorLogTaskSuccess func(...interface{}) string
	formatters          map[darius.LogLevel]func(int, string) string
)

func init() {
	colorLogName = colorize(color.Bold, color.FgYellow)
	colorLogSystem = colorize(color.FgMagenta)
	colorLogStdOut = func(values ...interface{}) string {
		return values[0].(string)
	}

	colorLogStdErr = colorize(color.FgRed)
	colorLogCommand = colorize(color.Bold, color.FgGreen)
	colorLogCommandFail = colorize(color.Bold, color.FgWhite, color.BgRed)
	colorLogContext = colorize(color.FgCyan)
	colorLogRescue = colorize(color.Bold, color.FgYellow)
	colorLogEnsure = colorize(color.Bold, color.FgYellow)
	colorLogTaskFail = colorize(color.Bold, color.FgWhite, color.BgRed)
	colorLogTaskSuccess = colorize(color.Bold, color.FgWhite, color.BgGreen)

	formatters = map[darius.LogLevel]func(int, string) string{
		darius.LogName: func(level int, message string) string {
			return format(level, "", "# ", message, colorLogName)
		},

		darius.LogSystem: func(level int, message string) string {
			return format(level, "", "% ", message, colorLogSystem)
		},

		darius.LogCommand: func(level int, message string) string {
			return format(level, "", "$ ", message, colorLogCommand)
		},

		darius.LogStdOut: func(level int, message string) string {
			return format(level, "  ", "> ", message, colorLogStdOut)
		},

		darius.LogStdErr: func(level int, message string) string {
			return format(level, "  ", "! ", message, colorLogStdErr)
		},

		darius.LogCommandFail: func(level int, message string) string {
			return format(level, "", " ** ", message, colorLogCommandFail)
		},

		darius.LogContext: func(level int, message string) string {
			return format(level, "", "? ", message, colorLogContext)
		},

		darius.LogRescue: func(level int, message string) string {
			return format(level, "", "", message, colorLogRescue)
		},

		darius.LogEnsure: func(level int, message string) string {
			return format(level, "", "", message, colorLogEnsure)
		},

		darius.LogTaskFail: func(level int, message string) string {
			return format(level, "", "", message, colorLogTaskFail)
		},

		darius.LogTaskSuccess: func(level int, message string) string {
			return format(level, "", "", message, colorLogTaskSuccess)
		},
	}
}

func colorize(colors ...color.Attribute) func(...interface{}) string {
	color := color.New(colors...)
	color.EnableColor()
	return color.SprintFunc()
}

func format(
	level int,
	prefix string,
	coloredPrefix string,
	message string,
	colorize func(...interface{}) string,
) string {
	lines := []string{message}
	size, _, err := terminal.GetSize(int(os.Stdin.Fd()))
	if err != nil || size == -1 {
		size = 80
	}

	lines = wrap(message, size)

	indent := strings.Repeat(" ", len(coloredPrefix))
	prefix = strings.Repeat("  ", level) + prefix

	for index, line := range lines {
		lines[index] = prefix + colorize(coloredPrefix+line)
		coloredPrefix = indent
	}

	return strings.Join(lines, "\n")
}

func wrap(message string, width int) []string {
	lines := strings.Split(message, "\n")
	result := []string{}
	for _, line := range lines {
		for index := 0; index < len(line); index += width {
			next := index + width
			if next > len(line) {
				next = len(line)
			}

			result = append(result, line[index:next])
		}
	}

	return result
}
