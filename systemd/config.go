package systemd

import (
  "path"
)

type Config struct {
  HelperDir string
  TargetDir string
  User string
  Group string
  WorkingDirectory string
}

func (config *Config) unitPath(name string) string {
  return path.Join(config.TargetDir, name + ".service")
}

func (config *Config) helperPath(name string) string {
  return path.Join(config.HelperDir, name + ".sh")
}