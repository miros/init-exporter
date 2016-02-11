package procfile

import (
  "github.com/miros/systemd-exporter/systemd"
  "errors"
  "bytes"
  "bufio"
  "strings"
  "regexp"
)

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