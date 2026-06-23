package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSchemasStorageUsage_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_schemas_storage_usage.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGaussDBInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSchemasStorageUsage_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "schema_volumes.#"),
					resource.TestCheckResourceAttrSet(dataSource, "schema_volumes.0.schema_size"),
					resource.TestCheckResourceAttrSet(dataSource, "schema_volumes.0.table_count"),
					resource.TestCheckResourceAttrSet(dataSource, "schema_volumes.0.user_name"),
					resource.TestCheckResourceAttrSet(dataSource, "schema_volumes.0.schema_name"),
					resource.TestCheckOutput("schema_name_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceSchemasStorageUsage_basic() string {
	return fmt.Sprintf(`

data "huaweicloud_gaussdb_schemas_storage_usage" "test" {
  instance_id   = "%[1]s"
  database_name = "test_db"
}

data "huaweicloud_gaussdb_schemas_storage_usage" "schema_name_filter" {
  instance_id   = "%[1]s"
  database_name = "test_db"
  schema_name   = "test_schema"
}
output "schema_name_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_schemas_storage_usage.schema_name_filter) > 0 && alltrue(
    [for v in data.huaweicloud_gaussdb_schemas_storage_usage.schema_name_filter[*].schema_name : v == "test_schema"]
  )
}
`, acceptance.HW_GAUSSDB_INSTANCE_ID)
}
