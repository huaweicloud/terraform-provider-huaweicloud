package fmtp

import (
	"fmt"
	"strings"
)

const REPLACE_STR = "huawei"

func Sprintf(format string, a ...interface{}) string {
	newFormat := strings.Replace(format, REPLACE_STR, "", -1)
	return fmt.Sprintf(newFormat, a)
}
