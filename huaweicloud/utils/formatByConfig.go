package utils

import "regexp"

const REPLACE_REG = "(?i)huawei"

var buildCloudCompany string

var re = regexp.MustCompile(REPLACE_REG)

func BuildNewFormatByConfig(format string) string {
	newFormat := format
	if buildCloudCompany != "" {
		newFormat = re.ReplaceAllString(format, buildCloudCompany)
	}
	return newFormat
}
