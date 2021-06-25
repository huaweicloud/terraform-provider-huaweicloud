package logp

import (
	"log"
	"strings"
)

const REPLACE_STR = "huawei"

func Printf(format string, v ...interface{}) {
	newFormat := strings.Replace(format, REPLACE_STR, "", -1)
	log.Printf(newFormat, v)

}
