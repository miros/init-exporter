package test_env

import (
	"github.com/miros/init-exporter/exporter"
	"github.com/miros/init-exporter/utils"
	"github.com/spf13/afero"
	"strings"
)

type TestEnv struct {
	ExecutedCommands []string
	fs               afero.Fs
}

func New() *TestEnv {
	env := new(TestEnv)
	env.fs = afero.NewMemMapFs()
	return env
}

func (env *TestEnv) FakeExecSystemCommand() utils.SystemExecutor {
	return func(name string, args ...string) error {
		env.ExecutedCommands = append(env.ExecutedCommands, name+" "+strings.Join(args, " "))
		return nil
	}
}

func (env *TestEnv) ReadFile(path string) string {
	data, err := afero.ReadFile(env.fs, path)

	if err != nil {
		panic(err)
	}

	return string(data)
}

func (env *TestEnv) FileExists(path string) bool {
	result, _ := afero.Exists(env.fs, path)
	return result
}

func (env *TestEnv) WriteFile(path string, data string) {
	utils.MustWriteFile(env.fs, path, data)
}

func (env *TestEnv) NewExporter(config exporter.Config, provider exporter.Provider) *exporter.Exporter {
	sys := exporter.New(config, provider)
	sys.Fs = env.fs

	return sys
}
