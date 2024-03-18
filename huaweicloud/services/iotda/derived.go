package iotda

import (
	"fmt"
	"log"
	"strings"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// WithDerivedAuth calculate whether derived authentication is required by the endpoint.
// currently, this method only applies for HuaweiCloud.
// A sample endpoint: https://9bc34xxxxx.st1.iotda-app.ap-southeast-1.myhuaweicloud.com
func WithDerivedAuth(cfg *config.Config, region string) bool {
	endpoint := config.GetServiceEndpoint(cfg, "iotda", region)
	if endpoint == "" {
		log.Printf("[WARN ]failed to get the endpoint of IoTDA service in region %s", region)
		return false
	}

	subStr := fmt.Sprintf(".iotda-app.%s.", region)
	if index := strings.Index(endpoint, subStr); index > 0 {
		return true
	}

	return false
}
