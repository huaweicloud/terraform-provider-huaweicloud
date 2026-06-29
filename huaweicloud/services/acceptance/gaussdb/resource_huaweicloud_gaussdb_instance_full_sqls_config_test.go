package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccGaussDbFullSqlsConfig_basic(t *testing.T) {
	var (
		resourceName  = "huaweicloud_gaussdb_instance_full_sqls_config.test"
		instanceID    = acceptance.HW_GAUSSDB_INSTANCE_ID
		logGroupName  = "GROUP_GAUSSDB_APS-" + instanceID
		logStreamName = "STREAM_APS_FULL_SQL-" + instanceID
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGaussDBInstanceId(t)
			acceptance.TestAccPreCheckHighCostAllow(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccGaussDbFullSqlsConfig_basic(instanceID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "instance_id", instanceID),
					resource.TestCheckResourceAttr(resourceName, "storage_mode", "LTS"),
					resource.TestCheckResourceAttr(resourceName, "save_days", "7"),
					resource.TestCheckResourceAttr(resourceName, "lts_config.0.log_group_name", logGroupName),
					resource.TestCheckResourceAttr(resourceName, "lts_config.0.log_stream_name", logStreamName),
					resource.TestCheckResourceAttrSet(resourceName, "is_exclude_sys_user"),
					resource.TestCheckResourceAttr(resourceName, "sql_type_range.0.category", "all"),
				),
			},
			{
				Config: testAccGaussDbFullSqlsConfig_update(instanceID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "instance_id", instanceID),
					resource.TestCheckResourceAttr(resourceName, "storage_mode", "LTS"),
					resource.TestCheckResourceAttr(resourceName, "save_days", "30"),
					resource.TestCheckResourceAttr(resourceName, "lts_config.0.log_group_name", logGroupName),
					resource.TestCheckResourceAttr(resourceName, "lts_config.0.log_stream_name", logStreamName),
					resource.TestCheckResourceAttr(resourceName, "is_exclude_sys_user", "false"),
					resource.TestCheckResourceAttr(resourceName, "sql_type_range.0.category", "custom"),
					resource.TestCheckResourceAttr(resourceName, "sql_type_range.0.prefixes.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "sql_type_range.0.prefixes.0", "gaussdb"),
					resource.TestCheckResourceAttr(resourceName, "sql_type_range.0.prefixes.1", "sql"),
				),
			},
		},
	})
}

func testAccGaussDbFullSqlsConfig_basic(instanceID string) string {
	return fmt.Sprintf(`
resource "huaweicloud_gaussdb_instance_full_sqls_config" "test" {
  instance_id         = "%[1]s"
  storage_mode        = "LTS"
  save_days           = 7
  is_exclude_sys_user = true

  lts_config {
    log_group_name  = "GROUP_GAUSSDB_APS-%[1]s"
    log_stream_name = "STREAM_APS_FULL_SQL-%[1]s"
  }

  sql_type_range {
    category = "all"
  }
}
`, instanceID)
}

func testAccGaussDbFullSqlsConfig_update(instanceID string) string {
	return fmt.Sprintf(`
resource "huaweicloud_gaussdb_instance_full_sqls_config" "test" {
  instance_id         = "%[1]s"
  storage_mode        = "LTS"
  save_days           = 30
  is_exclude_sys_user = false

  lts_config {
    log_group_name  = "GROUP_GAUSSDB_APS-%[1]s"
    log_stream_name = "STREAM_APS_FULL_SQL-%[1]s"
  }

  sql_type_range {
    category = "custom"
    prefixes = ["gaussdb", "sql"]
  }
}
`, instanceID)
}
