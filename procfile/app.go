package procfile

import (
	"github.com/miros/init-exporter/utils/validation"
)

type App struct {
	Services   []Service
	StartLevel string
	StopLevel  string
}

func (app App) Validate() error {
	if err := validation.RunLevel(app.StartLevel); err != nil {
		return err
	}

	if err := validation.RunLevel(app.StopLevel); err != nil {
		return err
	}

	for _, service := range app.Services {
		if err := service.Validate(); err != nil {
			return err
		}
	}

	return nil
}
