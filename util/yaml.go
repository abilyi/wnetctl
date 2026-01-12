package util

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

func ReadObject(path string, object interface{}) error {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return err
	}
	if fileInfo.IsDir() {
		return fmt.Errorf("%s is a directory", path)
	}
	in, err := os.OpenFile(path, os.O_RDONLY, 0600)
	if err != nil {
		return err
	} else {
		defer in.Close()
	}
	content := make([]byte, fileInfo.Size())
	_, err = in.Read(content)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(content, object)
}

func WriteObject(path string, object interface{}) error {
	content, err := yaml.Marshal(object)
	if err != nil {
		return err
	}
	out, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = out.Write(content)
	return err
}
