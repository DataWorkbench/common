package utils

type StrArray []string

func (sa StrArray) Contain(str string) bool {
	for _, s := range sa {
		if s == str {
			return true
		}
	}
	return false
}
