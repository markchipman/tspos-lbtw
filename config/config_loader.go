package config

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

// Bind Read config file and load it into config structure
func Bind(path string, config interface{}) error {

	if _, err := os.Stat(path); err != nil {
		return errors.New("config path not valid")
	}

	ymlData, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal([]byte(ymlData), config)
	return err
}
