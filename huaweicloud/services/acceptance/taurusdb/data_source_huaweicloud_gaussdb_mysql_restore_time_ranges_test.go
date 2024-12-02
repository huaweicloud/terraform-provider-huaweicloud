package taurusdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGaussdbMysqlRestoreTimeRanges_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_mysql_restore_time_ranges.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceGaussdbMysqlRestoreTimeRanges_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "restore_times.#"),
					resource.TestCheckResourceAttrSet(dataSource, "restore_times.0.start_time"),
					resource.TestCheckResourceAttrSet(dataSource, "restore_times.0.end_time"),

					resource.TestCheckOutput("date_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceGaussdbMysqlRestoreTimeRanges_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_gaussdb_mysql_restore_time_ranges" "test" {
  instance_id = huaweicloud_gaussdb_mysql_instance.test.id
}

locals {
  date = split("T", huaweicloud_gaussdb_mysql_backup.test.end_time)[0]
}
data "huaweicloud_gaussdb_mysql_restore_time_ranges" "date_filter" {
  instance_id = huaweicloud_gaussdb_mysql_instance.test.id
  date        = local.date
}
output "date_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_mysql_restore_time_ranges.date_filter.restore_times) > 0
}
`, testBackup_basic(name))
}
