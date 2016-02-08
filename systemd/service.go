package systemd

type Respawn struct {
  Count int
  Interval int
}

type ServiceOptions struct {
  WorkingDirectory string
  Env map[string]string
  User string
  Group string
  KillTimeout int
  Respawn Respawn
}

type Service struct {
  Name string
  Cmd string
  Options ServiceOptions
  helperPath string
}

func (service *Service) fullName(appName string) string {
  return appName + "_" + service.Name
}