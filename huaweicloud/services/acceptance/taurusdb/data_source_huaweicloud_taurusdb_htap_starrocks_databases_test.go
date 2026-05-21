package taurusdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceTaurusDBHtapStarrocksDatabases_basic(t *testing.T) {
	dataSource := "data.huaweicloud_taurusdb_htap_starrocks_databases.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckTaurusDBHtapInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceTaurusDBHtapStarrocksDatabases_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "databases.#"),
					resource.TestCheckOutput("filter_dbname_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceTaurusDBHtapStarrocksDatabases_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_taurusdb_htap_starrocks_databases" "test" {
  instance_id = "%[1]s"
}

data "huaweicloud_taurusdb_htap_starrocks_databases" "filter_by_dbname" {
  instance_id   = "%[1]s"
  database_name = "sys"

  depends_on = [data.huaweicloud_taurusdb_htap_starrocks_databases.test]
}

locals {
  filtered_databases_length = length(data.huaweicloud_taurusdb_htap_starrocks_databases.filter_by_dbname.databases)
  filtered_database_name    = data.huaweicloud_taurusdb_htap_starrocks_databases.filter_by_dbname.databases[0]
}

output "filter_dbname_is_useful" {
  value = local.filtered_databases_length == 1 && local.filtered_database_name == "sys"
}
`, acceptance.HW_TAURUSDB_HTAP_INSTANCE_ID)
}
