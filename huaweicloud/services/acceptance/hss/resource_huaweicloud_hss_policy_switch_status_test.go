package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccPolicySwitchStatus_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testPolicySwitchStatus_basic,
			},
			{
				Config: testPolicySwitchStatus_basic_update,
			},
		},
	})
}

const testPolicySwitchStatus_basic = `
resource "huaweicloud_hss_policy_switch_status" "test" {
  enterprise_project_id = "all_granted_eps"
  policy_name           = "sp_feature"
  enable                = true
}
`

const testPolicySwitchStatus_basic_update = `
resource "huaweicloud_hss_policy_switch_status" "test" {
  enterprise_project_id = "all_granted_eps"
  policy_name           = "sp_feature"
  enable                = false
}
`
