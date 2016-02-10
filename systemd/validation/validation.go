package validation

import (
  "regexp"
  "fmt"
  "errors"
)

type Validatable interface {
    Validate() error
}

func MustBeValid(item Validatable) {
  if err := item.Validate(); err != nil {
    panic(err)
  }
}

func Path(val string) error {
  return validateString(val, `\A[A-Za-z0-9_\-./]+\z`)
}

func NoSpecialSymbols(val string) error {
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