package procfile

type App struct {
	Services   []Service
	StartLevel string
	StopLevel  string
}
