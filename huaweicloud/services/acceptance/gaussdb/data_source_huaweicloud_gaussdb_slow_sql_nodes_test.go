package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGaussDbSlowSqlNodes_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_slow_sql_nodes.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGaussDBInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGaussDbSlowSqlNodes_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.#"),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.0.node_id"),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.0.node_name"),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.0.role"),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.0.component_type"),
				),
			},
		},
	})
}

func testDataSourceGaussDbSlowSqlNodes_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_gaussdb_slow_sql_nodes" "test" {
  instance_id = "%s"
  action      = "slow"
}
`, acceptance.HW_GAUSSDB_INSTANCE_ID)
}
