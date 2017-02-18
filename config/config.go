package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// ConfYaml is config structure.
type ConfYaml struct {
	Core SectionCore `yaml:"core"`
}

// SectionCore is sub section of config.
type SectionCore struct {
	Listener SectionListener `yaml:"listener"`
	Database SectionDatabase `yaml:"database"`
}

// SectionListener is sub section of SectionCore.
type SectionListener struct {
	Bind string `yaml:"bind"`
	Port string `yaml:"port"`
}

// SectionDatabase is sub section of SectionCore.
type SectionDatabase struct {
	Url string `yaml:"url"`
}

// BuildDefaultConf is default config setting.
func BuildDefaultConf() ConfYaml {
	var conf ConfYaml

	// Core Listener
	conf.Core.Listener.Bind = "localhost"
	conf.Core.Listener.Port = "2727"

	// Core Database
	conf.Core.Database.Url = "mongodb://localhost:27017/tspos_lbtw"

	return conf
}

// LoadConfYaml provide load yml config.
func LoadConfYaml(confPath string) (ConfYaml, error) {
	var conf ConfYaml

	confFile, err := ioutil.ReadFile(confPath)

	if err != nil {
		return conf, err
	}

	err = yaml.Unmarshal(confFile, &conf)

	if err != nil {
		return conf, err
	}

	return conf, nil
}
