package drs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccUpdateDataProgressRules_basic(t *testing.T) {
	resourceName := "huaweicloud_drs_update_data_progress_rules.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDrsJobId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testUpdateDataProgressRules_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
		},
	})
}

// Currently, the update API only succeeds when the boyd parameter is not passed.
// Otherwise, the query API will report an error instead of a task failure
func testUpdateDataProgressRules_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_drs_update_data_progress_rules" "test" {
  job_id = "%s"
}
`, acceptance.HW_DRS_JOB_ID)
}
