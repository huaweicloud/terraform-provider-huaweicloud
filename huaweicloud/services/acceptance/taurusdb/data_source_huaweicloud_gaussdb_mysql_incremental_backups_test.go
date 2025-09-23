package taurusdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGaussdbMysqlIncrementalBackups_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_mysql_incremental_backups.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceGaussdbMysqlIncrementalBackups_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "backups.#"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.begin_time"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.end_time"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.size"),

					resource.TestCheckOutput("backup_time_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceGaussdbMysqlIncrementalBackups_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_gaussdb_mysql_incremental_backups" "test" {
  instance_id = "%[1]s"
}

data "huaweicloud_gaussdb_mysql_incremental_backups" "backup_time_filter" {
  instance_id = "%[1]s"
  begin_time  = "%[2]s"
  end_time    = "%[3]s"
}
output "backup_time_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_mysql_incremental_backups.backup_time_filter.backups) > 0
}
`, acceptance.HW_GAUSSDB_MYSQL_INSTANCE_ID, acceptance.HW_GAUSSDB_MYSQL_BACKUP_BEGIN_TIME, acceptance.HW_GAUSSDB_MYSQL_BACKUP_END_TIME)
}
