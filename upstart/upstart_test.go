package upstart

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

const appUnitFilePath = "/units/some-app.conf"
const helperFilePath = "/helpers/some-app_some-service.sh"
const serviceUnitFilePath = "/units/some-app_some-service.conf"

func newExporter(env *test_env.TestEnv, config exporter.Config) *exporter.Exporter {
	provider := New()
	return env.NewExporter(config, provider)
}

func TestInstall(t *testing.T) {
	env := test_env.New()
	sys := newExporter(env, systemdConfig)

	sys.Install("some-app", procfile.App{Services: []procfile.Service{defaultService}})

	assert.True(t, env.FileExists(appUnitFilePath), "no app unit file")

	assert.True(t, env.FileExists(helperFilePath), "no helper file")
	helperFileData := env.ReadFile(helperFilePath)
	assert.Contains(t, helperFileData, "run-some-service")
	assert.Contains(t, helperFileData, "cd /projects/some-service")
	assert.Contains(t, helperFileData, "env_var=env_val")
	assert.Contains(t, helperFileData, "env_var2=env_val2")

	assert.True(t, env.FileExists(serviceUnitFilePath), "no service unit file")
	unitFileData := env.ReadFile(serviceUnitFilePath)

	assert.Contains(t, unitFileData, "kill timeout 12345")
	assert.Contains(t, unitFileData, "respawn limit 100 10")
	assert.Contains(t, unitFileData, "sudo -u run_user")
	assert.Contains(t, unitFileData, ">> /var/log/some-app/some-service.log")
}
