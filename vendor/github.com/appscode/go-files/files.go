package files

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
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
	err = json.Unmarshal(d, obj)
	if err != nil {
		return err
	}
	return nil
}

func WriteFile(path string, obj interface{}) error {
	d, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return err
	}
	EnsureDirectory(path)
	err = ioutil.WriteFile(path, d, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func WriteString(path string, data string) bool {
	EnsureDirectory(path)
	err := ioutil.WriteFile(path, []byte(data), os.ModePerm)

	if err != nil {
		return false
	}
	return true
}

func AppendToFile(path string, values string) error {
	EnsureDirectory(path)
	if _, err := os.Stat(path); err != nil {
		ioutil.WriteFile(path, []byte(""), os.ModePerm)
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	f.WriteString("\n")
	_, err = f.WriteString(values)
	if err != nil {
		return err
	}
	return nil
}

func EnsureDirectory(path string) {
	parent := filepath.Dir(path)
	if _, err := os.Stat(parent); err != nil {
		err = os.MkdirAll(parent, os.ModePerm)
	}
}

func IsFileExists(path string) bool {
	EnsureDirectory(path)
	if _, err := os.Stat(path); err != nil {
		return false
	}
	return true
}
