package fmtp

import (
	"fmt"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func Sprintf(format string, a ...interface{}) string {
	newFormat := utils.BuildNewFormatByConfig(format)
	return fmt.Sprintf(newFormat, a)
}
