package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceTableDefinition_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_table_definition.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGaussDBInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceTableDefinition_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "table_definitions.#"),
					resource.TestCheckResourceAttrSet(dataSource, "table_definitions.0.table_definition"),
					resource.TestCheckResourceAttrSet(dataSource, "table_definitions.0.schema_name"),
					resource.TestCheckOutput("schema_name_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceTableDefinition_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_gaussdb_table_definition" "test" {
  instance_id   = "%[1]s"
  database_name = "test_db"
  table_name    = "test_table"
}

data "huaweicloud_gaussdb_table_definition" "schema_name_filter" {
  instance_id   = "%[1]s"
  database_name = "test_db"
  table_name    = "test_table"
  schema_name   = "test_schema"
}
output "schema_name_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_table_definition.schema_name_filter.table_definitions) > 0 && alltrue(
  [for v in data.huaweicloud_gaussdb_table_definition.schema_name_filter.table_definitions[*].schema_name : v == "test_schema"]
  )
}
`, acceptance.HW_GAUSSDB_INSTANCE_ID)
}
