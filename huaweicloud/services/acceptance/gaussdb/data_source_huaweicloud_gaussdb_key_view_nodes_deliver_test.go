package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGaussdbKeyViewNodesDeliver_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_key_view_nodes_deliver.test"
	dc := acceptance.InitDataSourceCheck(dataSource)
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGaussdbKeyViewNodesDeliver_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.#"),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.0.node_id"),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.0.node_name"),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.0.role"),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.0.distributed_id"),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.0.component_id"),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.0.detail"),
				),
			},
		},
	})
}

func testDataSourceGaussdbKeyViewNodesDeliver_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_gaussdb_key_view_nodes_deliver" "test" {
  instance_id = huaweicloud_gaussdb_instance.test.id
}
`, testDataSourceGaussdbInstanceMetrics_base(name))
}
