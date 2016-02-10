package systemd

import (
  "testing"
  "github.com/stretchr/testify/assert"
  "github.com/spf13/afero"
  "strings"
)

type testEnv struct {
  executedCommands []string
  fs afero.Fs
}

func newTestEnv() *testEnv {
  env := new(testEnv)
  env.fs = afero.NewMemMapFs()
  return env
}

func (env *testEnv) fakeExecSystemCommand() systemExecutor {
  return func(name string, args ...string) error {
    env.executedCommands = append(env.executedCommands,  name + " " + strings.Join(args, " "))
    return nil
  }
}

func (env *testEnv) readFile(path string) string {
  data, err := afero.ReadFile(env.fs, path)

  if (err != nil) {
    panic(err)
  }

  return string(data)
}

func (env *testEnv) fileExists(path string) bool {
  result, _ := afero.Exists(env.fs, path)
  return result
}

func (env *testEnv) writeFile(path string, data string) {
  writeFile(env.fs, path, data)
}

func (env *testEnv) newSystemd(config Config) *Systemd {
  sys := New(config)
  sys.fs = env.fs
  sys.execSystemCommand = env.fakeExecSystemCommand()
  return sys
}

var systemdConfig Config = Config{
  HelperDir: "/helpers",
  TargetDir: "/units",
  User: "user",
  Group: "group",
  DefaultWorkingDirectory: "/tmp",
}

var services []Service = []Service{
  Service{
    Name: "some-service",
    Cmd: "run",
    Options: ServiceOptions{
      WorkingDirectory: "/projects/some-service",
      User: "run_user",
      Group: "run_group",
      KillTimeout: 12345,
      Respawn: Respawn{Interval: 10, Count: 100},
    },
  },
}

const appUnitFilePath = "/units/some-app.service"
const helperFilePath = "/helpers/some-app_some-service.sh"
const serviceUnitFilePath = "/units/some-app_some-service.service"

func TestInstall(t *testing.T) {
  env := newTestEnv()
  sys := env.newSystemd(systemdConfig)

  sys.Install("some-app", services)

  assert.True(t, env.fileExists(appUnitFilePath), "no app unit file")

  assert.True(t, env.fileExists(helperFilePath), "no helper file")
  helperFileData := env.readFile(helperFilePath)
  assert.Contains(t, helperFileData, "run")

  assert.True(t, env.fileExists(serviceUnitFilePath), "no service unit file")
  unitFileData := env.readFile(serviceUnitFilePath)
  assert.Contains(t, unitFileData, "PartOf=some-app.service")
  assert.Contains(t, unitFileData, "TimeoutStopSec=12345")
  assert.Contains(t, unitFileData, "StartLimitInterval=10")
  assert.Contains(t, unitFileData, "StartLimitBurst=100")
  assert.Contains(t, unitFileData, "WorkingDirectory=/projects/some-service")
  assert.Contains(t, unitFileData, "User=run_user")
  assert.Contains(t, unitFileData, "Group=run_group")
  assert.Contains(t, unitFileData, "WorkingDirectory=/projects/some-service")

  assert.Contains(t, env.executedCommands, "systemctl enable some-app.service")
}

func TestUnInstall(t *testing.T) {
  env := newTestEnv()
  sys := env.newSystemd(systemdConfig)

  sys.Install("some-app", services)

  env.writeFile("/helpers/file_to_keep.sh", "data")
  env.writeFile("/units/file_to_keep.service", "data")

  sys.Uninstall("some-app")

  assert.False(t, env.fileExists(appUnitFilePath), "app unit file exists")
  assert.False(t, env.fileExists(helperFilePath), "helper file exists")
  assert.False(t, env.fileExists(serviceUnitFilePath), "service unit file exists")

  assert.True(t, env.fileExists("/helpers/file_to_keep.sh"), "wrong file deleted")
  assert.True(t, env.fileExists("/units/file_to_keep.service"), "wrong file deleted")

  assert.Contains(t, env.executedCommands, "systemctl disable some-app.service")
}
