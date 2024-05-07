package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRdsErrorLogs_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rds_error_logs.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRdsInstanceId(t)
			acceptance.TestAccPreCheckRdsTimeRange(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRdsErrorLogs_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "error_logs.#"),
					resource.TestCheckResourceAttrSet(dataSource, "error_logs.0.time"),
					resource.TestCheckResourceAttrSet(dataSource, "error_logs.0.level"),
					resource.TestCheckResourceAttrSet(dataSource, "error_logs.0.content"),

					resource.TestCheckOutput("level_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceRdsErrorLogs_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_rds_error_logs" "test" {
  instance_id = "%[1]s"
  start_time  = "%[2]s"
  end_time    = "%[3]s"
}

locals {
  level = "WARNING"
}
data "huaweicloud_rds_error_logs" "level_filter" {
  instance_id = "%[1]s"
  start_time  = "%[2]s"
  end_time    = "%[3]s"
  level       = "WARNING"
}
output "level_filter_is_useful" {
  value = length(data.huaweicloud_rds_error_logs.level_filter.error_logs) > 0 && alltrue(
  [for v in data.huaweicloud_rds_error_logs.level_filter.error_logs[*].level : v == local.level]
  )
}
`, acceptance.HW_RDS_INSTANCE_ID, acceptance.HW_RDS_START_TIME, acceptance.HW_RDS_END_TIME)
}
