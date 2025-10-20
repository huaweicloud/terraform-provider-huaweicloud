package hss

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccHSSAppWhitelistPolicyProcess_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAppWhitelistPolicyProcess_basic,
				ExpectError: regexp.MustCompile(`Access denied.`),
			},
		},
	})
}

// The values ​​of policy_id and process_hash_list are fake data.
const testAppWhitelistPolicyProcess_basic = `
resource "huaweicloud_hss_app_whitelist_policy_process" "test" {
  policy_id      = "3bd2a82c-4b37-47f3-952b-fa323c22c8e6"
  process_status = "unknown"

  process_hash_list = [
    "3bd2a82c-4b37-47f3-952b-fa323c22c8e6",
  ]
}
`
