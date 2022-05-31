package io

import (
"gopkg.in/yaml.v3"
"io/ioutil"
"os"
)

func FileExist(file string) (bool, error) {
	_, err := os.Stat(file)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func ReadYaml(file string) (map[string]interface{}, error) {
	var err error
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	m := map[string]interface{}{}
	if err = yaml.Unmarshal(bytes, &m); err != nil {
		return m, nil
	}
	return nil, err
}

