package systemd

import (
  "os/exec"
)

func EnableService(appName string) error {
  return exec.Command("systemctl", "enable", serviceName(appName)).Run();
}

func MustEnableService(appName string) {
  must(EnableService(appName))
}

func DisableService(appName string) error {
  return exec.Command("systemctl", "disable", serviceName(appName)).Run();
}

func MustDisableService(appName string) {
  must(DisableService(appName))
}

func serviceName(name string) string {
  return name + ".service"
}

func must(err error) {
  if err != nil {
    panic(err)
  }
}