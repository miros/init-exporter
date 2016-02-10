package systemd

import (
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestServiceValidations(t *testing.T) {
  var service Service

  service = Service{Name: "!wrong!"}
  assert.Error(t, service.Validate())

  service = Service{Options: ServiceOptions{User: "!wrong!"}}
  assert.Error(t, service.Validate())

  service = Service{Options: ServiceOptions{Group: "!wrong!"}}
  assert.Error(t, service.Validate())

  service = Service{Options: ServiceOptions{WorkingDirectory: "!wrong!"}}
  assert.Error(t, service.Validate())
}

