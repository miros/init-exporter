package exporter

import (
	"github.com/miros/go-glob"
	"github.com/spf13/afero"
	"path"
)

func (self *Exporter) Uninstall(appName string) {
	self.uninstallHelpers(appName)
	self.uninstallUnits(appName)
	self.provider.DisableService(appName)
}

func (self *Exporter) uninstallUnits(appName string) {
	self.deleteByMask(self.Config.TargetDir, self.provider.UnitName(appName+"_*"))
	self.deleteByMask(self.Config.TargetDir, self.provider.UnitName(appName))
}

func (self *Exporter) uninstallHelpers(appName string) {
	self.deleteByMask(self.Config.HelperDir, appName+"_*.sh")
}

func (self *Exporter) deleteByMask(path, pattern string) {
	for _, path := range globFiles(self.Fs, path, pattern) {
		err := self.Fs.Remove(path)

		if err != nil {
			panic(err)
		}
	}
}

func globFiles(fs afero.Fs, targetPath, pattern string) []string {
	children, err := afero.ReadDir(fs, targetPath)

	if err != nil {
		panic(err)
	}

	var matches []string

	for _, child := range children {
		if glob.Glob(pattern, child.Name()) {
			matches = append(matches, path.Join(targetPath, child.Name()))
		}
	}

	return matches
}
