package taurusdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceTaurusDBHtapPrimaryInstanceTables_basic(t *testing.T) {
	dataSource1 := "data.huaweicloud_taurusdb_htap_primary_instance_tables.whitelist"
	dataSource2 := "data.huaweicloud_taurusdb_htap_primary_instance_tables.blacklist"
	dc1 := acceptance.InitDataSourceCheck(dataSource1)
	dc2 := acceptance.InitDataSourceCheck(dataSource2)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckTaurusDBHtapInstanceId(t)
			acceptance.TestAccPreCheckTaurusDBInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceTaurusDBHtapPrimaryInstanceTables_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc1.CheckResourceExists(),
					dc2.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSource1, "tables.#", "1"),
					resource.TestCheckResourceAttr(dataSource1, "tables.0", "db.table1"),
					resource.TestCheckResourceAttr(dataSource2, "tables.#", "1"),
					resource.TestCheckResourceAttr(dataSource2, "tables.0", "db.table3"),
				),
			},
		},
	})
}

func testDataSourceTaurusDBHtapPrimaryInstanceTables_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_taurusdb_htap_primary_instance_tables" "whitelist" {
  instance_id        = "%[1]s"
  source_instance_id = "%[2]s"
  filter_type        = "include_tables"

  database_tables {
    database = "db"
    tables   = ["table1"]
  }

  selected_tables {
    database = "db"
    tables   = ["table1", "table2"]
  }
}

data "huaweicloud_taurusdb_htap_primary_instance_tables" "blacklist" {
  instance_id        = "%[1]s"
  source_instance_id = "%[2]s"
  filter_type        = "exclude_tables"

  database_tables {
    database = "db"
    tables   = ["table3"]
  }

  selected_tables {
    database = "db"
    tables   = ["table1", "table2"]
  }
}
`, acceptance.HW_TAURUSDB_HTAP_INSTANCE_ID, acceptance.HW_TAURUSDB_INSTANCE_ID)
}
