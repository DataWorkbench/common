package schargs

// Encode for re-build the string with value-map
func Encode(data string, valueMap map[string]string) (string, error) {
	return newDecoder(data).encode(valueMap)
}

// ExtractVariables extract all valid variables from giving string.
func ExtractVariables(data string) ([]string, error) {
	return newDecoder(data).extractVariables()
}
