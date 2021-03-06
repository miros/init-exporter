package procfile

import (
	"errors"
	"fmt"
	"github.com/imdario/mergo"
	"github.com/miros/simpleyaml"
)

func parseProcfileV2(data []byte) (app App, err error) {
	yaml, err := simpleyaml.NewYaml(data)
	if err != nil {
		return
	}

	commands, err := yaml.Get("commands").Map()
	if err != nil {
		err = errors.New("commands missing in Procfile")
		return
	}

	app.Services = parseCommands(yaml, commands)
	app.StartLevel, _ = yaml.Get("start_on_runlevel").String()
	app.StopLevel, _ = yaml.Get("stop_on_runlevel").String()

	return app, nil
}

func parseCommands(yaml *simpleyaml.Yaml, commands map[interface{}]interface{}) []Service {
	commonOptions := getServiceOptions(yaml)
	services := make([]Service, 0, len(commands))

	for key, _ := range commands {
		name := toString(key)
		commandYaml := yaml.GetPath("commands", name)
		cmd, _ := commandYaml.Get("command").String()

		service := Service{
			Name:    name,
			Cmd:     cmd,
			Options: getServiceOptions(commandYaml),
		}
		mergo.Merge(&service.Options, commonOptions)

		services = append(services, service)
	}

	return services
}

func getServiceOptions(yaml *simpleyaml.Yaml) ServiceOptions {
	options := ServiceOptions{}

	options.WorkingDirectory, _ = yaml.Get("working_directory").String()
	options.User, _ = yaml.Get("user").String()
	options.Group, _ = yaml.Get("group").String()
	options.KillTimeout = mustGetInt(yaml, "kill_timeout")
	options.LogPath, _ = yaml.Get("log").String()
	options.Count = mustGetInt(yaml, "count")

	if value, err := yaml.Get("env").Map(); err == nil {
		options.Env = toStringMap(value)
	}

	if value := yaml.Get("respawn"); isPresent(value) {
		count := mustGetInt(value, "count")
		interval := mustGetInt(value, "interval")
		options.Respawn = Respawn{Count: count, Interval: interval}
	}

	return options
}

func mustGetInt(yaml *simpleyaml.Yaml, key string) int {
	if rawVal := yaml.Get(key); isPresent(rawVal) {
		if val, err := rawVal.Int(); err != nil {
			panic(fmt.Sprintf("%s is not integer", key))
		} else {
			return val
		}
	}

	return 0
}

func isPresent(yaml *simpleyaml.Yaml) bool {
	empty := simpleyaml.Yaml{}
	return *yaml != empty
}

func toStringMap(sourceMap map[interface{}]interface{}) map[string]string {
	results := make(map[string]string)

	for key, value := range sourceMap {
		results[key.(string)] = toString(value)
	}

	return results
}

func toString(value interface{}) string {
	return fmt.Sprintf("%v", value)
}
