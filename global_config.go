package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type GlobalConfig struct {
	RunUser          string `yaml:"run_user"`
	RunGroup         string `yaml:"run_group"`
	WorkingDirectory string `yaml:"working_directory"`
	HelperDir        string `yaml:"helper_dir"`
	TargetDir        string `yaml:"target_dir"`
	Prefix           string
}

func defaultConfig() GlobalConfig {
	return GlobalConfig{
		RunUser:          "service",
		RunGroup:         "service",
		WorkingDirectory: "/tmp",
		HelperDir:        "/var/local/upstart_helpers/",
		Prefix:           "fb-",
	}
}

func ReadGlobalConfig(path string) GlobalConfig {
	config := defaultConfig()
	data, err := ioutil.ReadFile(path)

	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		panic(err)
	}

	return config
}
