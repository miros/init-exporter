package environment

import (
	"os"
	"os/exec"
)

const SYSTEMD = "systemd"
const UPSTART = "upstart"

func DetectProvider(providerName string) string {
	if providerName != "" {
		return providerName
	}

	if providerName = detectByCurrentExecutable(); providerName != "" {
		return providerName
	}

	if providerName = detectByInstalledProvider(); providerName != "" {
		return providerName
	}

	panic("init format was not detected; explicitly pass --format=(upstart|systemd)")
}

func detectByCurrentExecutable() string {
	switch os.Args[0] {
	case "systemd-exporter":
		return SYSTEMD
	case "upstart-exporter":
		return UPSTART
	default:
		return ""
	}
}

func detectByInstalledProvider() string {
	switch {
	case executableExists("service"):
		return UPSTART
	case executableExists("systemctl"):
		return SYSTEMD
	default:
		return ""
	}
}

func executableExists(executable string) bool {
	return exec.Command("which", executable).Run() == nil
}
