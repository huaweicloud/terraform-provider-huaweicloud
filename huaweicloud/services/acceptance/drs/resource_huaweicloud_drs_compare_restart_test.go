package drs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDrsCompareRestart_basic(t *testing.T) {
	resourceName := "huaweicloud_drs_compare_restart.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDrsJobId(t)
			acceptance.TestAccPreCheckDrsCompareJobId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDrsCompareRestart_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "job_id", acceptance.HW_DRS_JOB_ID),
				),
			},
		},
	})
}

func testAccDrsCompareRestart_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_drs_compare_restart" "test" {
  job_id           = "%[1]s"
  compare_job_ids  = ["%[2]s"]
}
`, acceptance.HW_DRS_JOB_ID, acceptance.HW_DRS_COMPARE_JOB_ID)
}
