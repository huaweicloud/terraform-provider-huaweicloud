package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGaussDBInstanceDatabaseTables_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_instance_database_tables.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGaussDBInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGaussDBInstanceDatabaseTables_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "tables.#"),
					resource.TestCheckResourceAttrSet(dataSource, "tables.0.table_name"),
				),
			},
		},
	})
}

func testAccDataSourceGaussDBInstanceDatabaseTables_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_gaussdb_instance_database_tables" "test" {
  instance_id = "%[1]s"
  db_name     = "gaussdb_test"
  schema_name = "rds"
}

data "huaweicloud_gaussdb_instance_database_tables" "name_filter" {
  instance_id        = "%[1]s"
  db_name            = "gaussdb_test"
  schema_name        = "rds"
  table_name_keyword = "test"
}
output "name_filter" {
  value = length(data.huaweicloud_gaussdb_instance_database_tables.name_filter.tables) > 0
}
`, acceptance.HW_GAUSSDB_INSTANCE_ID)
}
