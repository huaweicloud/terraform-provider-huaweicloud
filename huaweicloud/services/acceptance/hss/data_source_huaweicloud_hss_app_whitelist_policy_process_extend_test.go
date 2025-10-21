package hss

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAppWhitelistPolicyProcessExtend_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceAppWhitelistPolicyProcessExtend_basic(),
				ExpectError: regexp.MustCompile(`Access denied.`),
			},
		},
	})
}

// The values ​​of policy_id and host_id are fake data.
func testAccDataSourceAppWhitelistPolicyProcessExtend_basic() string {
	return `
data "huaweicloud_hss_app_whitelist_policy_process_extend" "test" {
  policy_id = "3bd2a82c-4b37-47f3-952b-fa323c22c8e6"
  host_id   = "3bd2a82c-4b37-47f3-952b-fa323c22c8e6"
}
`
}
