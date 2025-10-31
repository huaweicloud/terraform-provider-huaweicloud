package logp

import (
	"log"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// Deprecated: Please use log.Printf instead and don't contain `huaweicloud` in the format string
func Printf(format string, v ...interface{}) {
	newFormat := utils.BuildNewFormatByConfig(format)
	log.Printf(newFormat, v...)
}
