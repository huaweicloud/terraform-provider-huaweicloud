package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGaussdbInstanceRealTimeSessionStatistics_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_instance_real_time_session_statistics.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGaussDBInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGaussdbInstanceRealTimeSessionStatistics_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "statistics_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "statistics_list.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "statistics_list.0.active_num"),
					resource.TestCheckResourceAttrSet(dataSource, "statistics_list.0.total_num"),
					resource.TestCheckOutput("query_filter", "true"),
				),
			},
		},
	})
}

func testDataSourceGaussdbInstanceRealTimeSessionStatistics_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_gaussdb_instance_real_time_session_statistics" "test" {
  instance_id = "%[1]s"
  dimension   = "usename"
}


data "huaweicloud_gaussdb_instance_real_time_session_statistics" "query_filter" {
  instance_id = "%[1]s"
  dimension   = "usename"
  order_field = "active_num"
  order       = "DESC"
}

output "query_filter" {
  value = length(data.huaweicloud_gaussdb_instance_real_time_session_statistics.query_filter.statistics_list) > 0
}
`, acceptance.HW_GAUSSDB_INSTANCE_ID)
}
