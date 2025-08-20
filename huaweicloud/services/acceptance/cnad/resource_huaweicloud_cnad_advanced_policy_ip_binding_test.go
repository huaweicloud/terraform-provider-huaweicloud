package cnad

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Note: Due to limited test conditions, this test case only verifies the expected error scenario.
func TestAccResourcePolicyIpBinding_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testPolicyIpBinding_basic,
				ExpectError: regexp.MustCompile(`Resource Not Found.|资源不存在。`),
			},
		},
	})
}

// The values of fields `policy_id` and `ip_list` are mock data.
const testPolicyIpBinding_basic = `
resource "huaweicloud_cnad_advanced_policy_ip_binding" "test" {
  policy_id = "1d8c03c4-a1d0-49cf-808a-ab50bfd7f0d8"
  ip_list   = ["192.168.1.1", "192.168.1.2"]
}`
