package procfile

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProcfileV1(t *testing.T) {
	data := `
    cmd1: run-cmd1
    # comment
    cmd2: run-cmd2
  `
	app, _ := parseProcfile([]byte(data))

	assert.Equal(t, []Service{
		Service{Name: "cmd1", Cmd: "run-cmd1"},
		Service{Name: "cmd2", Cmd: "run-cmd2"},
	}, app.Services)
}

func TestProcfileV2(t *testing.T) {
	data := `
    version: 2
    env:
      env1: env1-val
    working_directory: /working-dir
    commands:
      cmd1:
        command: run-cmd1
        kill_timeout: 60
        respawn:
          count: 5
          interval: 10
        log: /path/to/log
        count: 2
      cmd2:
        command: run-cmd2
        working_directory: /working-dir2
        env:
          env1: env1-val-redefined
          env2: env2-val
  `
	app, _ := parseProcfile([]byte(data))

	assert.Contains(t, app.Services, Service{
		Name: "cmd2",
		Cmd:  "run-cmd2",
		Options: ServiceOptions{
			WorkingDirectory: "/working-dir2",
			Env:              map[string]string{"env1": "env1-val-redefined", "env2": "env2-val"},
		},
	})

	assert.Contains(t, app.Services, Service{
		Name: "cmd1",
		Cmd:  "run-cmd1",
		Options: ServiceOptions{
			KillTimeout: 60,
			Respawn: Respawn{
				Count:    5,
				Interval: 10,
			},
			LogPath:          "/path/to/log",
			Count:            2,
			WorkingDirectory: "/working-dir",
			Env:              map[string]string{"env1": "env1-val"},
		},
	})

}
