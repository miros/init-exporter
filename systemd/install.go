package systemd

import (
  "io/ioutil"
  "os"
  "github.com/imdario/mergo"
)

import "github.com/davecgh/go-spew/spew"
var _ = spew.Dump

func InstallAndEnable(appName string, config Config, services []Service) {
  Install(appName, config, services)
  MustEnableService(appName)
}

func Install(appName string, config Config, services []Service) {
  setServiceDefaults(services, config)

  validateAppName(appName)
  mustBeValid(&config)
  validateServices(services)

  installServices(appName, config, services)
  writeAppUnit(appName, config, services)
}

func validateAppName(appName string) {
  if err := validateNoSpecialSymbols(appName); err != nil {
    panic(err)
  }
}

func installServices(appName string, config Config, services []Service) {
  error := os.MkdirAll(config.HelperDir, 0755)
  if error != nil {
    panic(error)
  }

  for _, service := range(services) {
    writeServiceUnit(appName, config, service)
  }
}

func setServiceDefaults(services []Service, config Config) {
  for i, _ := range services {
    defaults := ServiceOptions{User: config.User, Group: config.Group, WorkingDirectory: config.WorkingDirectory}
    mergo.Merge(&services[i].Options, defaults)
  }
}

func writeAppUnit(appName string, config Config, services []Service) {
  path := config.unitPath(appName)
  data := renderAppTemplate(appName, config, services)
  writeFile(path, data)
}

func writeServiceUnit(appName string, config Config, service Service) {
  fullServiceName := service.fullName(appName)

  service.helperPath = config.helperPath(fullServiceName)
  helperData := renderHelperTemplate(service.Cmd)
  writeFile(service.helperPath, helperData)

  unitPath := config.unitPath(fullServiceName)
  writeFile(unitPath, renderServiceTemplate(appName, service))
}

func writeFile(path string, data string) {
  error := ioutil.WriteFile(path, []byte(data), 0644)
  if error != nil {
    panic(error)
  }
}
