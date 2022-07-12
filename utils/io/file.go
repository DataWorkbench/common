package io

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

func ReadYaml(file string) (map[string]interface{}, error) {
	var err error
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	m := map[string]interface{}{}
	if err = yaml.Unmarshal(bytes, &m); err != nil {
		return nil, err
	}
	return m, nil
}
