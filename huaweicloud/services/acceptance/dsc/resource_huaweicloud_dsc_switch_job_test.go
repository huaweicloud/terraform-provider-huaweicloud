package dsc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceDscSwitchJob_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDscScanJobId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDscSwitchJob_basic(),
			},
		},
	})
}

func testAccResourceDscSwitchJob_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_dsc_switch_job" "test" {
  job_id = "%[1]s"
  open   = true
}
`, acceptance.HW_DSC_SCAN_JOB_ID)
}
