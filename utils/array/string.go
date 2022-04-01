package array

func ContainStr(strArray []string, str string) bool {
	for _, s := range strArray {
		if s == str {
			return true
		}
	}
	return false
}
