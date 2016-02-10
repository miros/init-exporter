package systemd

import (
  "os/exec"
  "github.com/spf13/afero"
)

type systemExecutor func(name string, arg ...string) error

type Systemd struct {
  Config Config
  fs afero.Fs
  execSystemCommand systemExecutor
}

func New(config Config) *Systemd {
  return &Systemd{
    Config: config,
    fs: afero.NewOsFs(),
    execSystemCommand: execSystemCommand,
  }
}

func execSystemCommand(name string, arg ...string) error {
  return exec.Command(name, arg...).Run();
}
