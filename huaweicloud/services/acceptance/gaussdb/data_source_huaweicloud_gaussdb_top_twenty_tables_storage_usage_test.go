package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceTopTwentyTablesStorageUsage_basic(t *testing.T) {
	dataSource1 := "data.huaweicloud_gaussdb_top_twenty_tables_storage_usage.test1"
	dataSource2 := "data.huaweicloud_gaussdb_top_twenty_tables_storage_usage.test2"
	dc1 := acceptance.InitDataSourceCheck(dataSource1)
	dc2 := acceptance.InitDataSourceCheck(dataSource2)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGaussDBInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceTopTwentyTablesStorageUsage_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc1.CheckResourceExists(),
					dc2.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource1, "job_id"),
					resource.TestCheckResourceAttrSet(dataSource1, "node_id"),
					resource.TestCheckResourceAttrSet(dataSource2, "table_volumes.#"),
					resource.TestCheckResourceAttrSet(dataSource2, "table_volumes.0.id"),
					resource.TestCheckResourceAttrSet(dataSource2, "table_volumes.0.table_name"),
					resource.TestCheckResourceAttrSet(dataSource2, "table_volumes.0.table_owner"),
					resource.TestCheckResourceAttrSet(dataSource2, "table_volumes.0.database_name"),
					resource.TestCheckResourceAttrSet(dataSource2, "table_volumes.0.schema_name"),
					resource.TestCheckResourceAttrSet(dataSource2, "table_volumes.0.is_part_type"),
					resource.TestCheckResourceAttrSet(dataSource2, "table_volumes.0.is_hash_cluster_key"),
					resource.TestCheckResourceAttrSet(dataSource2, "table_volumes.0.tuples"),
					resource.TestCheckResourceAttrSet(dataSource2, "table_volumes.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSource2, "table_volumes.0.update_time"),
					resource.TestCheckResourceAttrSet(dataSource2, "table_volumes.0.table_size"),
					resource.TestCheckResourceAttrSet(dataSource2, "state"),
				),
			},
		},
	})
}

func testDataSourceTopTwentyTablesStorageUsage_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_gaussdb_top_twenty_tables_storage_usage" "test1" {
  instance_id = "%[1]s"
}

data "huaweicloud_gaussdb_top_twenty_tables_storage_usage" "test2" {
  instance_id = "%[1]s"
  job_id      = data.huaweicloud_gaussdb_top_twenty_tables_storage_usage.test1.job_id
  node_id     = data.huaweicloud_gaussdb_top_twenty_tables_storage_usage.test1.node_id
}
`, acceptance.HW_GAUSSDB_INSTANCE_ID)
}
