package exporter

import (
	"github.com/miros/init-exporter/procfile"
)

type Provider interface {
	UnitName(name string) string

	RenderAppTemplate(appName string, config Config, services []procfile.Service) string
	RenderServiceTemplate(appName string, service procfile.Service) string
	RenderHelperTemplate(service procfile.Service) string

	MustEnableService(appName string)
	DisableService(appName string) error

	DefaultTargetDir() string
}
