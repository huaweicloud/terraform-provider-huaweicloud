package taurusdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/taurusdb"
)

func getHtapInstanceReplicationFunc(cfg *config.Config, r *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("gaussdb", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating TaurusDB client: %s", err)
	}
	instanceId := r.Primary.Attributes["instance_id"]
	taskName := r.Primary.Attributes["task_name"]
	details, err := taurusdb.GetHtapReplicationStatus(client, instanceId, taskName)
	if err != nil {
		return nil, golangsdk.ErrDefault404{}
	}
	return details, nil
}

func TestAccTaurusDBHtapStarrocksReplication_basic(t *testing.T) {
	var obj interface{}

	resourceName := "huaweicloud_taurusdb_htap_starrocks_replication.test"
	rName := acceptance.RandomAccResourceName()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getHtapInstanceReplicationFunc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckTaurusDBInstanceId(t)
			acceptance.TestAccPreCheckTaurusDBHtapInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccTaurusDBHtapStarrocksReplication_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "task_name", rName),
					resource.TestCheckResourceAttr(resourceName, "instance_id", acceptance.HW_TAURUSDB_HTAP_INSTANCE_ID),
					resource.TestCheckResourceAttr(resourceName, "source_instance_id", acceptance.HW_TAURUSDB_INSTANCE_ID),
					resource.TestCheckResourceAttr(resourceName, "source_database_name", "__taurus_sys__"),
					resource.TestCheckResourceAttr(resourceName, "target_database_name", "__taurus_sys__"),
					resource.TestCheckResourceAttr(resourceName, "is_instance_level_sync", "false"),
					resource.TestCheckResourceAttr(resourceName, "database_repl_scope", "part"),
					resource.TestCheckResourceAttr(resourceName, "is_support_reg_exp", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "source_node_id"),
					resource.TestCheckResourceAttr(resourceName, "database_info.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "database_info.0.database_name", "__taurus_sys__"),
					resource.TestCheckResourceAttrSet(resourceName, "database_info.0.db_config_check_results.#"),
					resource.TestCheckResourceAttrSet(resourceName, "database_info.0.db_config_check_results.0.param_name"),
					resource.TestCheckResourceAttrSet(resourceName, "database_info.0.db_config_check_results.0.value"),
					resource.TestCheckResourceAttr(resourceName, "database_info.0.db_config_check_results.0.check_result", "success"),
					resource.TestCheckResourceAttrSet(resourceName, "table_infos.#"),
					resource.TestCheckResourceAttrSet(resourceName, "new_table_repl_config.#"),
					resource.TestCheckResourceAttr(resourceName, "status", "Yes"),
					resource.TestCheckResourceAttr(resourceName, "stage", "Wait"),
					resource.TestCheckResourceAttr(resourceName, "percentage", "0"),
					resource.TestCheckResourceAttr(resourceName, "is_need_repair", "false"),
					resource.TestCheckResourceAttr(resourceName, "is_main_task", "true"),
				),
			},
			{
				Config: testAccTaurusDBHtapStarrocksReplication_basicUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "task_name", rName),
					resource.TestCheckResourceAttr(resourceName, "instance_id", acceptance.HW_TAURUSDB_HTAP_INSTANCE_ID),
					resource.TestCheckResourceAttr(resourceName, "source_instance_id", acceptance.HW_TAURUSDB_INSTANCE_ID),
					resource.TestCheckResourceAttr(resourceName, "source_database_name", "__taurus_sys__"),
					resource.TestCheckResourceAttr(resourceName, "target_database_name", "__taurus_sys__"),
					resource.TestCheckResourceAttr(resourceName, "is_instance_level_sync", "false"),
					resource.TestCheckResourceAttr(resourceName, "database_repl_scope", "part"),
					resource.TestCheckResourceAttr(resourceName, "is_support_reg_exp", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "source_node_id"),
					resource.TestCheckResourceAttr(resourceName, "database_info.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "database_info.0.database_name", "__taurus_sys__"),
					resource.TestCheckResourceAttrSet(resourceName, "database_info.0.db_config_check_results.#"),
					resource.TestCheckResourceAttrSet(resourceName, "database_info.0.db_config_check_results.0.param_name"),
					resource.TestCheckResourceAttrSet(resourceName, "database_info.0.db_config_check_results.0.value"),
					resource.TestCheckResourceAttr(resourceName, "database_info.0.db_config_check_results.0.check_result", "success"),
					resource.TestCheckResourceAttrSet(resourceName, "table_infos.#"),
					resource.TestCheckResourceAttrSet(resourceName, "new_table_repl_config.#"),
					resource.TestCheckResourceAttr(resourceName, "status", "Yes"),
					resource.TestCheckResourceAttr(resourceName, "stage", "Incremental"),
					resource.TestCheckResourceAttr(resourceName, "percentage", "100"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testAccTaurusDBHtapStarrocksReplicationImportStateIdFunc(),
				ImportStateVerifyIgnore: []string{"db_configs", "tables_configs", "enable_sync", "sync_action"},
			},
		},
	})
}

