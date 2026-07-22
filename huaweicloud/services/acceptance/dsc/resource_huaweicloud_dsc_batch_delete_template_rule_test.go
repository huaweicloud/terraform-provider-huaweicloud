package dsc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceBatchDeleteTemplateRule_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDSCScanTemplateID(t)
			acceptance.TestAccPreCheckDscRuleIds(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceBatchDeleteTemplateRule_basic(),
			},
		},
	})
}

func testAccResourceBatchDeleteTemplateRule_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_dsc_batch_delete_template_rule" "test" {
  template_id = "%s"
  rule_ids    = "%s"
}
`, acceptance.HW_DSC_SCAN_TEMPLATE_ID, acceptance.HW_DSC_RULE_IDS)
}
