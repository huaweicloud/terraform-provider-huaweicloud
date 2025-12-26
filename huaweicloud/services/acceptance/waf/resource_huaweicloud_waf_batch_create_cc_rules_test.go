package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccBatchCreateCcRules_basic(t *testing.T) {
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
				Config: testAccBatchCreateCcRules_basic(),
			},
		},
	})
}

func testAccBatchCreateCcRules_basic() string {
	name := acceptance.RandomAccResourceName()
	return fmt.Sprintf(`
resource "huaweicloud_waf_batch_create_cc_rules" "test" {
  name                  = "%[1]s"
  description           = "CC protection rule for admin paths"
  tag_type              = "ip"
  limit_num             = 100
  limit_period          = 60
  policy_ids            = ["%[2]s"]
  enterprise_project_id = "%[3]s"
  lock_time             = 2
  region_aggregation    = true
  domain_aggregation    = true
  
  conditions {
    category        = "url"
    logic_operation = "prefix"
    contents        = ["/admin"]
  }

  action {
    category = "block"
	detail {
	  response {
		content_type = "application/json"
		content      = "block application json"
	  }
	}
  }
}
`, name, acceptance.HW_WAF_POLICY_ID, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
