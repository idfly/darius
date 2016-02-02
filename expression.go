package darius

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type Expression struct {
	State            State
	ExpandExpression func(string, string) (interface{}, error)

	stack []string
}

var (
	expandRegexp = regexp.MustCompile(`\$+\{.*?\}`)
)

func NewExpression(
	state State,
	expand func(string, string) (interface{}, error),

) *Expression {
	return &Expression{state, expand, []string{}}
}

func (expander *Expression) Expand(
	value interface{},
	recursive bool,
) (interface{}, error) {
	str, ok := value.(string)
	if !ok {
		return value, nil
	}

	expr := expandRegexp.FindString(str)
	if expr == str {
		result, err := expander.expand(expr)
		if err != nil {
			return nil, err
		}

		if recursive {
			return expander.expandRecursive(result)
		}

		return result, nil
	}

	var errs []string
	result := expandRegexp.ReplaceAllStringFunc(str, func(expr string) string {
		result, err := expander.expand(expr)
		if err != nil {
			errs = append(errs, err.Error())
			return expr
		}

		return fmt.Sprint(result)
	})

	if len(errs) > 0 {
		return nil, errors.New(strings.Join(errs, "\n"))
	}

	return result, nil
}

func (expander *Expression) expandRecursive(
	value interface{},
) (interface{}, error) {
	mapping, ok := value.(map[interface{}]interface{})
	if ok {
		result := map[interface{}]interface{}{}
		for key, value := range mapping {
			var err error
			result[key], err = expander.Expand(value, true)
			if err != nil {
				return nil, err
			}
		}

		return result, nil
	}

	array, ok := value.([]interface{})
	if ok {
		result := make([]interface{}, len(array))
		for index, value := range array {
			var err error
			result[index], err = expander.Expand(value, true)
			if err != nil {
				return nil, err
			}
		}

		return result, nil
	}

	return value, nil
}

func (expander *Expression) expand(expr string) (interface{}, error) {
	prefix := regexp.MustCompile(`^\$+`).FindString(expr)
	truncatedPrefix := strings.Repeat("$", len(prefix)/2)
	if len(prefix)%2 == 0 {
		return truncatedPrefix + expr[len(prefix):], nil
	}

	for _, value := range expander.stack {
		if value == expr {
			return nil, errors.New("recursive expression detected: " +
				strings.Join(append(expander.stack, expr), " -> "))
		}
	}

	expander.stack = append(expander.stack, expr)

	expr = strings.TrimLeft(expr, "${")
	expr = expr[:len(expr)-1]

	parts := regexp.MustCompile(`^(\w+)(?:\W(.*?)$|$)`).FindStringSubmatch(expr)
	if len(parts) == 0 {
		return expr, errors.New("wrong expression: " + expr)
	}

	var err error
	var result interface{}
	if parts[1] == "vars" {
		result, err = expander.expandVariable(strings.Split(parts[2], "."))
	} else if parts[1] == "args" {
		result, err = expander.expandArguments(strings.Split(parts[2], "."))
	} else {
		result, err = expander.ExpandExpression(parts[1], parts[2])
	}

	if err != nil {
		return nil, err
	}

	if truncatedPrefix != "" {
		return fmt.Sprint(truncatedPrefix, result), nil
	}

	return result, nil

}

func (expander *Expression) expandVariable(expr []string) (interface{}, error) {
	current := expander.State
	for current != nil {
		task := current.Task()
		if task != nil {
			taskVars, ok := task["vars"]
			if ok {
				result, found, err := ExpandMap(current, taskVars, expr)

				if found || err != nil {
					return result, err
				}
			}
		}

		var ok bool
		current, ok = current.Parent()
		if !ok {
			current = nil
		}
	}

	config, ok := expander.State.Config()["vars"]
	if ok {
		result, found, err := ExpandMap(expander.State, config, expr)
		if found || err != nil {
			return result, err
		}
	}

	return nil, errors.New("undefined variable vars." + strings.Join(expr, "."))
}

func (expander *Expression) expandArguments(expr []string) (interface{}, error) {
	result, found, err := ExpandMap(expander.State, expander.State.Args(), expr)
	if err != nil {
		return nil, err
	}

	if !found {
		return nil, errors.New("undefined variable args." +
			strings.Join(expr, "."))
	}

	return result, nil
}

func ExpandMap(
	state State,
	value interface{},
	expr []string,
) (interface{}, bool, error) {
	expanded, err := state.Expand(value, false)
	if err != nil {
		return nil, false, err
	}

	mapping, ok := expanded.(map[interface{}]interface{})
	if !ok {
		return nil, false, errors.New("expanded element should be map")
	}

	var current interface{} = mapping
	for _, variable := range expr {
		mapping, ok := current.(map[interface{}]interface{})
		if !ok {
			return nil, false, errors.New("variable should be map")
		}

		currentRaw, ok := mapping[variable]
		if !ok {
			return nil, false, nil
		}

		current, err = state.Expand(currentRaw, false)
		if err != nil {
			return nil, false, err
		}
	}

	return current, true, nil
}
