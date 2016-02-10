package systemd

import (
  "text/template"
  "bytes"
  "strings"
)

const helperTemplate = `#!/bin/bash
[[ -r /etc/profile.d/rbenv.sh ]] && source /etc/profile.d/rbenv.sh
exec {{.cmd}}
`

func RenderHelperTemplate(cmd string) string {
  data := make(map[string]interface{})
  data["cmd"] = cmd

  return renderTemplate("helper", data)
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
func RenderAppTemplate(appName string, config Config, services []Service) string {
  data := make(map[string]interface{})
  data["app_name"] = appName
  data["user"] = config.User
  data["group"] = config.Group
  data["wants"] = renderWantsClause(appName, services)

  return renderTemplate("app", data)
}

func renderWantsClause(appName string, services []Service) string {
  names := make([]string, 0, len(services))
  for _, service := range(services) {
    names = append(names, service.fullName(appName) + ".service")
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

ExecStartPre=/bin/touch /var/log/{{.app_name}}/{{.cmd_name}}.log
ExecStartPre=/bin/chown {{.user}} /var/log/{{.app_name}}/{{.cmd_name}}.log
ExecStartPre=/bin/chgrp {{.group}} /var/log/{{.app_name}}/{{.cmd_name}}.log
ExecStartPre=/bin/chmod g+w /var/log/{{.app_name}}/{{.cmd_name}}.log

User={{.user}}
Group={{.group}}
WorkingDirectory={{.working_directory}}
Environment={{.env}}
ExecStart=/bin/sh {{.helper_path}} >> /var/log/{{.app_name}}/{{.cmd_name}}.log 2>&1
`

func RenderServiceTemplate(appName string, service Service) string {
  data := make(map[string]interface{})
  data["app_name"] = appName
  data["cmd_name"] = service.Name
  data["kill_timeout"] = service.Options.KillTimeout
  data["respawn_interval"] = service.Options.Respawn.Interval
  data["respawn_count"] = service.Options.Respawn.Count
  data["user"] = service.Options.User
  data["group"] = service.Options.Group
  data["helper_path"] = service.helperPath
  data["working_directory"] = service.Options.WorkingDirectory
  data["env"] = renderEnvClause(service.Options.Env)

  return renderTemplate("service", data)
}

func renderEnvClause(env map[string]string) string {
  clauses := make([]string, 0, len(env))
  for name, value := range(env) {
    clauses = append(clauses, name + "=" + value)
  }
  return strings.Join(clauses, " ")
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