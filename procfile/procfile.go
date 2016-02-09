package procfile

import (
  "systemd-exporter/systemd"
  "io/ioutil"
  "regexp"
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