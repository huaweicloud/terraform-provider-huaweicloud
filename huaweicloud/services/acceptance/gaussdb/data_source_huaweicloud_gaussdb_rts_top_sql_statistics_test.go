package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGaussdbRtsTopSqlStatistics_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_rts_top_sql_statistics.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGaussDBInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGaussdbRtsTopSqlStatistics_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "top_sql_info.#"),
					resource.TestCheckResourceAttrSet(dataSource, "top_sql_info.0.node_name"),
					resource.TestCheckResourceAttrSet(dataSource, "top_sql_info.0.unique_sql_id"),
					resource.TestCheckResourceAttrSet(dataSource, "top_sql_info.0.query"),
					resource.TestCheckResourceAttrSet(dataSource, "top_sql_info.0.count"),
				),
			},
		},
	})
}

func testDataSourceGaussdbRtsTopSqlStatistics_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_gaussdb_rts_top_sql_statistics" "test" {
  instance_id = "%s"
}
`, acceptance.HW_GAUSSDB_INSTANCE_ID)
}
