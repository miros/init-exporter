package exporter

import (
	"fmt"
	"github.com/imdario/mergo"
	"github.com/miros/init-exporter/procfile"
	"github.com/miros/init-exporter/utils"
	"github.com/miros/init-exporter/utils/validation"
)

func (self *Exporter) Install(appName string, app procfile.App) {
	app.Services = handleServiceCounts(app.Services)
	setServiceDefaults(appName, app, self.Config)
	validateParams(appName, self.Config, app)

	self.doInstall(appName, app)
}

func (self *Exporter) doInstall(appName string, app procfile.App) {
	self.writeServices(appName, app.Services)
	self.writeAppUnit(appName, app)
	self.provider.MustEnableService(appName)
}

func (self *Exporter) writeServices(appName string, services []procfile.Service) {
	error := self.Fs.MkdirAll(self.Config.HelperDir, 0755)
	if error != nil {
		panic(error)
	}

	for _, service := range services {
		self.writeServiceUnit(appName, service)
	}
}

func (self *Exporter) writeAppUnit(appName string, app procfile.App) {
	path := self.UnitPath(appName)
	data := self.provider.RenderAppTemplate(appName, self.Config, app)
	utils.MustWriteFile(self.Fs, path, data)
}

func (self *Exporter) writeServiceUnit(appName string, service procfile.Service) {
	fullServiceName := service.FullName(appName)

	service.HelperPath = self.Config.HelperPath(fullServiceName)
	helperData := self.provider.RenderHelperTemplate(service)
	utils.MustWriteFile(self.Fs, service.HelperPath, helperData)

	unitPath := self.UnitPath(fullServiceName)
	utils.MustWriteFile(self.Fs, unitPath, self.provider.RenderServiceTemplate(appName, service))
}

func handleServiceCounts(services []procfile.Service) []procfile.Service {
	newServices := make([]procfile.Service, 0, len(services))

	for _, service := range services {
		if count := service.Options.Count; count > 1 {
			for i := 1; i <= count; i++ {
				newService := service
				newService.Name = fmt.Sprintf("%s%d", service.Name, i)
				newServices = append(newServices, newService)
			}
		} else {
			newServices = append(newServices, service)
		}
	}

	return newServices
}

func setServiceDefaults(appName string, app procfile.App, config Config) {
	for i, _ := range app.Services {
		service := &app.Services[i]

		defaults := procfile.ServiceOptions{
			User:             config.User,
			Group:            config.Group,
			WorkingDirectory: config.DefaultWorkingDirectory,
			LogPath:          fmt.Sprintf("/var/log/%s/%s.log", appName, service.Name),
			StartLevel:       app.StartLevel,
			StopLevel:        app.StopLevel,
		}
		mergo.Merge(&service.Options, defaults)
	}
}

func validateParams(appName string, config Config, app procfile.App) {
	validateAppName(appName)
	validateConfig(config)
	validateServices(app.Services)
}

func validateAppName(appName string) {
	if err := validation.NoSpecialSymbols(appName); err != nil {
		panic(err)
	}
}

func validateConfig(config Config) {
	validation.MustBeValid(&config)
}

func validateServices(services []procfile.Service) {
	for _, service := range services {
		validation.MustBeValid(&service)
	}
}
