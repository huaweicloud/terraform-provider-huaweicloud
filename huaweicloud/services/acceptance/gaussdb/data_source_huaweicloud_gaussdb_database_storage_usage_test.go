package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDatabaseStorageUsage_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_database_storage_usage.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGaussDBInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDatabaseStorageUsage_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "database_volumes.#"),
					resource.TestCheckResourceAttrSet(dataSource, "database_volumes.0.database_name"),
					resource.TestCheckResourceAttrSet(dataSource, "database_volumes.0.table_space_name"),
					resource.TestCheckResourceAttrSet(dataSource, "database_volumes.0.user_name"),
					resource.TestCheckResourceAttrSet(dataSource, "database_volumes.0.database_size"),
					resource.TestCheckOutput("database_name_filter_is_useful", "true"),
					resource.TestCheckOutput("table_space_name_filter_is_useful", "true"),
					resource.TestCheckOutput("user_name_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDatabaseStorageUsage_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_gaussdb_database_storage_usage" "test" {
  instance_id = "%[1]s"
}

data "huaweicloud_gaussdb_database_storage_usage" "database_name_filter" {
  instance_id   = "%[1]s"
  database_name = data.huaweicloud_gaussdb_database_storage_usage.test.database_volumes[0].database_name
}
locals {
  database_name = data.huaweicloud_gaussdb_database_storage_usage.test.database_volumes[0].database_name
}
output "database_name_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_database_storage_usage.database_name_filter) > 0 && alltrue(
    [for v in data.huaweicloud_gaussdb_database_storage_usage.database_name_filter[*].database_name : v == local.database_name]
  )
}

data "huaweicloud_gaussdb_database_storage_usage" "table_space_name_filter" {
  instance_id      = "%[1]s"
  table_space_name = data.huaweicloud_gaussdb_database_storage_usage.test.database_volumes[0].table_space_name
}
locals {
  table_space_name = data.huaweicloud_gaussdb_database_storage_usage.test.database_volumes[0].table_space_name
}
output "table_space_name_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_database_storage_usage.table_space_name_filter) > 0 && alltrue(
    [for v in data.huaweicloud_gaussdb_database_storage_usage.table_space_name_filter[*].table_space_name :
    v == local.table_space_name]
  )
}

data "huaweicloud_gaussdb_database_storage_usage" "user_name_filter" {
  instance_id = "%[1]s"
  user_name   = data.huaweicloud_gaussdb_database_storage_usage.test.database_volumes[0].user_name
}
locals {
  user_name = data.huaweicloud_gaussdb_database_storage_usage.test.database_volumes[0].user_name
}
output "user_name_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_database_storage_usage.user_name_filter) > 0 && alltrue(
    [for v in data.huaweicloud_gaussdb_database_storage_usage.user_name_filter[*].user_name : v == local.user_name]
  )
}
`, acceptance.HW_GAUSSDB_INSTANCE_ID)
}
