package taurusdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceTaurusDBHtapStarrocksReplicationConfig_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	dataSource := "data.huaweicloud_taurusdb_htap_starrocks_replication_config.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckTaurusDBInstanceId(t)
			acceptance.TestAccPreCheckTaurusDBHtapInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTaurusDBHtapStarrocksReplicationConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSource, "task_name", rName),
					resource.TestCheckResourceAttr(dataSource, "instance_id", acceptance.HW_TAURUSDB_HTAP_INSTANCE_ID),
					resource.TestCheckResourceAttr(dataSource, "source_instance_id", acceptance.HW_TAURUSDB_INSTANCE_ID),
					resource.TestCheckResourceAttr(dataSource, "target_database_name", "ALL"),
					resource.TestCheckResourceAttr(dataSource, "is_instance_level_sync", "true"),
					resource.TestCheckResourceAttr(dataSource, "database_repl_scope", "all"),
					resource.TestCheckResourceAttr(dataSource, "is_support_reg_exp", "false"),
					resource.TestCheckResourceAttrSet(dataSource, "source_node_id"),
					resource.TestCheckResourceAttr(dataSource, "database_info.#", "1"),
					resource.TestCheckResourceAttr(dataSource, "database_info.0.database_name", "ALL"),
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

func testAccDataSourceTaurusDBHtapStarrocksReplicationConfig_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_taurusdb_htap_starrocks_replication_config" "test" {
  instance_id = huaweicloud_taurusdb_htap_starrocks_replication.test.instance_id
  task_name   = huaweicloud_taurusdb_htap_starrocks_replication.test.task_name
}
`, testAccTaurusDBHtapStarrocksReplication_base(rName))
}
