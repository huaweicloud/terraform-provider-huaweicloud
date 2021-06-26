package fmtp

import (
	"fmt"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func Errorf(format string, a ...interface{}) error {
	newFormat := utils.BuildNewFormatByConfig(format)
	return fmt.Errorf(newFormat, a)
}