func TestAccTaurusDBHtapStarrocksReplication_instanceAllDbs(t *testing.T) {
	var obj interface{}

	resourceName := "huaweicloud_taurusdb_htap_starrocks_replication.test"
	rName := acceptance.RandomAccResourceName()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getHtapInstanceReplicationFunc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckTaurusDBInstanceId(t)
			acceptance.TestAccPreCheckTaurusDBHtapInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccTaurusDBHtapStarrocksReplication_instanceAllDbs(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "task_name", rName),
					resource.TestCheckResourceAttr(resourceName, "instance_id", acceptance.HW_TAURUSDB_HTAP_INSTANCE_ID),
					resource.TestCheckResourceAttr(resourceName, "source_instance_id", acceptance.HW_TAURUSDB_INSTANCE_ID),
					resource.TestCheckResourceAttr(resourceName, "source_database_name", "ALL"),
					resource.TestCheckResourceAttr(resourceName, "target_database_name", "ALL"),
					resource.TestCheckResourceAttr(resourceName, "is_instance_level_sync", "true"),
					resource.TestCheckResourceAttr(resourceName, "database_repl_scope", "all"),
					resource.TestCheckResourceAttr(resourceName, "is_support_reg_exp", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "source_node_id"),
					resource.TestCheckResourceAttr(resourceName, "database_info.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "database_info.0.database_name", "ALL"),
					resource.TestCheckResourceAttrSet(resourceName, "database_info.0.db_config_check_results.#"),
					resource.TestCheckResourceAttrSet(resourceName, "database_info.0.db_config_check_results.0.param_name"),
					resource.TestCheckResourceAttrSet(resourceName, "database_info.0.db_config_check_results.0.value"),
					resource.TestCheckResourceAttr(resourceName, "database_info.0.db_config_check_results.0.check_result", "success"),
					resource.TestCheckResourceAttrSet(resourceName, "table_infos.#"),
					resource.TestCheckResourceAttrSet(resourceName, "new_table_repl_config.#"),
					resource.TestCheckResourceAttr(resourceName, "status", "Yes"),
					resource.TestCheckResourceAttr(resourceName, "stage", "Incremental"),
					resource.TestCheckResourceAttr(resourceName, "percentage", "100"),
					resource.TestCheckResourceAttr(resourceName, "is_need_repair", "false"),
					resource.TestCheckResourceAttr(resourceName, "is_main_task", "true"),
				),
			},
			{
				Config: testAccTaurusDBHtapStarrocksReplication_instanceAllDbsPauseTask(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "task_name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "Yes"),
					resource.TestCheckResourceAttr(resourceName, "stage", "Paused"),
					resource.TestCheckResourceAttr(resourceName, "percentage", "100"),
				),
			},
			{
				Config: testAccTaurusDBHtapStarrocksReplication_instanceAllDbsResumeTask(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "task_name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "Yes"),
					resource.TestCheckResourceAttr(resourceName, "stage", "Incremental"),
					resource.TestCheckResourceAttr(resourceName, "percentage", "100"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testAccTaurusDBHtapStarrocksReplicationImportStateIdFunc(),
				ImportStateVerifyIgnore: []string{"db_configs", "tables_configs", "enable_sync", "sync_action"},
			},
		},
	})
}

func testAccTaurusDBHtapStarrocksReplicationImportStateIdFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var instanceId, taskName string
		for _, rs := range s.RootModule().Resources {
			if rs.Type == "huaweicloud_taurusdb_htap_starrocks_replication" {
				instanceId = rs.Primary.Attributes["instance_id"]
				taskName = rs.Primary.Attributes["task_name"]
			}
		}
		if instanceId == "" || taskName == "" {
			return "", fmt.Errorf("resource not found: %s/%s", instanceId, taskName)
		}
		return fmt.Sprintf("%s/%s", instanceId, taskName), nil
	}
}

