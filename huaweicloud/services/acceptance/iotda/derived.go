package iotda

import (
	"fmt"
	"strings"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func withDerivedAuth() bool {
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
