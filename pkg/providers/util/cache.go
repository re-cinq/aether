package util

import "fmt"

func CacheKey(z, s, n string) string {
	return fmt.Sprintf("%s-%s-%s", z, s, n)
}
