package hss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccAntivirusCreatePayPerScanTask_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case requires setting a host ID with host protection enabled,
			// and the host is under the default enterprise project.
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAntivirusCreatePayPerScanTask_basic(),
			},
		},
	})
}

func testAntivirusCreatePayPerScanTask_basic() string {
	taskName := acceptance.RandomAccResourceName()

	return fmt.Sprintf(`
resource "huaweicloud_hss_antivirus_create_pay_per_scan_task" "test" {
  task_name             = "%[1]s"
  scan_type             = "quick"
  action                = "auto"
  host_ids              = ["%[2]s"]
  file_types            = [0]
  enterprise_project_id = "0"
}
`, taskName, acceptance.HW_HSS_HOST_PROTECTION_HOST_ID)
}
