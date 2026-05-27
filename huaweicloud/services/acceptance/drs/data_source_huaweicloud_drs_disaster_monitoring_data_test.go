package drs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDrsDisasterMonitoringData_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_drs_disaster_monitoring_data.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDrsJobIds(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDrsDisasterMonitoringData_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.data_guard_monitor.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.data_guard_monitor.0.src_normal"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.data_guard_monitor.0.dst_normal"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.data_guard_monitor.0.sr_delay"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.data_guard_monitor.0.dst_delay"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.data_guard_monitor.0.src_rps"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.data_guard_monitor.0.dst_rps"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.data_guard_monitor.0.bandwidth"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.data_guard_monitor.0.cpu_used_percent"),
				),
			},
		},
	})
}

func testAccDataSourceDrsDisasterMonitoringData_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_drs_disaster_monitoring_data" "test" {
  job_ids = split(",", "%s")
}
`, acceptance.HW_DRS_JOB_IDS)
}
