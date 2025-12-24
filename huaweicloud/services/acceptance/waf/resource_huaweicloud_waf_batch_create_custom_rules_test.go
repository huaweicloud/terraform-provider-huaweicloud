package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccBatchCreateCustomRules_basic(t *testing.T) {
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
				Config: testAccBatchCreateCustomRules_basic(),
			},
		},
	})
}

func testAccBatchCreateCustomRules_basic() string {
	name := acceptance.RandomAccResourceName()
	return fmt.Sprintf(`
resource "huaweicloud_waf_batch_create_custom_rules" "test" {
  name        = "%[1]s"
  description = "test description"
  time        = false
  priority    = 10
  policy_ids  = ["%[2]s"]
  
  conditions {
    category        = "url"
    logic_operation = "equal"
    contents        = ["/[1]s"]
  }

  action {
    category = "block"
  }
}
`, name, acceptance.HW_WAF_POLICY_ID)
}
