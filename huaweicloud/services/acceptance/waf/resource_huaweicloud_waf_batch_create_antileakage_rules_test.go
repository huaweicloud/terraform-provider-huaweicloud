package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccBatchCreateAntileakageRules_basic(t *testing.T) {
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
				Config: testAccBatchCreateAntileakageRules_basic(),
			},
		},
	})
}

func testAccBatchCreateAntileakageRules_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_waf_batch_create_antileakage_rules" "test" {
  url                   = "/test"
  category              = "sensitive"
  contents              = ["id_card", "phone"]
  policy_ids            = ["%[1]s"]
  description           = "test description"
  enterprise_project_id = "%[2]s"

  action {
    category = "block"
  }
}
`, acceptance.HW_WAF_POLICY_ID, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
