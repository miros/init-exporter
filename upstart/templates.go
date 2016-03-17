package upstart

import (
	"github.com/miros/init-exporter/exporter"
	"github.com/miros/init-exporter/procfile"
)

const helperTemplate = `#!/bin/bash
[[ -r /etc/profile.d/rbenv.sh ]] && source /etc/profile.d/rbenv.sh
cd {{.working_directory}} && exec {{.env}} {{.cmd}}
`

func (self *Upstart) RenderHelperTemplate(service procfile.Service) string {
	data := make(map[string]interface{})

	data["working_directory"] = service.Options.WorkingDirectory
	data["env"] = exporter.RenderEnvClause(service.Options.Env)
	data["cmd"] = service.Cmd

	return exporter.RenderTemplate("helper", helperTemplate, data)
}

const appTemplate = `
start on {{.start_on}}
stop on {{.stop_on}}

pre-start script

bash << "EOF"
  mkdir -p /var/log/{{.app_name}}
  chown -R {{.user}} /var/log/{{.app_name}}
  chgrp -R {{.group}} /var/log/{{.app_name}}
  chmod -R g+w /var/log/{{.app_name}}
EOF

end script
`

func (self *Upstart) RenderAppTemplate(appName string, config exporter.Config, services []procfile.Service) string {
	data := make(map[string]interface{})

	data["app_name"] = appName
	data["user"] = config.User
	data["group"] = config.Group
	data["start_on"] = "[3]"
	data["stop_on"] = "[3]"

	return exporter.RenderTemplate("app", appTemplate, data)
}

const serviceTemplate = `
start on {{.start_on}}
stop on {{.stop_on}}
respawn
respawn limit {{.respawn_count}} {{.respawn_interval}}
kill timeout {{.kill_timeout}}

script
  touch {{.log_path}}
  chown {{.user}} {{.log_path}}
  chgrp {{.group}} {{.log_path}}
  chmod g+w {{.log_path}}
  exec sudo -u {{.user}} /bin/sh {{.helper_path}} >> {{.log_path}} 2>&1
end script
`

func (self *Upstart) RenderServiceTemplate(appName string, service procfile.Service) string {
	data := make(map[string]interface{})

	data["app_name"] = appName
	data["cmd_name"] = service.Name
	data["kill_timeout"] = service.Options.KillTimeout
	data["respawn_interval"] = service.Options.Respawn.Interval
	data["respawn_count"] = service.Options.Respawn.Count
	data["user"] = service.Options.User
	data["group"] = service.Options.Group
	data["helper_path"] = service.HelperPath
	data["working_directory"] = service.Options.WorkingDirectory
	data["log_path"] = service.Options.LogPath
	data["start_on"] = "[3]"
	data["stop_on"] = "[3]"

	return exporter.RenderTemplate("service", serviceTemplate, data)
}
