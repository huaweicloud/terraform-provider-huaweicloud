package hss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccPolicyGroupDeploy_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please prepare a valid HSS target policy ID and config it to the environment variable.
			acceptance.TestAccPreCheckHSSTargetPolicyGroupId(t)
			// Please prepare a valid HSS default target policy ID and config it to the environment variable.
			acceptance.TestAccPreCheckHSSDefaultTargetPolicyGroupId(t)
			// Please prepare a valid HSS host protection host id and config it to the environment variable.
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testPolicyGroupDeploy_basic(),
			},
		},
	})
}

func testPolicyGroupDeploy_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_hss_policy_group_deploy" "test" {
  target_policy_group_id = "%s"
  enterprise_project_id  = "all_granted_eps"
  operate_all            = true
}

resource "huaweicloud_hss_policy_group_deploy" "test1" {
  target_policy_group_id = "%s"
  operate_all            = false

  host_id_list = ["%s"]
}
`, acceptance.HW_HSS_TARGET_POLICY_GROUP_ID, acceptance.HW_HSS_DEFAULT_TARGET_POLICY_GROUP_ID, acceptance.HW_HSS_HOST_PROTECTION_HOST_ID)
}
