package systemd

import (
	"github.com/miros/init-exporter/exporter"
	"github.com/miros/init-exporter/procfile"
	"strings"
)

const helperTemplate = `#!/bin/bash
[[ -r /etc/profile.d/rbenv.sh ]] && source /etc/profile.d/rbenv.sh
exec {{.cmd}}
`

func (sys *Systemd) RenderHelperTemplate(service procfile.Service) string {
	data := make(map[string]interface{})
	data["cmd"] = service.Cmd

	return exporter.RenderTemplate("helper", helperTemplate, data)
}

const appTemplate = `[Unit]
Description={{.app_name}}
After=network.target
Wants={{.wants}}

[Service]
Type=oneshot
RemainAfterExit=true

ExecStartPre=/bin/mkdir -p /var/log/{{.app_name}}
ExecStartPre=/bin/chown -R {{.user}} /var/log/{{.app_name}}
ExecStartPre=/bin/chgrp -R {{.group}} /var/log/{{.app_name}}
ExecStartPre=/bin/chmod -R g+w /var/log/{{.app_name}}

ExecStart=/bin/echo "{{.app_name}} started"
ExecStop=/bin/echo "{{.app_name}} stoped"

[Install]
WantedBy=multi-user.target
`

func (sys *Systemd) RenderAppTemplate(appName string, config exporter.Config, app procfile.App) string {
	data := make(map[string]interface{})

	data["app_name"] = appName
	data["user"] = config.User
	data["group"] = config.Group
	data["wants"] = renderWantsClause(appName, app.Services)

	return exporter.RenderTemplate("app", appTemplate, data)
}

func renderWantsClause(appName string, services []procfile.Service) string {
	names := make([]string, 0, len(services))
	for _, service := range services {
		names = append(names, service.FullName(appName)+".service")
	}
	return strings.Join(names, " ")
}

const serviceTemplate = `[Unit]
Description={{.app_name}}/{{.cmd_name}}
PartOf={{.app_name}}.service

[Service]
Type=simple

TimeoutStopSec={{.kill_timeout}}
Restart=on-failure
StartLimitInterval={{.respawn_interval}}
StartLimitBurst={{.respawn_count}}

ExecStartPre=/bin/touch {{.log_path}}
ExecStartPre=/bin/chown {{.user}} {{.log_path}}
ExecStartPre=/bin/chgrp {{.group}} {{.log_path}}
ExecStartPre=/bin/chmod g+w {{.log_path}}

User={{.user}}
Group={{.group}}
WorkingDirectory={{.working_directory}}
Environment={{.env}}
ExecStart=/bin/sh {{.helper_path}} >> {{.log_path}} 2>&1
`

func (sys *Systemd) RenderServiceTemplate(appName string, service procfile.Service) string {
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
	data["env"] = exporter.RenderEnvClause(service.Options.Env)

	return exporter.RenderTemplate("service", serviceTemplate, data)
}
