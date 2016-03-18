package procfile

import (
	"io/ioutil"
	"regexp"
)

import "github.com/davecgh/go-spew/spew"

var _ = spew.Dump

func ReadProcfile(path string) (app App, err error) {
	data, err := ioutil.ReadFile(path)

	if err != nil {
		return
	}

	return parseProcfile(data)
}

func parseProcfile(data []byte) (app App, err error) {
	if isV2(data) {
		app, err = parseProcfileV2(data)
	} else {
		app.Services, err = parseProcfileV1(data)
	}

	return
}

func isV2(data []byte) bool {
	re := regexp.MustCompile(`(?m)^\s*version:\s*2\s*$`)
	return re.Find(data) != nil
}
