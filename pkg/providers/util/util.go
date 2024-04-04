package util

import "fmt"

func Key(z, s, n string) string {
	return fmt.Sprintf("%s-%s-%s", z, s, n)
}
