package funcutil

import "errors"

// ContainStr check if the str is in strArray
func ContainStr(strArray []string, str string) bool {
	for _, s := range strArray {
		if s == str {
			return true
		}
	}
	return false
}

// StringDiffSet calculate the difference set of slice a and slice b
func StringDiffSet(a []string, b []string) []string {
	m := make(map[string]struct{}, len(b))
	for _, s := range b {
		m[s] = struct{}{}
	}
	var c []string
	for _, s := range a {
		if _, ok := m[s]; !ok {
			c = append(c, s)
		}
	}
	return c
}

func InterfaceToStringArray(i interface{}) ([]string, error) {
	varray, ok := i.([]interface{})
	if !ok {
		return nil, errors.New("the format of interface is not array")
	}
	var arr []string
	for _, v := range varray {
		str, ok := v.(string)
		if !ok {
			return nil, errors.New("the format of interface-value is not string")
		}
		arr = append(arr, str)
	}
	return arr, nil
}
