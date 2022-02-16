package file

import (
	"os"
)

func IsExists(path string) bool {
	f, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !f.IsDir()
}

func Unlink(path string) {
	if err := os.Remove(path); err != nil {
		panic(err)
	}
}
