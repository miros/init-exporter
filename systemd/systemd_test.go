package systemd

import (
	"github.com/miros/init-exporter/exporter"
	"github.com/miros/init-exporter/procfile"
	"github.com/miros/init-exporter/utils/test_env"
	"github.com/stretchr/testify/assert"
	"testing"
)

import "github.com/davecgh/go-spew/spew"

var _ = spew.Dump

var systemdConfig exporter.Config = exporter.Config{
	HelperDir: "/helpers",
	TargetDir: "/units",
	User:      "user",
	Group:     "group",
	DefaultWorkingDirectory: "/tmp",
}

var defaultService = procfile.Service{
	Name: "some-service",
	Cmd:  "run-some-service",
	Options: procfile.ServiceOptions{
		WorkingDirectory: "/projects/some-service",
		User:             "run_user",
		Group:            "run_group",
		KillTimeout:      12345,
		Respawn:          procfile.Respawn{Interval: 10, Count: 100},
		Env:              map[string]string{"env_var": "env_val", "env_var2": "env_val2"},
	},
}

const appUnitFilePath = "/units/some-app.service"
const helperFilePath = "/helpers/some-app_some-service.sh"
const serviceUnitFilePath = "/units/some-app_some-service.service"

func newExporter(env *test_env.TestEnv, config exporter.Config) *exporter.Exporter {
	provider := New()
	provider.execSystemCommand = env.FakeExecSystemCommand()

	return env.NewExporter(config, provider)
}

func TestInstall(t *testing.T) {
	env := test_env.New()
	sys := newExporter(env, systemdConfig)

	sys.Install("some-app", []procfile.Service{defaultService})

	assert.True(t, env.FileExists(appUnitFilePath), "no app unit file")

	assert.True(t, env.FileExists(helperFilePath), "no helper file")
	helperFileData := env.ReadFile(helperFilePath)
	assert.Contains(t, helperFileData, "run-some-service")

	assert.True(t, env.FileExists(serviceUnitFilePath), "no service unit file")
	unitFileData := env.ReadFile(serviceUnitFilePath)
	assert.Contains(t, unitFileData, "PartOf=some-app.service")
	assert.Contains(t, unitFileData, "TimeoutStopSec=12345")
	assert.Contains(t, unitFileData, "StartLimitInterval=10")
	assert.Contains(t, unitFileData, "StartLimitBurst=100")
	assert.Contains(t, unitFileData, "WorkingDirectory=/projects/some-service")
	assert.Contains(t, unitFileData, "User=run_user")
	assert.Contains(t, unitFileData, "Group=run_group")
	assert.Contains(t, unitFileData, "env_var=env_val")
	assert.Contains(t, unitFileData, "env_var2=env_val2")
	assert.Contains(t, unitFileData, ">> /var/log/some-app/some-service.log")

	assert.Contains(t, env.ExecutedCommands, "systemctl enable some-app.service")
}

func TestInstallMultiCount(t *testing.T) {
	env := test_env.New()
	sys := newExporter(env, systemdConfig)

	multiService := defaultService
	multiService.Options.Count = 2

	sys.Install("some-app", []procfile.Service{multiService})

	assert.True(t, env.FileExists("/units/some-app_some-service1.service"), "no service unit file")
	assert.True(t, env.FileExists("/units/some-app_some-service2.service"), "no service unit file")
}

func TestUnInstall(t *testing.T) {
	env := test_env.New()
	sys := newExporter(env, systemdConfig)

	sys.Install("some-app", []procfile.Service{defaultService})

	env.WriteFile("/helpers/file_to_keep.sh", "data")
	env.WriteFile("/units/file_to_keep.service", "data")

	sys.Uninstall("some-app")

	assert.False(t, env.FileExists(appUnitFilePath), "app unit file exists")
	assert.False(t, env.FileExists(helperFilePath), "helper file exists")
	assert.False(t, env.FileExists(serviceUnitFilePath), "service unit file exists")

	assert.True(t, env.FileExists("/helpers/file_to_keep.sh"), "wrong file deleted")
	assert.True(t, env.FileExists("/units/file_to_keep.service"), "wrong file deleted")

	assert.Contains(t, env.ExecutedCommands, "systemctl disable some-app.service")
}
