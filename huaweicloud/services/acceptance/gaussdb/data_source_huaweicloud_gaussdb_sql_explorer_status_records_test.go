package gaussdb

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGaussDbSqlExplorerStatusRecords_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_gaussdb_sql_explorer_status_records.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGaussDBInstanceId(t)
			acceptance.TestAccPreCheckHighCostAllow(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGaussDbSqlExplorerStatusRecords_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "region"),
					resource.TestMatchResourceAttr(dataSourceName, "full_sql_switches.#", regexp.MustCompile(`^[0-9]+$`)),
					resource.TestCheckResourceAttrSet(dataSourceName, "full_sql_switches.0.is_open"),
					resource.TestCheckResourceAttrSet(dataSourceName, "full_sql_switches.0.begin_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "full_sql_switches.0.end_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "full_sql_switches.0.save_days"),
					resource.TestCheckResourceAttrSet(dataSourceName, "full_sql_switches.0.storage_mode"),
					resource.TestCheckResourceAttrSet(dataSourceName, "full_sql_switches.0.is_exclude_sys_user"),
					resource.TestCheckResourceAttrSet(dataSourceName, "full_sql_switches.0.lts_config.0.log_group_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "full_sql_switches.0.lts_config.0.log_stream_name"),
					resource.TestMatchResourceAttr(dataSourceName, "allowed_sql_types.#", regexp.MustCompile(`^[0-9]+$`)),
					resource.TestCheckResourceAttrSet(dataSourceName, "allowed_sql_types.0.category"),
					resource.TestCheckResourceAttrSet(dataSourceName, "allowed_sql_types.0.is_preset"),
				),
			},
		},
	})
}

func testDataSourceGaussDbSqlExplorerStatusRecords_basic() string {
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

data "huaweicloud_gaussdb_sql_explorer_status_records" "test" {
  instance_id = "%[1]s"
}
`, acceptance.HW_GAUSSDB_INSTANCE_ID)
}
