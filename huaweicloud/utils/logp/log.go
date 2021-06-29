package logp

import (
	"log"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func Printf(format string, v ...interface{}) {
	newFormat := utils.BuildNewFormatByConfig(format)
	log.Printf(newFormat, v...)
}
