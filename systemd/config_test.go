package systemd

import (
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestConfigValidations(t *testing.T) {
  var config Config

  config = Config{User: "!wrong!"}
  assert.Error(t, config.Validate())

  config = Config{Group: "!wrong!"}
  assert.Error(t, config.Validate())

  config = Config{DefaultWorkingDirectory: "!wrong!"}
  assert.Error(t, config.Validate())
}

