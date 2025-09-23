package dbss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDbssDatabases_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_dbss_databases.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)

		byStatus   = "data.huaweicloud_dbss_databases.filter_by_status"
		dcByStatus = acceptance.InitDataSourceCheck(byStatus)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Before running this test case, please prepare an audit instance associated with database.
			acceptance.TestAccPrecheckDbssInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDbssDatabases_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "databases.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "databases.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "databases.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "databases.0.type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "databases.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "databases.0.ip"),
					resource.TestCheckResourceAttrSet(dataSourceName, "databases.0.os"),
					resource.TestCheckResourceAttrSet(dataSourceName, "databases.0.port"),
					resource.TestCheckResourceAttrSet(dataSourceName, "databases.0.version"),
					resource.TestCheckResourceAttrSet(dataSourceName, "databases.0.charset"),
					resource.TestCheckResourceAttrSet(dataSourceName, "databases.0.instance_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "databases.0.db_classification"),

					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("status_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDbssDatabases_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dbss_databases" "test" {
  instance_id = "%[1]s"
}

locals {
  status = data.huaweicloud_dbss_databases.test.databases[0].status
}

data "huaweicloud_dbss_databases" "filter_by_status" {
  instance_id = "%[1]s"
  status      = local.status
}

output "status_filter_useful" {
  value = length(data.huaweicloud_dbss_databases.filter_by_status.databases) > 0 && alltrue(
    [for v in data.huaweicloud_dbss_databases.filter_by_status.databases[*].status : v == local.status]
  )
}
`, acceptance.HW_DBSS_INSATNCE_ID)
}
