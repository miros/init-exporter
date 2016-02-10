package systemd

import (
  "testing"
  "github.com/stretchr/testify/assert"
  "github.com/spf13/afero"
  "strings"
)

import "github.com/davecgh/go-spew/spew"
var _ = spew.Dump

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

var defaultService = Service{
  Name: "some-service",
  Cmd: "run-some-service",
  Options: ServiceOptions{
    WorkingDirectory: "/projects/some-service",
    User: "run_user",
    Group: "run_group",
    KillTimeout: 12345,
    Respawn: Respawn{Interval: 10, Count: 100},
    Env: map[string]string{"env_var": "env_val", "env_var2": "env_val2"},
  },
}

const appUnitFilePath = "/units/some-app.service"
const helperFilePath = "/helpers/some-app_some-service.sh"
const serviceUnitFilePath = "/units/some-app_some-service.service"

func TestInstall(t *testing.T) {
  env := newTestEnv()
  sys := env.newSystemd(systemdConfig)

  sys.Install("some-app", []Service{defaultService})

  assert.True(t, env.fileExists(appUnitFilePath), "no app unit file")

  assert.True(t, env.fileExists(helperFilePath), "no helper file")
  helperFileData := env.readFile(helperFilePath)
  assert.Contains(t, helperFileData, "run-some-service")

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
  assert.Contains(t, unitFileData, "env_var=env_val")
  assert.Contains(t, unitFileData, "env_var2=env_val2")
  assert.Contains(t, unitFileData, ">> /var/log/some-app/some-service.log")

  assert.Contains(t, env.executedCommands, "systemctl enable some-app.service")
}

func TestInstallMultiCount(t *testing.T) {
  env := newTestEnv()
  sys := env.newSystemd(systemdConfig)

  multiService := defaultService
  multiService.Options.Count = 2

  sys.Install("some-app", []Service{multiService})

  assert.True(t, env.fileExists("/units/some-app_some-service1.service"), "no service unit file")
  assert.True(t, env.fileExists("/units/some-app_some-service2.service"), "no service unit file")
}

func TestUnInstall(t *testing.T) {
  env := newTestEnv()
  sys := env.newSystemd(systemdConfig)

  sys.Install("some-app", []Service{defaultService})

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
