package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccBatchCreateIpReputationRules_basic(t *testing.T) {
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
				Config: testAccBatchCreateIpReputationRules_basic(),
			},
		},
	})
}

func testAccBatchCreateIpReputationRules_basic() string {
	name := acceptance.RandomAccResourceName()
	return fmt.Sprintf(`
resource "huaweicloud_waf_batch_create_ip_reputation_rules" "test" {
  name                  = "%[1]s"
  type                  = "idc"
  tags                  = ["Tencent"]
  policy_ids            = ["%[2]s"]
  enterprise_project_id = "%[3]s"
  description           = "test_ip_reputation_rule"

  action {
    category = "block"
  }
}
`, name, acceptance.HW_WAF_POLICY_ID, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
