package main

import (
  "systemd-exporter/systemd"
  "errors"
  "fmt"
  "io/ioutil"
  "regexp"
  "bytes"
  "bufio"
  "strings"
  "github.com/imdario/mergo"
  "github.com/smallfish/simpleyaml"
)

import "github.com/davecgh/go-spew/spew"
var _ = spew.Dump

func ReadProcfile(path string) (services []systemd.Service, err error) {
  data, err := ioutil.ReadFile(path)

  if err != nil {
    return
  }

  if isV2(data) {
    services, err = parseProcfileV2(data)
  } else {
    services, err = parseProcfileV1(data)
  }

  return
}

func isV2(data []byte) bool {
  re := regexp.MustCompile(`(?m)^version:\s*2\s*$`)
  return re.Find(data) != nil
}

func parseProcfileV1(data []byte) (services []systemd.Service, err error) {
  scanner := bufio.NewScanner(bytes.NewReader(data))

  for scanner.Scan() {
    line := strings.TrimSpace(scanner.Text())

    switch {
    case line == "":
      // nop
    case strings.HasPrefix(line, "#"):
      // comment
    default:
      if service := parseV1Line(line); service != nil {
        services = append(services, *service)
      } else {
        err = errors.New("procfile v1 should have format: 'some_label: command'")
      }
    }
  }

  if err := scanner.Err(); err != nil {
    panic(err)
  }

  return
}

func parseV1Line(line string) *systemd.Service {
  re := regexp.MustCompile(`^([A-z\d_]+):\s*(.+)`)
  matches := re.FindStringSubmatch(line)

  if len(matches) != 3 {
    return nil
  }

  name := matches[1]
  cmd := matches[2]

  return &systemd.Service{Name: name, Cmd: cmd}
}

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