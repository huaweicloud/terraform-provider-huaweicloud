package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccBatchCreateAntitamperRules_basic(t *testing.T) {
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
				Config: testAccBatchCreateAntitamperRules_basic(),
			},
		},
	})
}

func testAccBatchCreateAntitamperRules_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_waf_batch_create_antitamper_rules" "test" {
  hostname              = "www.test.com"
  url                   = "/test"
  policy_ids            = ["%[1]s"]
  enterprise_project_id = "%[2]s"
  description           = "test description"
}
`, acceptance.HW_WAF_POLICY_ID, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
