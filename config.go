package darius

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	includeStack []string
	ReadFile     func(string) ([]byte, error)
	Glob         func(string) ([]string, error)
}

func (config Config) Load(file string) (map[interface{}]interface{}, error) {
	plain, err := config.load(file)
	if err != nil {
		return nil, err
	}

	result, ok := plain.(map[interface{}]interface{})
	if !ok {
		return nil, errors.New("config value should be map in " + file)
	}

	return result, nil
}

func (config Config) load(file string) (interface{}, error) {
	for _, fileInProcess := range config.includeStack {
		if fileInProcess == file {
			return nil, errors.New("recursive include detected: " +
				strings.Join(config.includeStack, " <- "))
		}
	}

	config.includeStack = append(config.includeStack, file)
	defer func() {
		config.includeStack = config.includeStack[:len(config.includeStack)-1]
	}()

	contents, err := config.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var raw interface{} = nil
	slice := yaml.MapSlice{}
	err = yaml.Unmarshal(contents, &slice)
	rootIncluded := false

	if err == nil {
		raw = slice
	} else {
		trimmed := strings.Trim(string(contents), " \n")
		if strings.HasPrefix(trimmed, "-") {
			contents = []byte("__root: \n" + string(contents))
		} else {
			contents = []byte("__root: " + string(contents))
		}

		err = yaml.Unmarshal(contents, &slice)
		if err != nil {
			return nil, err
		}

		rootIncluded = true
		raw = slice
	}

	parsed, err := config.parseYaml(raw)
	if err != nil {
		return nil, err
	}

	if rootIncluded {
		parsed = parsed.(map[interface{}]interface{})["__root"]
	}

	result, err := config.include(parsed)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (config Config) parseYaml(raw interface{}) (interface{}, error) {
	slice, ok := raw.(yaml.MapSlice)
	if ok {
		result := map[interface{}]interface{}{}
		for _, element := range slice {
			parsed, err := config.parseYaml(element.Value)
			if err != nil {
				return nil, err
			}

			result[element.Key] = parsed
		}

		return result, nil
	}

	array, ok := raw.([]interface{})
	if ok {
		for index, element := range array {
			var err error
			array[index], err = config.parseYaml(element)
			if err != nil {
				return nil, err
			}
		}

		return array, nil
	}

	return raw, nil
}

var (
	includeRegexp    = regexp.MustCompile(`^\$\{include (.*?)\}$`)
	includeRegexpOld = regexp.MustCompile(`^\$include (.*?)$`)
)

func (config Config) include(
	value interface{},
) (interface{}, error) {
	var err error

	str, ok := value.(string)
	if ok {
		match := includeRegexp.FindStringSubmatch(str)
		if len(match) == 0 {
			match = includeRegexpOld.FindStringSubmatch(str)
		}

		if len(match) > 0 {
			file := match[1]

			if !strings.HasPrefix(file, "/") {
				current := filepath.Dir(file)
				if len(config.includeStack) > 0 {
					current = config.includeStack[len(config.includeStack)-1]
				}

				file = filepath.Join(filepath.Dir(current), file)
			}

			if !strings.Contains(file, "*") {
				return config.load(file)
			}

			files, err := config.Glob(file)
			if err != nil {
				return nil, err
			}

			return config.loadFiles(files)
		}

		return value, nil
	}

	mapping, ok := value.(map[interface{}]interface{})
	if ok {
		for key, value := range mapping {
			mapping[key], err = config.include(value)
			if err != nil {
				return nil, err
			}
		}

		return mapping, nil
	}

	array, ok := value.([]interface{})
	if ok {
		for index, value := range array {
			array[index], err = config.include(value)
			if err != nil {
				return nil, err
			}
		}

		return array, nil
	}

	return value, nil
}

func (config *Config) loadFiles(files []string) (interface{}, error) {
	result := map[interface{}]interface{}{}
	for _, file := range files {
		current, err := config.load(file)
		if err != nil {
			return nil, err
		}

		mapping, ok := current.(map[interface{}]interface{})
		if !ok {
			return nil, errors.New("config in " + file + " should be map")
		}

		for key, value := range mapping {
			_, ok := result[key]
			if ok {
				return nil, errors.New("value is alredy present in " +
					fmt.Sprint(key) + " (loading " + file + ")")
			}

			result[key] = value
		}
	}

	return result, nil
}
