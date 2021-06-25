package logp

import (
	"log"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const REPLACE_STR = "huawei"

func Printf(format string, v ...interface{}) {
	newFormat := utils.BuildNewFormatByConfig(format)
	log.Printf(newFormat, v)
}
