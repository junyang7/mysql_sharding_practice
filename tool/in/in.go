package in

func StringList(s string, sList []string) bool {
	for _, v := range sList {
		if s == v {
			return true
		}
	}
	return false
}
