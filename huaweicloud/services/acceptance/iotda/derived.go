package iotda

import (
	"fmt"
	"strings"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func WithDerivedAuth() bool {
	endpoint := acceptance.HW_IOTDA_ACCESS_ADDRESS
	if endpoint == "" {
		return false
	}

	subStr := fmt.Sprintf(".iotda-app.%s.", acceptance.HW_REGION_NAME)
	if index := strings.Index(endpoint, subStr); index > 0 {
		return true
	}

	return false
}

// When accessing an IoTDA standard or enterprise edition instance, you need to specify
// the IoTDA service endpoint in `provider` block.
// You can login to the IoTDA console, choose the instance Overview and click Access Details
// to view the HTTPS application access address. An example of the access address might be
// "9bc34xxxxx.st1.iotda-app.ap-southeast-1.myhuaweicloud.com".
func buildIoTDAEndpoint() string {
	endpoint := acceptance.HW_IOTDA_ACCESS_ADDRESS
	if endpoint == "" {
		return ""
	}

	// lintignore:AT004
	return fmt.Sprintf(`
provider "huaweicloud" {
  endpoints = {
    iotda = "%s"
  }
}
`, endpoint)
}
