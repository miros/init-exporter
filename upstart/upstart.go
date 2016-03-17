package upstart

type Upstart struct{}

func New() *Upstart {
	return &Upstart{}
}

func (_ *Upstart) UnitName(name string) string {
	return name + ".conf"
}

func (_ *Upstart) DefaultTargetDir() string {
	return "/etc/init/"
}

func (_ *Upstart) MustEnableService(appName string) {
	// nop
}

func (_ *Upstart) DisableService(appName string) error {
	// nop
	return nil
}
