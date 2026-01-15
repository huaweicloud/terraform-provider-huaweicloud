package dds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceConnectionStatistics_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_dds_connection_statistics.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDDSInstanceID(t)
			acceptance.TestAccPreCheckDDSNodeId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceConnectionStatistics_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "total_connections"),
					resource.TestCheckResourceAttrSet(dataSourceName, "total_inner_connections"),
					resource.TestCheckResourceAttrSet(dataSourceName, "total_outer_connections"),
					resource.TestCheckResourceAttrSet(dataSourceName, "inner_connections.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "inner_connections.0.client_ip"),
					resource.TestCheckResourceAttrSet(dataSourceName, "inner_connections.0.count"),
					resource.TestCheckResourceAttrSet(dataSourceName, "outer_connections.#"),

					resource.TestCheckOutput("is_node_id_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceConnectionStatistics_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dds_connection_statistics" "test" {
  instance_id = "%[1]s"
}

// filter by node ID
data "huaweicloud_dds_connection_statistics" "filter_by_node_id" {
  instance_id = "%[1]s"
  node_id     = "%[2]s"
}

output "is_node_id_filter_useful" {
  value = length(data.huaweicloud_dds_connection_statistics.filter_by_node_id.inner_connections) > 0
}
`, acceptance.HW_DDS_INSTANCE_ID, acceptance.HW_DDS_NODE_ID)
}
