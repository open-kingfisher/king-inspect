package util

import "strings"

func Contains(list []string, name string) bool {
	for _, l := range list {
		if strings.TrimSpace(l) == name {
			return true
		}
	}
	return false
}
