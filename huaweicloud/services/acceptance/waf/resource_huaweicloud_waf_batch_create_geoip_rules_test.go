package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccBatchCreateGeoipRules_basic(t *testing.T) {
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
				Config: testAccBatchCreateGeoipRules_basic(),
			},
		},
	})
}

func testAccBatchCreateGeoipRules_basic() string {
	name := acceptance.RandomAccResourceName()
	return fmt.Sprintf(`
resource "huaweicloud_waf_batch_create_geoip_rules" "test" {
  name                  = "%[1]s"
  geoip                 = "CA"
  white                 = 1
  ip_type               = "v4"
  policy_ids            = ["%[2]s"]
  enterprise_project_id = "%[3]s"
  description           = "test description"
}
`, name, acceptance.HW_WAF_POLICY_ID, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
