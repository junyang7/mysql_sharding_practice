package convert

import "strconv"

func StringToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}

func IntToString(i int) string {
	return strconv.Itoa(i)
}
