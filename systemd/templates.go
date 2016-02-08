package systemd

import (
  "text/template"
  "bytes"
)

const helperTemplate = `#!/bin/bash
[[ -r /etc/profile.d/rbenv.sh ]] && source /etc/profile.d/rbenv.sh
{{.cmd}}
`

func renderHelperTemplate(cmd string) string {
  data := make(map[string]interface{})
  data["cmd"] = cmd

  return renderTemplate("helper", data)
}

const appTemplate = `[Unit]
Description={{.app_name}}
After=network.target

[Service]
Type=oneshot
RemainAfterExit=true

ExecStartPre=mkdir -p /var/log/{{.app_name}}
ExecStartPre=chown -R {{.user}} /var/log/{{.app_name}}
ExecStartPre=chgrp -R {{.group}} /var/log/{{.app_name}}
ExecStartPre=chmod -R g+w /var/log/{{.app_name}}

[Install]
WantedBy=multi-user.target
`
func renderAppTemplate(appName string, config Config) string {
  data := make(map[string]interface{})
  data["app_name"] = appName
  data["user"] = config.User
  data["group"] = config.Group

  return renderTemplate("app", data)
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

ExecStartPre=touch /var/log/{{.app_name}}/{{.cmd_name}}.log
ExecStartPre=chown {{.user}} /var/log/{{.app_name}}/{{.cmd_name}}.log
ExecStartPre=chgrp {{.group}} /var/log/{{.app_name}}/{{.cmd_name}}.log
ExecStartPre=mkdir chmod g+w /var/log/{{.app_name}}/{{.cmd_name}}.log

ExecStart=exec sudo -u {{.user}} /bin/sh {{.helper_path}} >> /var/log/{{.app_name}}/{{.cmd_name}}.log 2>&1
`

func renderServiceTemplate(appName string, service Service) string {
  data := make(map[string]interface{})
  data["app_name"] = appName
  data["cmd_name"] = service.Name
  data["kill_timeout"] = service.Options.KillTimeout
  data["respawn_interval"] = service.Options.Respawn.Interval
  data["respawn_count"] = service.Options.Respawn.Count
  data["user"] = service.Options.User
  data["group"] = service.Options.Group
  data["helper_path"] = service.helperPath

  return renderTemplate("service", data)
}

func renderTemplate(Name string, data map[string]interface{}) string {
  var templates = map[string]string{
    "helper": helperTemplate,
    "app": appTemplate,
    "service": serviceTemplate,
  }

  tmpl, err := template.New(Name).Parse(templates[Name])
  if err != nil {
    panic(err)
  }
  return renderTemplateToString(tmpl, data)
}

func renderTemplateToString(template *template.Template, data interface{}) string {
  buffer := new(bytes.Buffer)

  err := template.Execute(buffer, data)
  if err != nil {
    panic(err)
  }

  return buffer.String()
}