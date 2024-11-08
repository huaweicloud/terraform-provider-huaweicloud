package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRdsRestoredDatabases_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rds_restored_databases.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRdsInstanceId(t)
			acceptance.TestAccPreCheckRdsInstanceName(t)
			acceptance.TestAccPreCheckRdsDatabaseName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRdsRestoredDatabases_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "database_limit"),
					resource.TestCheckResourceAttrSet(dataSource, "table_limit"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.total_tables"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.databases.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.databases.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.databases.0.total_tables"),

					resource.TestCheckOutput("instance_name_like_filter_is_useful", "true"),
					resource.TestCheckOutput("database_name_like_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceRdsRestoredDatabases_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_rds_restore_time_ranges" "test" {
  instance_id = "%[1]s"
}

data "huaweicloud_rds_restored_databases" "test" {
  engine       = "postgresql"
  instance_ids = ["%[1]s"]
  restore_time = data.huaweicloud_rds_restore_time_ranges.test.restore_time[0].end_time
}

locals {
  instance_name_like = "%[2]s"
}
data "huaweicloud_rds_restored_databases" "instance_name_like_filter" {
  engine             = "postgresql"
  instance_ids       = ["%[1]s"]
  restore_time       = data.huaweicloud_rds_restore_time_ranges.test.restore_time[0].end_time
  instance_name_like = "%[2]s"
}
output "instance_name_like_filter_is_useful" {
  value = length(data.huaweicloud_rds_restored_databases.instance_name_like_filter.instances) > 0
}

locals {
  database_name_like = "%[3]s"
}
data "huaweicloud_rds_restored_databases" "database_name_like_filter" {
  engine             = "postgresql"
  instance_ids       = ["%[1]s"]
  restore_time       = data.huaweicloud_rds_restore_time_ranges.test.restore_time[0].end_time
  database_name_like = "%[3]s"
}
output "database_name_like_filter_is_useful" {
  value = length(data.huaweicloud_rds_restored_databases.database_name_like_filter.instances) > 0 && alltrue(
  [for v in data.huaweicloud_rds_restored_databases.database_name_like_filter.instances[*].databases : length(v) > 0]
  )
}
`, acceptance.HW_RDS_INSTANCE_ID, acceptance.HW_RDS_INSTANCE_NAME, acceptance.HW_RDS_DATABASE_NAME)
}
