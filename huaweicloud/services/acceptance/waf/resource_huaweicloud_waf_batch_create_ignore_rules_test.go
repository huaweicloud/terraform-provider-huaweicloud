package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccBatchCreateIgnoreRules_basic(t *testing.T) {
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
				Config: testAccBatchCreateIgnoreRules_basic(),
			},
		},
	})
}

func testAccBatchCreateIgnoreRules_basic() string {
	domainName := fmt.Sprintf("%s.huaweicloud.com", acceptance.RandomAccResourceName())
	return fmt.Sprintf(`
resource "huaweicloud_waf_batch_create_ignore_rules" "test" {
  rule                  = "xss"
  description           = "test description"
  policy_ids            = ["%[1]s"]
  enterprise_project_id = "%[2]s"
  domain                = ["%[3]s"]

  conditions {
    category        = "url"
    contents        = ["/admin"]
    logic_operation = "equal"
  }

  advanced {
    index    = "params"
    contents = ["test-param"]
  }
}
`, acceptance.HW_WAF_POLICY_ID, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST, domainName)
}
