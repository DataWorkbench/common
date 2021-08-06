package utils

import "os"

func GetEnv(env, defaultVal string) string {
	val := os.Getenv(env)
	if val == "" {
		return defaultVal
	}
	return val
}
