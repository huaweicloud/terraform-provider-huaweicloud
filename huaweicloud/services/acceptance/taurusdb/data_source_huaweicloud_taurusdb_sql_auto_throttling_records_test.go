package taurusdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSqlAutoThrottlingRecords_basic(t *testing.T) {
	dataSource := "data.huaweicloud_taurusdb_sql_auto_throttling_records.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckTaurusDBInstanceId(t)
			acceptance.TestAccPreCheckTaurusDBNodeId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSqlAutoThrottlingRecords_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "logs.#"),
					resource.TestCheckResourceAttrSet(dataSource, "logs.0.node_id"),
					resource.TestCheckResourceAttrSet(dataSource, "logs.0.pattern"),
					resource.TestCheckResourceAttrSet(dataSource, "logs.0.sql_type"),
					resource.TestCheckResourceAttrSet(dataSource, "logs.0.max_concurrency"),
					resource.TestCheckResourceAttrSet(dataSource, "logs.0.create_at"),
					resource.TestCheckResourceAttrSet(dataSource, "logs.0.expire_at"),
				),
			},
		},
	})
}

func testDataSourceSqlAutoThrottlingRecords_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_taurusdb_sql_auto_throttling_records" "test" {
  instance_id = "%[1]s"
  node_id     = "%[2]s"
}
`, acceptance.HW_TAURUSDB_INSTANCE_ID, acceptance.HW_TAURUSDB_NODE_ID)
}
