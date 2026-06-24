package secmaster

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceCreateRetryPolicy_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config:      testAccCreateRetryPolicy_basic(),
				ExpectError: regexp.MustCompile(`该操作连接当前流程不支持`),
			},
		},
	})
}

// The field values of `account_scope`, `eps_scope`, `region_scope` and `defense_connection_id` are fake data for testing
// a failure scenario.
func testAccCreateRetryPolicy_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_secmaster_create_retry_policy" "test" {
  workspace_id     = "%[1]s"
  action_type      = "create"
  version          = "25.8.0"
  block_target     = "192.168.0.0"
  policy_category  = "BLOCK"
  policy_direction = "INGRESS,EGRESS"
  account_scope    = "0970d7b7d400f2470fbec00316a03561"
  eps_scope        = "571f1c69-6d32-43d9-8ba5-eff73fc3eebf"
  region_scope     = "09d5805b3300f46e2f65c00326c7fcc0"

  block_age {
    is_block_ageing = false
  }

  defense_policy_list {
    defense_connection_id = "b3d704f3-cf7a-3aba-a8bd-2394b4a3d143"
  }

  policy_type {
    policy_type = "Source Ip"
  }
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID)
}
