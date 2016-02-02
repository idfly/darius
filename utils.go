package darius

func CreateTask(task interface{}) map[interface{}]interface{} {
	mapping, ok := task.(map[interface{}]interface{})
	if ok {
		return mapping
	}

	return map[interface{}]interface{}{"command": task}
}

func ExpandTask(
	state State,
	task map[interface{}]interface{},
) error {
	var err error
	var ok bool

	_, ok = task["name"]
	if ok {
		task["name"], err = state.Expand(task["name"], true)
		if err != nil {
			return err
		}
	}

	_, ok = task["host"]
	if ok {
		task["host"], err = state.Expand(task["host"], true)
		if err != nil {
			return err
		}
	}

	_, ok = task["context"]
	if ok {
		task["context"], err = state.Expand(task["context"], false)
		if err != nil {
			return err
		}
	}

	_, ok = task["job"]
	if ok {
		task["job"], err = state.Expand(task["job"], false)
		if err != nil {
			return err
		}
	}

	_, ok = task["command"]
	if ok {
		task["command"], err = state.Expand(task["command"], false)
		if err != nil {
			return err
		}
	}

	_, ok = task["rescue"]
	if ok {
		task["rescue"], err = state.Expand(task["rescue"], false)
		if err != nil {
			return err
		}
	}

	_, ok = task["ensure"]
	if ok {
		task["ensure"], err = state.Expand(task["ensure"], false)
		if err != nil {
			return err
		}
	}

	return nil
}

func Copy(value interface{}) interface{} {
	mapping, ok := value.(map[interface{}]interface{})
	if ok {
		result := map[interface{}]interface{}{}
		for key, value := range mapping {
			result[key] = Copy(value)
		}

		return result
	}

	array, ok := value.([]interface{})
	if ok {
		result := make([]interface{}, len(array))
		for index, value := range array {
			result[index] = Copy(value)
		}

		return result
	}

	return value
}
