package systemd

import (
  "os/exec"
  "path/filepath"
  "github.com/spf13/afero"
)

type fileGlobber func(string) ([]string, error)
type systemExecutor func(name string, arg ...string) error

type Systemd struct {
  Config Config
  fs afero.Fs
  globFilepath fileGlobber
  execSystemCommand systemExecutor
}

func New(config Config) *Systemd {
  return &Systemd{
    Config: config,
    fs: afero.NewOsFs(),
    globFilepath: globFilepath,
    execSystemCommand: execSystemCommand,
  }
}

func globFilepath(pattern string) ([]string, error) {
  return filepath.Glob(pattern)
}

func execSystemCommand(name string, arg ...string) error {
  return exec.Command(name, arg...).Run();
}
