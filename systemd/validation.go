package systemd

import (
  "regexp"
  "fmt"
  "errors"
)

type validatable interface {
    validate() error
}

func mustBeValid(item validatable) {
  if err := item.validate(); err != nil {
    panic(err)
  }
}

func validatePath(val string) error {
  return validateString(val, `\A[A-Za-z0-9_\-./]+\z`)
}

func validateNoSpecialSymbols(val string) error {
  return validateString(val, `\A[A-Za-z0-9_\-]+\z`)
}

func validateString(val string, reString string) error {
  if val == "" {
    return nil
  }

  if re := regexp.MustCompile(reString); !re.MatchString(val) {
    return errors.New(fmt.Sprintf("value %s is insecure and can't be accepted", val))
  }

  return nil
}