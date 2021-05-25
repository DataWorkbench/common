package constants

// JSONString used to convert a json format type to a string.
type JSONString string

func (s *JSONString) UnmarshalJSON(b []byte) error {
	*s = JSONString(b)
	return nil
}

func (s JSONString) MarshalJSON() ([]byte, error) {
	return []byte(s), nil
}
