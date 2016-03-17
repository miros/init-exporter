package procfile

import (
	"github.com/miros/init-exporter/utils/validation"
)

type Respawn struct {
	Interval int
	Count    int
}

type ServiceOptions struct {
	WorkingDirectory string
	Env              map[string]string
	User             string
	Group            string
	KillTimeout      int
	Respawn          Respawn
	LogPath          string
	Count            int
}

func (options *ServiceOptions) Validate() error {
	if err := validation.Path(options.WorkingDirectory); err != nil {
		return err
	}

	if err := validation.NoSpecialSymbols(options.User); err != nil {
		return err
	}

	if err := validation.NoSpecialSymbols(options.Group); err != nil {
		return err
	}

	if err := validation.Path(options.LogPath); err != nil {
		return err
	}

	return nil
}

type Service struct {
	Name       string
	Cmd        string
	Options    ServiceOptions
	HelperPath string
}

func (service *Service) Validate() error {
	if err := validation.NoSpecialSymbols(service.Name); err != nil {
		return err
	}

	if err := service.Options.Validate(); err != nil {
		return err
	}

	return nil
}

func (service *Service) FullName(appName string) string {
	return appName + "_" + service.Name
}
