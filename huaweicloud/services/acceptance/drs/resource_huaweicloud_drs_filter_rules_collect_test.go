package drs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDrsFilterRulesCollect_basic(t *testing.T) {
	resourceName := "huaweicloud_drs_filter_rules_collect.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDrsJobId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDrsFilterRulesCollect_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "job_id", acceptance.HW_DRS_JOB_ID),
					resource.TestCheckResourceAttr(resourceName, "filter_type", "common"),
				),
			},
		},
	})
}

func testAccDrsFilterRulesCollect_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_drs_filter_rules_collect" "test" {
  job_id       = "%s"
  filter_type  = "common"
  is_all_rules = true
}
`, acceptance.HW_DRS_JOB_ID)
}
