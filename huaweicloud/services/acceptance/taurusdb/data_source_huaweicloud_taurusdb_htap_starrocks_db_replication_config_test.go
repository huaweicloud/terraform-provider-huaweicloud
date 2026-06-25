package taurusdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceTaurusDBHtapStarrocksDbReplicationConfig_basic(t *testing.T) {
	dataSource := "data.huaweicloud_taurusdb_htap_starrocks_db_replication_config.test"
	dc := acceptance.InitDataSourceCheck(dataSource)
	rName := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckTaurusDBInstanceId(t)
			acceptance.TestAccPreCheckTaurusDBHtapInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTaurusDBHtapStarrocksDbReplicationConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSource, "instance_id", acceptance.HW_TAURUSDB_HTAP_INSTANCE_ID),
					resource.TestCheckResourceAttr(dataSource, "database", "__taurus_sys__"),
					resource.TestCheckResourceAttr(dataSource, "source_instance_id", acceptance.HW_TAURUSDB_INSTANCE_ID),
					resource.TestCheckResourceAttr(dataSource, "target_database_name", "__taurus_sys__"),
					resource.TestCheckResourceAttr(dataSource, "is_instance_level_sync", "false"),
					resource.TestCheckResourceAttr(dataSource, "database_repl_scope", "part"),
					resource.TestCheckResourceAttr(dataSource, "is_support_reg_exp", "false"),
					resource.TestCheckResourceAttrSet(dataSource, "source_node_id"),
					resource.TestCheckResourceAttr(dataSource, "database_info.#", "1"),
					resource.TestCheckResourceAttr(dataSource, "database_info.0.database_name", "__taurus_sys__"),
					resource.TestCheckResourceAttrSet(dataSource, "database_info.0.db_config_check_results.#"),
					resource.TestCheckResourceAttrSet(dataSource, "database_info.0.db_config_check_results.0.param_name"),
					resource.TestCheckResourceAttrSet(dataSource, "database_info.0.db_config_check_results.0.value"),
					resource.TestCheckResourceAttr(dataSource, "database_info.0.db_config_check_results.0.check_result", "success"),
					resource.TestCheckResourceAttrSet(dataSource, "table_infos.#"),
					resource.TestCheckResourceAttrSet(dataSource, "new_table_repl_config.#"),
				),
			},
		},
	})
}

func testAccDataSourceTaurusDBHtapStarrocksDbReplicationConfig_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_taurusdb_htap_starrocks_db_replication_config" "test" {
  instance_id = huaweicloud_taurusdb_htap_starrocks_replication.test.instance_id
  database    = huaweicloud_taurusdb_htap_starrocks_replication.test.target_database_name
}
`, testAccTaurusDBHtapStarrocksReplication_basic(rName))
}
