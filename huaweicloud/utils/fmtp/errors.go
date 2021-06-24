package fmtp

import (
	"fmt"
	"strings"
)

func Errorf(format string, a ...interface{}) error {
	newFormat := strings.Replace(format, REPLACE_STR, "", -1)
	return fmt.Errorf(newFormat, a)
}
