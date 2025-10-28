package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRdsSqlStatistics_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rds_sql_statistics.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRdsSqlStatistics_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "list.0.query"),
					resource.TestCheckResourceAttrSet(dataSource, "list.0.rows"),
					resource.TestCheckResourceAttrSet(dataSource, "list.0.can_use"),
					resource.TestCheckResourceAttrSet(dataSource, "list.0.user_name"),
					resource.TestCheckResourceAttrSet(dataSource, "list.0.database"),
					resource.TestCheckResourceAttrSet(dataSource, "list.0.query_id"),
					resource.TestCheckResourceAttrSet(dataSource, "list.0.calls"),
				),
			},
		},
	})
}

func testDataSourceRdsSqlStatistics_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_rds_sql_statistics" "test" {
  instance_id = huaweicloud_rds_instance.test.id
}
`, testSqlStatisticsViewReset_base(name))
}
