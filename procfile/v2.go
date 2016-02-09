package procfile

import (
  "systemd-exporter/systemd"
  "errors"
  "fmt"
  "github.com/imdario/mergo"
  "github.com/smallfish/simpleyaml"
)

func parseProcfileV2(data []byte) (services []systemd.Service, err error) {
  // TODO this is too long: refactor

  yaml, err := simpleyaml.NewYaml(data)

  if err != nil {
    return
  }

  commonOptions := getServiceOptions(yaml)

  commands, err := yaml.Get("commands").Map()
  if (err != nil) {
    err = errors.New("commands missing in Procfile")
    return
  }

  for key, _ := range(commands) {
    name := toString(key)
    commandYaml := yaml.GetPath("commands", name)

    cmd, _ := commandYaml.Get("command").String()

    service := systemd.Service{
      Name: name,
      Cmd: cmd,
      Options: getServiceOptions(commandYaml),
    }
    mergo.Merge(&service.Options, commonOptions)

    services = append(services, service)
  }

  return services, nil
}

func getServiceOptions(yaml *simpleyaml.Yaml) systemd.ServiceOptions {

  options := systemd.ServiceOptions{}

  if value := yaml.Get("working_directory"); isPresent(value) {
    options.WorkingDirectory, _ = value.String()
  }

  if value := yaml.Get("user"); isPresent(value) {
    options.User, _ = value.String()
  }

  if value := yaml.Get("group"); isPresent(value) {
    options.Group, _ = value.String()
  }

  if value := yaml.Get("kill_timeout"); isPresent(value) {
    options.KillTimeout, _ = value.Int()
  }

  if value := yaml.Get("env"); isPresent(value) {
    if envMap, err := value.Map(); err == nil {
      options.Env = toStringMap(envMap)
    }
  }

  if value := yaml.Get("respawn"); isPresent(value) {
    count, _ := value.Get("count").Int()
    interval, _ := value.Get("interval").Int()

    options.Respawn = systemd.Respawn{Count: count, Interval: interval}
  }

  return options
}

func isPresent(yaml *simpleyaml.Yaml) bool {
  empty := simpleyaml.Yaml{}
  return *yaml != empty
}

func toStringMap(sourceMap map[interface{}]interface{}) map[string]string {
  results := make(map[string]string)

  for key, value := range(sourceMap) {
    results[key.(string)] = toString(value)
  }

  return results
}

func toString(value interface{}) string {
  return fmt.Sprintf("%v", value)
}