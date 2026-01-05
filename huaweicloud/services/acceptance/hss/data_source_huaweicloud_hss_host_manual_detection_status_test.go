package hss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceHostManualDetectionStatus_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_hss_host_manual_detection_status.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case needs to ensure the existence of a host with host protection enabled,
			// and the host is under the default enterprise project.
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceHostManualDetectionStatus_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "scan_status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "scaned_time"),
				),
			},
		},
	})
}

func testDataSourceHostManualDetectionStatus_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_hss_host_manual_detection_status" "test" {
  host_id               = "%s"
  type                  = "pwd"
  enterprise_project_id = "0"
}
`, acceptance.HW_HSS_HOST_PROTECTION_HOST_ID)
}
