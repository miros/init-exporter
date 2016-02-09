package main

import (
  "os"
  "github.com/codegangsta/cli"
  "systemd-exporter/systemd"
  "systemd-exporter/procfile"
)

// import "github.com/davecgh/go-spew/spew"

const defaultConfigPath = "/etc/systemd-exporter.yaml"

func RunCLI(version string) {
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
  systemdConfig := newSystemdConfig(globalConfig)

  if cliContext.Bool("uninstall") {
    uninstall(appName, systemdConfig)
  } else {
    install(appName, systemdConfig, cliContext.String("procfile"))
  }
}

func uninstall(appName string, config systemd.Config) {
  systemd.Uninstall(appName, config)
  println("systemd service uninstalled")
}

func install(appName string, systemdConfig systemd.Config, pathToProcfile string) {
  if (pathToProcfile == "") {
    panic("No procfile given")
  }

  if services, err := procfile.ReadProcfile(pathToProcfile); err == nil {
    systemd.InstallAndEnable(appName, systemdConfig, services)
    println("systemd service installed")
  } else {
    panic(err)
  }
}

func newSystemdConfig(config GlobalConfig) systemd.Config {
  return systemd.Config{
    HelperDir: config.HelperDir,
    TargetDir: config.TargetDir,
    User: config.RunUser,
    Group: config.RunGroup,
  }
}