package taurusdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSqlControlHistoryRules_basic(t *testing.T) {
	dataSource := "data.huaweicloud_taurusdb_sql_control_history_rules.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckTaurusDBInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSqlControlHistoryRules_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "sql_filter_rules.#"),
					resource.TestCheckResourceAttrSet(dataSource, "sql_filter_rules.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "sql_filter_rules.0.node_id"),
					resource.TestCheckResourceAttrSet(dataSource, "sql_filter_rules.0.pattern"),
					resource.TestCheckResourceAttrSet(dataSource, "sql_filter_rules.0.sql_type"),
					resource.TestCheckResourceAttrSet(dataSource, "sql_filter_rules.0.max_concurrency"),
					resource.TestCheckResourceAttrSet(dataSource, "sql_filter_rules.0.create_at"),
					resource.TestCheckResourceAttrSet(dataSource, "sql_filter_rules.0.expire_at"),
					resource.TestCheckResourceAttrSet(dataSource, "sql_filter_rules.0.delete_at"),
					resource.TestCheckOutput("type_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceSqlControlHistoryRules_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_taurusdb_sql_control_history_rules" "test" {
  instance_id = "%[1]s"
  node_id     = "%[2]s"
}

data "huaweicloud_taurusdb_sql_control_history_rules" "type_fiter" {
  instance_id = "%[1]s"
  node_id     = "%[2]s"
  sql_type    = "UPDATE"

  depends_on = [data.huaweicloud_taurusdb_sql_control_history_rules.test]
}

output "type_filter_is_useful" {
  value = length(data.huaweicloud_taurusdb_sql_control_history_rules.type_fiter.sql_filter_rules) > 0 && alltrue(
    [for v in data.huaweicloud_taurusdb_sql_control_history_rules.type_fiter.sql_filter_rules[*].sql_type : v == "UPDATE"]
  )
}
`, acceptance.HW_TAURUSDB_INSTANCE_ID, acceptance.HW_TAURUSDB_NODE_ID)
}
