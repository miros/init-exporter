package systemd

import (
  "path/filepath"
  "os"
)

func Uninstall(appName string, config Config) {
  uninstallHelpers(appName, config)
  uninstallUnits(appName, config)

  DisableService(appName)
}

func uninstallUnits(appName string, config Config) {
  pattern := config.unitPath(appMask(appName))
  deleteByMask(pattern)
}

func uninstallHelpers(appName string, config Config) {
  pattern := config.helperPath(appMask(appName))
  deleteByMask(pattern)

  deleteByMask(config.unitPath(appName))
}

func appMask(appName string) string {
  return appName + "_*"
}

func deleteByMask(pattern string) {
  for _, path := range mustGlob(pattern) {
    err := os.Remove(path)

    if (err != nil) {
      panic(err)
    }
  }
}

func mustGlob(pattern string) []string {
  matches, err := filepath.Glob(pattern)

  if (err != nil) {
    panic(err)
  }

  return matches
}