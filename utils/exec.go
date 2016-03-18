package utils

import (
	"os/exec"
)

type SystemExecutor func(name string, arg ...string) error

func ExecSystemCommand(name string, arg ...string) error {
	return exec.Command(name, arg...).Run()
}
