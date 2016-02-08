package main

import (
  "io/ioutil"
  "os"
  "gopkg.in/yaml.v2"
)

type GlobalConfig struct {
    RunUser     string `yaml:"run_user"`
    RunGroup    string `yaml:"run_group"`
    HelperDir   string `yaml:"helper_dir"`
    TargetDir   string `yaml:"target_dir"`
    Prefix      string
}

func defaultConfig() GlobalConfig {
  return GlobalConfig{
    RunUser: "service",
    RunGroup: "service",
    HelperDir: "/users/miros/systemd/var/local/upstart_helpers/",
    TargetDir: "/users/miros/systemd/etc/systemd/system/",
    Prefix: "fb-",
  }
}

func ReadGlobalConfig(path string) GlobalConfig {
  config := defaultConfig()
  data, err := ioutil.ReadFile(path)

  if (err != nil && !os.IsNotExist(err)) {
    panic(err)
  }

  if err := yaml.Unmarshal(data, &config); err != nil {
    panic(err)
  }

  return config
}



