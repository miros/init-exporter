package utils

import (
	"github.com/spf13/afero"
)

func MustWriteFile(fs afero.Fs, path string, data string) {
	error := afero.WriteFile(fs, path, []byte(data), 0644)
	if error != nil {
		panic(error)
	}
}
