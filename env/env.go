package env

import (
	"os"
	"strconv"
)

func OrDefault(key string, def string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return def
}

func OrDefaultInt(key string, def int) int {
	if v, ok := os.LookupEnv(key); ok {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}
