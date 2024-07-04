package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRdsRestoreTimeRanges_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rds_restore_time_ranges.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRdsInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceRdsRestoreTimeRanges_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "restore_time.#"),
					resource.TestCheckResourceAttrSet(dataSource, "restore_time.0.start_time"),
					resource.TestCheckResourceAttrSet(dataSource, "restore_time.0.end_time"),

					resource.TestCheckOutput("date_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceRdsRestoreTimeRanges_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_rds_restore_time_ranges" "test" {
  instance_id = "%[1]s"
}

locals {
  date = "2024-05-20"
}
data "huaweicloud_rds_restore_time_ranges" "date_filter" {
  instance_id = "%[1]s"
  date        = local.date
}
output "date_filter_is_useful" {
  value = length(data.huaweicloud_rds_restore_time_ranges.date_filter.restore_time) > 0
}
`, acceptance.HW_RDS_INSTANCE_ID)
}
