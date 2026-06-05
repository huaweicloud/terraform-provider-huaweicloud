package gaussdb

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGaussDbSlowSqlDetail_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_slow_sql_detail.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGaussDBInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGaussDbSlowSqlDetail_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "slow_sql_details.#"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_sql_details.0.db_name"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_sql_details.0.schema_name"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_sql_details.0.sql_id"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_sql_details.0.user_name"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_sql_details.0.node_id"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_sql_details.0.node_name"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_sql_details.0.sql_text"),
				),
			},
		},
	})
}

func testDataSourceGaussDbSlowSqlDetail_basic() string {
	now := time.Now().UTC()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	startTime := startOfDay.UnixMilli()

	endOfDay := startOfDay.Add(24 * time.Hour)
	endTime := endOfDay.UnixMilli()

	return fmt.Sprintf(`
data "huaweicloud_gaussdb_instance_nodes" "test" {
  instance_id = "%s"
}

data "huaweicloud_gaussdb_slow_sql_list" "test_list" {
  instance_id = "%s"
  start_time  = "%d"
  end_time    = "%d"
  threshold   = 5
  node_ids    = data.huaweicloud_gaussdb_instance_nodes.test.nodes[*].id
}

data "huaweicloud_gaussdb_slow_sql_detail" "test" {
  instance_id = "%s"
  start_time  = "%d"
  end_time    = "%d"
  sql_id      = data.huaweicloud_gaussdb_slow_sql_list.test_list.slow_sql_infos.0.sql_id
  node_ids    = data.huaweicloud_gaussdb_instance_nodes.test.nodes[*].id

  multi_queries {
    name      = "query"
    condition = "AND"
    is_fuzzy  = true
    values    = ["SELECT"]
  }
}
`, acceptance.HW_GAUSSDB_INSTANCE_ID, acceptance.HW_GAUSSDB_INSTANCE_ID, startTime, endTime,
		acceptance.HW_GAUSSDB_INSTANCE_ID, startTime, endTime)
}
