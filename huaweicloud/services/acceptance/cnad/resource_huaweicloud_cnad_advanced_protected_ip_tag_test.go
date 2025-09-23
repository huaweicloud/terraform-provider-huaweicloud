package cnad

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Note: Due to limited test conditions, this test case only verifies the expected error scenario.
func TestAccResourceProtectedIpTag_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testProtectedIpTag_basic,
				ExpectError: regexp.MustCompile(`Resource Not Found.|资源不存在。`),
			},
		},
	})
}

// The value of field `protected_ip_id` is mock data.
const testProtectedIpTag_basic = `
resource "huaweicloud_cnad_advanced_protected_ip_tag" "test" {
  protected_ip_id = "79189f77-bc57-46ff-a69d-17168d95c970"
  tag             = "test"
}`
