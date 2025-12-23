package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccBatchCreateWhiteBlackIpRules_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckWafInstance(t)
			acceptance.TestAccPreCheckWafPolicyId(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccBatchCreateWhiteBlackIpRules_basic(name),
			},
		},
	})
}

func testAccBatchCreateWhiteBlackIpRules_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_waf_batch_create_whiteblackip_rules" "test" {
  name                  = "%[1]s"
  white                 = 1
  policy_ids            = ["%[2]s"]
  addr                  = "127.0.0.1"
  description           = "test description"
  time_mode             = "permanent"
  enterprise_project_id = "%[3]s"
}
`, name, acceptance.HW_WAF_POLICY_ID, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
