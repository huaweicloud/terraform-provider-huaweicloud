package dsc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDscScanTasks_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dsc_scan_tasks.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDscScanJobId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDscScanTasks_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.#"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.category"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.progress"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.asset_name"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.asset_id"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.start_time"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.end_time"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.scanned_object_num"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.to_be_scanned_object_num"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.scan_speed"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.skip_object_num"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.last_scan_risk"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.security_level_color"),
				),
			},
		},
	})
}

func testAccDataSourceDscScanTasks_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dsc_scan_tasks" "test" {
  job_id = "%[1]s"
}
`, acceptance.HW_DSC_SCAN_JOB_ID)
}
