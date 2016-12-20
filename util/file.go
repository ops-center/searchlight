package util

import (
	"encoding/json"
	"io/ioutil"
)

func ReadFile(path string) (string, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func ReadFileAs(path string, obj interface{}) error {
	d, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(d, obj); err != nil {
		return err
	}
	return nil
}
