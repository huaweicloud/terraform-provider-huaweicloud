package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSqlTextSchemaTable_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_sql_text_schema_table.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGaussDBInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSqlTextSchemaTable_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "database_tables.#"),
					resource.TestCheckResourceAttrSet(dataSource, "database_tables.0.table_name"),
					resource.TestCheckResourceAttrSet(dataSource, "database_tables.0.schema_name"),
				),
			},
		},
	})
}

func testDataSourceSqlTextSchemaTable_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_gaussdb_sql_text_schema_table" "test" {
  instance_id = "%s"
  sql_text    = "SELECT * FROM public.users u LEFT JOIN sales.orders o ON u.id = o.user_id;"
}
`, acceptance.HW_GAUSSDB_INSTANCE_ID)
}
