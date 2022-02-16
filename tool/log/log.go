package log

import (
	"fmt"
	"tool/datetime"
)

func Info(message ...interface{}) {
	fmt.Println(datetime.Get(), message)
}