func testAccTaurusDBHtapStarrocksReplication_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_taurusdb_htap_starrocks_replication" "test" {
  task_name              = "%[1]s"
  instance_id            = "%[2]s"
  source_instance_id     = "%[3]s"
  is_instance_level_sync = "false"
  database_repl_scope    = "part"
  source_database_name   = "__taurus_sys__"
  target_database_name   = "__taurus_sys__"

  db_configs {
    param_name = "binlog_expire_logs_seconds"
    value      = "3600"
  }

  table_repl_config {
    repl_type  = "include_tables"
    repl_scope = "part"
    tables     = ["tenant"]
  }

  tables_configs {
    table_name   = "tenant"
    table_config = "order by tenant_name"
  }
}
`, name, acceptance.HW_TAURUSDB_HTAP_INSTANCE_ID, acceptance.HW_TAURUSDB_INSTANCE_ID)
}

func testAccTaurusDBHtapStarrocksReplication_basicUpdate(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_taurusdb_htap_starrocks_replication" "test" {
  task_name              = "%[1]s"
  instance_id            = "%[2]s"
  source_instance_id     = "%[3]s"
  is_instance_level_sync = "false"
  database_repl_scope    = "part"
  source_database_name   = "__taurus_sys__"
  target_database_name   = "__taurus_sys__"

  db_configs {
    param_name = "binlog_expire_logs_seconds"
    value      = "0"
  }
  
  db_configs {
    param_name = "max_full_sync_task_threads_num"
    value      = "4"
  }

  table_repl_config {
    repl_type  = "include_tables"
    repl_scope = "part"
    tables     = ["tenant"]
  }

  tables_configs {
    table_name   = "tenant"
    table_config = "order by tenant_name;key columns tenant_name"
  }

  tables_configs {
    table_name   = "tenant_db"
    table_config = "order by db_name,tenant_name"
  }

  enable_sync            = "true"
}
`, name, acceptance.HW_TAURUSDB_HTAP_INSTANCE_ID, acceptance.HW_TAURUSDB_INSTANCE_ID)
}

func testAccTaurusDBHtapStarrocksReplication_instanceAllDbs(name string) string {
	return testAccTaurusDBHtapStarrocksReplication_base(name)
}

func testAccTaurusDBHtapStarrocksReplication_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_taurusdb_htap_starrocks_replication" "test" {
  task_name              = "%[1]s"
  instance_id            = "%[2]s"
  source_instance_id     = "%[3]s"
  is_instance_level_sync = "true"
  database_repl_scope    = "all"
  source_database_name   = "ALL"
  target_database_name   = "ALL"
  enable_sync            = "true"

  db_configs {
    param_name = "binlog_expire_logs_seconds"
    value      = "0"
  }
  
  table_repl_config {
    repl_type  = "include_tables"
    repl_scope = "all"
  }
}`, name, acceptance.HW_TAURUSDB_HTAP_INSTANCE_ID, acceptance.HW_TAURUSDB_INSTANCE_ID)
}

func testAccTaurusDBHtapStarrocksReplication_instanceAllDbsPauseTask(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_taurusdb_htap_starrocks_replication" "test" {
  task_name              = "%[1]s"
  instance_id            = "%[2]s"
  source_instance_id     = "%[3]s"
  is_instance_level_sync = "true"
  database_repl_scope    = "all"
  source_database_name   = "ALL"
  target_database_name   = "ALL"
  enable_sync            = "true"
  sync_action            = "pause"

  db_configs {
    param_name = "binlog_expire_logs_seconds"
    value      = "0"
  }
  
  table_repl_config {
    repl_type  = "include_tables"
    repl_scope = "all"
  }
}
`, name, acceptance.HW_TAURUSDB_HTAP_INSTANCE_ID, acceptance.HW_TAURUSDB_INSTANCE_ID)
}

func testAccTaurusDBHtapStarrocksReplication_instanceAllDbsResumeTask(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_taurusdb_htap_starrocks_replication" "test" {
  task_name              = "%[1]s"
  instance_id            = "%[2]s"
  source_instance_id     = "%[3]s"
  is_instance_level_sync = "true"
  database_repl_scope    = "all"
  source_database_name   = "ALL"
  target_database_name   = "ALL"
  enable_sync            = "true"
  sync_action            = "resume"

  db_configs {
    param_name = "binlog_expire_logs_seconds"
    value      = "0"
  }
  
  table_repl_config {
    repl_type  = "include_tables"
    repl_scope = "all"
  }
}
`, name, acceptance.HW_TAURUSDB_HTAP_INSTANCE_ID, acceptance.HW_TAURUSDB_INSTANCE_ID)
}
