package exception

import "strings"

func Have(message string, err error) bool {
	return strings.Contains(err.Error(), message)
}
