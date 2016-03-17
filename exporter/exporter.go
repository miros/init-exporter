package exporter

import (
	"github.com/spf13/afero"
	"path"
)

type Exporter struct {
	Config   Config
	Fs       afero.Fs
	provider Provider
}

func New(config Config, provider Provider) *Exporter {
	return &Exporter{
		Config:   config,
		Fs:       afero.NewOsFs(),
		provider: provider,
	}
}

func (self *Exporter) UnitPath(name string) string {
	return path.Join(self.Config.TargetDir, self.provider.UnitName(name))
}
