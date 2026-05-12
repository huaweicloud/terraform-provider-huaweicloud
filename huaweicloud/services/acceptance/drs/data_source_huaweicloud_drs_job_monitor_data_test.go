package drs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDrsJobMonitorData_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_drs_job_monitor_data.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDrsJobId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDrsJobMonitorData_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "bandwidth"),
					resource.TestCheckResourceAttrSet(dataSourceName, "update_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "node_volume_size"),
				),
			},
		},
	})
}

func testAccDataSourceDrsJobMonitorData_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_drs_job_monitor_data" "test" {
  job_id = "%s"
}
`, acceptance.HW_DRS_JOB_ID)
}
