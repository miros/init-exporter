package systemd

import (
  "path"
  "github.com/ryanuber/go-glob"
  "github.com/spf13/afero"
)

func (sys *Systemd) Uninstall(appName string) {
  sys.uninstallHelpers(appName)
  sys.uninstallUnits(appName)
  sys.DisableService(appName)
}

func (sys *Systemd) uninstallUnits(appName string) {
  sys.deleteByMask(sys.Config.TargetDir, appName + "_*.service")
  sys.deleteByMask(sys.Config.TargetDir, appName + ".service")
}

func (sys *Systemd) uninstallHelpers(appName string) {
  sys.deleteByMask(sys.Config.HelperDir, appName + "_*.sh")
}

func (sys *Systemd) deleteByMask(path, pattern string) {
  for _, path := range globFiles(sys.fs, path, pattern) {
    err := sys.fs.Remove(path)

    if (err != nil) {
      panic(err)
    }
  }
}

func globFiles(fs afero.Fs, targetPath, pattern string) []string {
  children, err := afero.ReadDir(fs, targetPath)

  if (err != nil) {
    panic(err)
  }

  var matches []string

  for _, child := range(children) {
    if glob.Glob(pattern, child.Name()) {
      matches = append(matches, path.Join(targetPath, child.Name()))
    }
  }

  return matches
}
