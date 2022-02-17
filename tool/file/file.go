package file

import (
	"encoding/json"
	"io/fs"
	"io/ioutil"
	"os"
	"tool/convert"
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
func ReadByByteList(path string) []byte {
	bList, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return bList
}
func ReadByString(path string) string {
	return string(ReadByByteList(path))
}
func ReadByInt(path string) int {
	return convert.StringToInt(ReadByString(path))
}
func ReadByJson(path string, v interface{}) {
	if err := json.Unmarshal(ReadByByteList(path), v); err != nil {
		panic(err)
	}
}
func SaveByString(path string, s string, mode fs.FileMode) {
	if err := ioutil.WriteFile(path, []byte(s), mode); err != nil {
		panic(err)
	}
}
func SaveByInt(path string, i int, mode fs.FileMode) {
	SaveByString(path, convert.IntToString(i), mode)
}
