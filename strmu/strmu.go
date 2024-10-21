package strmu

import (
	"strings"
)

func Contains(s string, a ...string) bool {
	for _, v := range a {
		if strings.Contains(s, v) {
			return true
		}
	}
	return false
}
