package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/miros/init-exporter/environment"
	"github.com/miros/init-exporter/exporter"
	"github.com/miros/init-exporter/procfile"
	"github.com/miros/init-exporter/systemd"
	"github.com/miros/init-exporter/upstart"
	"os"
)

import "github.com/davecgh/go-spew/spew"

var _ = spew.Dump

const version = "0.0.2"
const defaultConfigPath = "/etc/init-exporter.yaml"

func main() {
	defer prettyPrintPanics()

	app := cli.NewApp()
	describeApp(app, version)
	app.Action = runAction
	app.Run(os.Args)
}

func prettyPrintPanics() {
	if os.Getenv("DEBUG") == "true" {
		return
	}

	if r := recover(); r != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", r)
		os.Exit(1)
	}
}

func describeApp(app *cli.App, version string) {
	app.Name = "init-exporter"
	app.Usage = "exports services described by Procfile to systemd"
	app.Version = version

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "n, appname",
			Usage: "Application name (This name only affects the names of generated files)",
		},
		cli.BoolFlag{
			Name:  "c, uninstall",
			Usage: "Remove scripts and helpers for a particular application",
		},
		cli.StringFlag{
			Name:  "config",
			Value: defaultConfigPath,
			Usage: "path to configuration file",
		},
		cli.StringFlag{
			Name:  "p, procfile",
			Usage: "path to procfile",
		},
		cli.StringFlag{
			Name:  "f, format",
			Usage: "Format of init files (upstart | systemd)",
		},
	}
}

func runAction(cliContext *cli.Context) {
	appName := cliContext.String("appname")

	if appName == "" {
		panic("No application name specified")
		return
	}

	globalConfig := environment.ReadGlobalConfig(cliContext.String("config"))
	appName = globalConfig.Prefix + appName

	providerName := environment.DetectProvider(cliContext.String("format"))
	exporter := newExporter(globalConfig, providerName)

	if cliContext.Bool("uninstall") {
		uninstall(exporter, appName)
	} else {
		install(exporter, appName, cliContext.String("procfile"))
	}
}

func newExporter(config environment.GlobalConfig, providerName string) *exporter.Exporter {
	exporterConfig := exporter.Config{
		HelperDir: config.HelperDir,
		TargetDir: config.TargetDirFor(providerName),
		User:      config.RunUser,
		Group:     config.RunGroup,
		DefaultWorkingDirectory: config.WorkingDirectory,
	}

	provider := newProvider(providerName)
	return exporter.New(exporterConfig, provider)
}

func newProvider(providerName string) exporter.Provider {
	switch providerName {
	case environment.SYSTEMD:
		return systemd.New()
	case environment.UPSTART:
		return upstart.New()
	default:
		panic("unknown init provider " + providerName)
	}
}

func uninstall(exporter *exporter.Exporter, appName string) {
	exporter.Uninstall(appName)
	fmt.Println("service uninstalled")
}

func install(exporter *exporter.Exporter, appName string, pathToProcfile string) {
	if pathToProcfile == "" {
		panic("No procfile given")
	}

	if app, err := procfile.ReadProcfile(pathToProcfile); err == nil {
		exporter.Install(appName, app)
		fmt.Println("service installed to", exporter.Config.TargetDir)
	} else {
		panic(err)
	}
}
