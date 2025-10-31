package fmtp

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// Deprecated: please use fmt.Errorf instead and don't contain `huaweicloud` in the format string
func Errorf(format string, a ...interface{}) error {
	newFormat := utils.BuildNewFormatByConfig(format)
	return fmt.Errorf(newFormat, a...)
}

// DiagErrorf wraps fmtp.Errorf into diag.Diagnostics

// Deprecated: please use diag.Errorf instead and don't contain `huaweicloud` in the format string
func DiagErrorf(format string, a ...interface{}) diag.Diagnostics {
	return diag.FromErr(Errorf(format, a...))
}
