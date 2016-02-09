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

func (service *Service) validate() error {
  if err := validateNoSpecialSymbols(service.Name); err != nil {
    return err
  }

  if err := service.Options.validate(); err != nil {
    return err
  }

  return nil
}

func (options *ServiceOptions) validate() error {
  if err := validatePath(options.WorkingDirectory); err != nil {
    return err
  }

  if err := validateNoSpecialSymbols(options.User); err != nil {
    return err
  }

  if err := validateNoSpecialSymbols(options.Group); err != nil {
    return err
  }

  return nil
}

func validateServices(services []Service) {
  for _, service := range(services) {
    mustBeValid(&service)
  }
}
