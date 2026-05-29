package drs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceCancelCompareJob_basic(t *testing.T) {
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
				Config: testAccResourceCancelCompareJob_basic(),
			},
		},
	})
}

func testAccResourceCancelCompareJob_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_drs_compare_job_cancel" "test" {
  job_id         = "%s"
  compare_job_id = "%s"
}
`, acceptance.HW_DRS_JOB_ID, acceptance.HW_DRS_COMPARE_JOB_ID)
}
