package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceTablesStorageUsage_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_tables_storage_usage.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGaussDBInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceTablesStorageUsage_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "table_volumes.#"),
					resource.TestCheckResourceAttrSet(dataSource, "table_volumes.0.table_size"),
					resource.TestCheckResourceAttrSet(dataSource, "table_volumes.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "table_volumes.0.table_name"),
					resource.TestCheckResourceAttrSet(dataSource, "table_volumes.0.table_owner"),
					resource.TestCheckResourceAttrSet(dataSource, "table_volumes.0.schema_name"),
					resource.TestCheckResourceAttrSet(dataSource, "table_volumes.0.database_name"),
					resource.TestCheckResourceAttrSet(dataSource, "table_volumes.0.is_part_type"),
					resource.TestCheckResourceAttrSet(dataSource, "table_volumes.0.is_hash_cluster_key"),
					resource.TestCheckResourceAttrSet(dataSource, "table_volumes.0.tuples"),
					resource.TestCheckResourceAttrSet(dataSource, "table_volumes.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "table_volumes.0.update_time"),
					resource.TestCheckOutput("table_name_filter_is_useful", "true"),
					resource.TestCheckOutput("sort_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceTablesStorageUsage_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_gaussdb_tables_storage_usage" "test" {
  instance_id   = "%[1]s"
  database_name = "test_db"
  schema_names  = ["test_schema"]
}


data "huaweicloud_gaussdb_tables_storage_usage" "table_name_filter" {
  instance_id   = "%[1]s"
  database_name = "test_db"
  schema_names  = ["test_schema"]
  table_name    = "test_table"
}
output "table_name_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_tables_storage_usage.table_name_filter) > 0 && alltrue(
    [for v in data.huaweicloud_gaussdb_tables_storage_usage.table_name_filter[*].table_name : v == "test_table"]
  )
}

data "huaweicloud_gaussdb_tables_storage_usage" "sort_filter" {
  instance_id   = "%[1]s"
  database_name = "test_db"
  schema_names  = ["test_schema"]
  sort_key      = "id"
  sort_order    = "DESC"
}
output "sort_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_tables_storage_usage.sort_filter) > 0
}
`, acceptance.HW_GAUSSDB_INSTANCE_ID)
}
