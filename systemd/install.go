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
  setServiceOptions(services, config)

  writeAppUnit(appName, config)

  error := os.MkdirAll(config.HelperDir, 0755)
  if error != nil {
    panic(error)
  }

  for _, service := range(services) {
    writeServiceUnit(appName, config, service)
  }
}

func setServiceOptions(services []Service, config Config) {
  for i, _ := range services {
    service := &services[i]
    mergo.Merge(&service.Options, config)
  }
}

func writeAppUnit(appName string, config Config) {
  path := config.unitPath(appName)
  data := renderAppTemplate(appName, config)
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
