package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSingleSqlDetail_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_gaussdb_single_sql_details.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGaussDBInstanceId(t)
			acceptance.TestAccPreCheckGaussDBSqlId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSingleSqlDetail_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "components.#"),
					resource.TestCheckResourceAttrSet(dataSource, "components.0.component_id"),
					resource.TestCheckResourceAttrSet(dataSource, "components.0.db_name"),
					resource.TestCheckResourceAttrSet(dataSource, "components.0.username"),
					resource.TestCheckResourceAttrSet(dataSource, "components.0.sql_id"),
					resource.TestCheckResourceAttrSet(dataSource, "components.0.sql_exec_id"),
					resource.TestCheckResourceAttrSet(dataSource, "components.0.query"),
					resource.TestCheckResourceAttrSet(dataSource, "components.0.start_time"),
					resource.TestCheckResourceAttrSet(dataSource, "components.0.finish_time"),
					resource.TestCheckResourceAttrSet(dataSource, "components.0.is_slow_sql"),
					resource.TestCheckResourceAttrSet(dataSource, "components.0.finish_status"),
					resource.TestCheckResourceAttrSet(dataSource, "trace_statistics.#"),

					resource.TestCheckOutput("sql_id_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceSingleSqlDetail_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_gaussdb_single_sql_details" "test" {
  instance_id = "%[1]s"
  sql_exec_id = "%[2]s"
}

locals {
  sql_id = data.huaweicloud_gaussdb_single_sql_details.test.components[0].sql_id
}

data "huaweicloud_gaussdb_single_sql_details" "sql_id_filter" {
  instance_id = "%[1]s"
  sql_exec_id = "%[2]s"
  sql_id      = local.sql_id
}

output "sql_id_filter_useful" {
  value = length(data.huaweicloud_gaussdb_single_sql_details.sql_id_filter.components) > 0
}
`, acceptance.HW_GAUSSDB_INSTANCE_ID, acceptance.HW_GAUSSDB_SQL_ID)
}
