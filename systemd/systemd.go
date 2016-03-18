package systemd

import (
	"github.com/miros/init-exporter/utils"
)

type Systemd struct {
	execSystemCommand utils.SystemExecutor
}

func New() *Systemd {
	return &Systemd{
		execSystemCommand: utils.ExecSystemCommand,
	}
}

func (self *Systemd) UnitName(name string) string {
	return name + ".service"
}
