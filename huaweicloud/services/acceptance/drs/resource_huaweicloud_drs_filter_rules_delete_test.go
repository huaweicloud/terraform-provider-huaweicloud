package drs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDrsFilterRulesDelete_basic(t *testing.T) {
	resourceName := "huaweicloud_drs_filter_rules_delete.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDrsJobId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDrsFilterRulesDelete_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "job_id", acceptance.HW_DRS_JOB_ID),
					resource.TestCheckResourceAttr(resourceName, "type", "common"),
				),
			},
		},
	})
}

func testAccDrsFilterRulesDelete_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_drs_filter_rules_delete" "test" {
  job_id   = "%s"
  rule_ids = [1, 2]
  type     = "common"
}
`, acceptance.HW_DRS_JOB_ID)
}
