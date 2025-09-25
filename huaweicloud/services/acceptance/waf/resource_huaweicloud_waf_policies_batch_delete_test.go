package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccPoliciesBatchDelete_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case requires the preparation of a policy ID under the default enterprise project.
			acceptance.TestAccPreCheckWafPolicyId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPoliciesBatchDelete_basic(),
			},
		},
	})
}

func testAccPoliciesBatchDelete_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_waf_policies_batch_delete" "test" {
  policy_ids            = ["%[1]s"]
  enterprise_project_id = "0"
}
`, acceptance.HW_WAF_POLICY_ID)
}
