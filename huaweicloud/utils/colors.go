package utils

import "fmt"

const (
	greenCode  = "\033[1;32m"
	yellowCode = "\033[1;33m"
	resetCode  = "\033[0m"
)

func Green(str interface{}) string {
	return fmt.Sprintf("%s%#v%s", greenCode, str, resetCode)
}

func Yellow(str interface{}) string {
	return fmt.Sprintf("%s%#v%s", yellowCode, str, resetCode)
}
