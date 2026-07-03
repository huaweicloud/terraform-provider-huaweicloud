package taurusdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceTaurusDBInstanceDbTableInfo_basic(t *testing.T) {
	var (
		dataSource  = "data.huaweicloud_taurusdb_instance_db_table_info.test"
		dataSource1 = "data.huaweicloud_taurusdb_instance_db_table_info.dbTest"
		dataSource2 = "data.huaweicloud_taurusdb_instance_db_table_info.tableTest"
		dc          = acceptance.InitDataSourceCheck(dataSource)
		dc1         = acceptance.InitDataSourceCheck(dataSource1)
		dc2         = acceptance.InitDataSourceCheck(dataSource2)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckTaurusDBInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTaurusDBInstanceDbTableInfo_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSource, "database_names.#", "1"),
					resource.TestCheckResourceAttr(dataSource, "database_names.0", "__taurus_sys__"),
					resource.TestCheckResourceAttr(dataSource, "table_names.#", "0"),
					resource.TestCheckResourceAttr(dataSource, "table_meta_infos.#", "0"),
				),
			},
			{
				Config: testAccDataSourceTaurusDBInstanceDbTableInfo_dbTest(),
				Check: resource.ComposeTestCheckFunc(
					dc1.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSource1, "database_names.#", "0"),
					resource.TestCheckResourceAttr(dataSource1, "table_names.#", "8"),
					resource.TestCheckResourceAttr(dataSource1, "table_meta_infos.#", "0"),
				),
			},
			{
				Config: testAccDataSourceTaurusDBInstanceDbTableInfo_tableTest(),
				Check: resource.ComposeTestCheckFunc(
					dc2.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSource2, "database_names.#", "0"),
					resource.TestCheckResourceAttr(dataSource2, "table_names.#", "0"),
					resource.TestCheckResourceAttr(dataSource2, "table_meta_infos.#", "4"),
					resource.TestCheckResourceAttrSet(dataSource2, "table_meta_infos.0.column_name"),
					resource.TestCheckResourceAttrSet(dataSource2, "table_meta_infos.0.column_type"),
					resource.TestCheckResourceAttrSet(dataSource2, "table_meta_infos.0.column_default"),
					resource.TestCheckResourceAttrSet(dataSource2, "table_meta_infos.0.is_nullable"),
				),
			},
		},
	})
}

func testAccDataSourceTaurusDBInstanceDbTableInfo_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_taurusdb_instance_db_table_info" "test" {
  instance_id = "%s"
}
`, acceptance.HW_TAURUSDB_INSTANCE_ID)
}

func testAccDataSourceTaurusDBInstanceDbTableInfo_dbTest() string {
	return fmt.Sprintf(`
data "huaweicloud_taurusdb_instance_db_table_info" "dbTest" {
  instance_id   = "%s"
  database_name = "__taurus_sys__"
}
`, acceptance.HW_TAURUSDB_INSTANCE_ID)
}

func testAccDataSourceTaurusDBInstanceDbTableInfo_tableTest() string {
	return fmt.Sprintf(`
data "huaweicloud_taurusdb_instance_db_table_info" "tableTest" {
  instance_id   = "%s"
  database_name = "__taurus_sys__"
  table_name    = "tenant"
}
`, acceptance.HW_TAURUSDB_INSTANCE_ID)
}
