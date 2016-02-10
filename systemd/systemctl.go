package systemd

func (sys *Systemd) EnableService(appName string) error {
  return sys.execSystemCommand("systemctl", "enable", serviceName(appName));
}

func (sys *Systemd) MustEnableService(appName string) {
  must(sys.EnableService(appName))
}

func (sys *Systemd) DisableService(appName string) error {
  return sys.execSystemCommand("systemctl", "disable", serviceName(appName));
}

func (sys *Systemd) MustDisableService(appName string) {
  must(sys.DisableService(appName))
}

func serviceName(name string) string {
  return name + ".service"
}

func must(err error) {
  if err != nil {
    panic(err)
  }
}

