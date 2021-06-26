package utils

import "regexp"

const REPLACE_REG = "(?i)huaweicloud"

var PackageName string

var re = regexp.MustCompile(REPLACE_REG)

func BuildNewFormatByConfig(format string) string {
	newFormat := format
	if PackageName != "" {
		newFormat = re.ReplaceAllString(format, PackageName)
	}
	return newFormat
}
