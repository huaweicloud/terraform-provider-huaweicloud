package antiddos

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Note: Due to limited test conditions, this test case only verifies the expected error scenario.
func TestAccResourceOpenProtection_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testOpenProtection_basic,
				ExpectError: regexp.MustCompile(`VPC access failed or EIP is not exist|访问VPC平台异常或EIP不存在`),
			},
		},
	})
}

// The value of field `floating_ip_id` is mock data.
const testOpenProtection_basic = `
resource "huaweicloud_antiddos_open_protection" "test" {
  floating_ip_id         = "85f523db-9d00-42af-abf6-24c64738e7a2"
  app_type_id            = 0
  cleaning_access_pos_id = 99
  enable_l7              = false
  http_request_pos_id    = 1
  traffic_pos_id         = 7
  antiddos_config_id     = "1234567890"
}
`
