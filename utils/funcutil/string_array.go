package funcutil

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
