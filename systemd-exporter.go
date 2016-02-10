package main

import (
  "os"
  "github.com/codegangsta/cli"
  "systemd-exporter/systemd"
  "systemd-exporter/procfile"
  "fmt"
)

import "github.com/davecgh/go-spew/spew"
var _ = spew.Dump

const version = "0.0.1"
const defaultConfigPath = "/etc/systemd-exporter.yaml"

func main() {
  app := cli.NewApp()
  describeApp(app, version)
  app.Action = runAction
  app.Run(os.Args)
}

func describeApp(app *cli.App, version string) {
  app.Name = "systemd-exporter"
  app.Usage = "exports services described by Procfile to systemd"
  app.Version = version

  app.Flags = []cli.Flag {
    cli.StringFlag{
      Name: "n, application_name",
      Usage: "Application name (This name only affects the names of generated files)",
    },
    cli.BoolFlag{
      Name: "c, uninstall",
      Usage: "Remove scripts and helpers for a particular application",
    },
    cli.StringFlag{
      Name: "config",
      Value: defaultConfigPath,
      Usage: "path to configuration file",
    },
    cli.StringFlag{
      Name: "p, procfile",
      Usage: "path to procfile",
    },
  }
}

func runAction(cliContext *cli.Context) {
  appName := cliContext.String("application_name")

  if appName == "" {
    panic("No application name specified")
    return
  }

  globalConfig := ReadGlobalConfig(cliContext.String("config"))
  appName = globalConfig.Prefix + appName;
  systemd := newSystemd(globalConfig)

  if cliContext.Bool("uninstall") {
    uninstall(systemd, appName)
  } else {
    install(systemd, appName, cliContext.String("procfile"))
  }
}

func newSystemd(config GlobalConfig) *systemd.Systemd {
  systemdConfig := systemd.Config{
    HelperDir: config.HelperDir,
    TargetDir: config.TargetDir,
    User: config.RunUser,
    Group: config.RunGroup,
    DefaultWorkingDirectory: config.WorkingDirectory,
  }

  return systemd.New(systemdConfig)
}

func uninstall(systemd *systemd.Systemd, appName string) {
  systemd.Uninstall(appName)
  fmt.Println("systemd service uninstalled")
}

func install(systemd *systemd.Systemd, appName string, pathToProcfile string) {
  if (pathToProcfile == "") {
    panic("No procfile given")
  }

  if services, err := procfile.ReadProcfile(pathToProcfile); err == nil {
    systemd.Install(appName, services)
    fmt.Println("systemd service installed to", systemd.Config.TargetDir)
  } else {
    panic(err)
  }
}
