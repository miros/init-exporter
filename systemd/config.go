package systemd

import (
  "path"
  "github.com/miros/systemd-exporter/systemd/validation"
)

type Config struct {
  HelperDir string
  TargetDir string
  User string
  Group string
  DefaultWorkingDirectory string
}

func (config *Config) unitPath(name string) string {
  return path.Join(config.TargetDir, name + ".service")
}

func (config *Config) helperPath(name string) string {
  return path.Join(config.HelperDir, name + ".sh")
}

func (config *Config) Validate() error {
  if err := validation.NoSpecialSymbols(config.User); err != nil {
    return err
  }

  if err := validation.NoSpecialSymbols(config.Group); err != nil {
    return err
  }

  if err := validation.Path(config.DefaultWorkingDirectory); err != nil {
    return err
  }

  return nil
}