package taurusdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceTaurusDBHtapPrimaryInstanceDatabases_basic(t *testing.T) {
	dataSource := "data.huaweicloud_taurusdb_htap_primary_instance_databases.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckTaurusDBHtapInstanceId(t)
			acceptance.TestAccPreCheckTaurusDBInstanceId(t)
			acceptance.TestAccPreCheckTaurusDBDatabaseName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceTaurusDBHtapPrimaryInstanceDatabases_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "database_names.#"),
					resource.TestCheckOutput("dbname_fuzzy_query_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceTaurusDBHtapPrimaryInstanceDatabases_basic() string {
	return fmt.Sprintf(`
locals {
  databse_name  = "%[3]s"
  fuzzy_pattern = substr(local.databse_name, 0, length(local.databse_name) - 1)
}

data "huaweicloud_taurusdb_htap_primary_instance_databases" "test" {
  instance_id        = "%[1]s"
  source_instance_id = "%[2]s"
  databases          = [local.fuzzy_pattern]
}

locals {
  database_names = data.huaweicloud_taurusdb_htap_primary_instance_databases.test.database_names
}

output "dbname_fuzzy_query_is_useful" {
  value = length(local.database_names) > 0 && alltrue(
    [for dbname in local.database_names : strcontains(dbname, local.fuzzy_pattern)]
  )
}
`, acceptance.HW_TAURUSDB_HTAP_INSTANCE_ID, acceptance.HW_TAURUSDB_INSTANCE_ID, acceptance.HW_TAURUSDB_DATABASE_NAME)
}
