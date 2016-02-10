package systemd

func (sys *Systemd) Uninstall(appName string) {
  sys.uninstallHelpers(appName)
  sys.uninstallUnits(appName)
  sys.DisableService(appName)
}

func (sys *Systemd) uninstallUnits(appName string) {
  pattern := sys.Config.unitPath(appMask(appName))
  sys.deleteByMask(pattern)
}

func (sys *Systemd) uninstallHelpers(appName string) {
  pattern := sys.Config.helperPath(appMask(appName))
  sys.deleteByMask(pattern)

  sys.deleteByMask(sys.Config.unitPath(appName))
}

func appMask(appName string) string {
  return appName + "_*"
}

func (sys *Systemd) deleteByMask(pattern string) {
  for _, path := range sys.mustGlob(pattern) {
    err := sys.fs.Remove(path)

    if (err != nil) {
      panic(err)
    }
  }
}

func (sys *Systemd) mustGlob(pattern string) []string {
  matches, err := sys.globFilepath(pattern)

  if (err != nil) {
    panic(err)
  }

  return matches
}

