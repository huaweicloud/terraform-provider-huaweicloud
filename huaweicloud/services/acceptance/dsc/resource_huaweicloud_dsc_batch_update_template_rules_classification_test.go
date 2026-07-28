package dsc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceDscBatchUpdateTemplateRulesClassification_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDSCScanTemplateID(t)
			acceptance.TestAccPreCheckDscClassificationId(t)
			acceptance.TestAccPreCheckDscRuleIds(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDscBatchUpdateTemplateRulesClassification_basic(),
			},
		},
	})
}

func testAccResourceDscBatchUpdateTemplateRulesClassification_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_dsc_batch_update_template_rules_classification" "test" {
  template_id       = "%s"
  classification_id = "%s"
  rule_id_list      = split(",", "%s")
}
`, acceptance.HW_DSC_SCAN_TEMPLATE_ID, acceptance.HW_DSC_CLASSIFICATION_ID, acceptance.HW_DSC_RULE_IDS)
}
