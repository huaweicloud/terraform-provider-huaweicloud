package gaussdb

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGaussDbSlowSqlList_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_slow_sql_list.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGaussDBInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGaussDbSlowSqlList_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "slow_sql_infos.#"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_sql_infos.0.db_name"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_sql_infos.0.schema_name"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_sql_infos.0.sql_id"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_sql_infos.0.user_name"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_sql_infos.0.node_id"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_sql_infos.0.node_name"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_sql_infos.0.sql_text"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_sql_infos.0.calls"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_sql_infos.0.avg_exec_time"),
				),
			},
		},
	})
}

func testDataSourceGaussDbSlowSqlList_basic() string {
	now := time.Now().UTC()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	startTime := startOfDay.UnixMilli()

	endOfDay := startOfDay.Add(24 * time.Hour)
	endTime := endOfDay.UnixMilli()

	return fmt.Sprintf(`
data "huaweicloud_gaussdb_instance_nodes" "test" {
  instance_id = "%s"
}

data "huaweicloud_gaussdb_slow_sql_list" "test" {
  instance_id = "%s"
  start_time  = "%d"
  end_time    = "%d"
  threshold   = 5
  node_ids    = data.huaweicloud_gaussdb_instance_nodes.test.nodes[*].id

  multi_queries {
    name      = "query"
    condition = "AND"
    is_fuzzy  = true
    values    = ["SELECT"]
  }
}
`, acceptance.HW_GAUSSDB_INSTANCE_ID, acceptance.HW_GAUSSDB_INSTANCE_ID, startTime, endTime)
}
