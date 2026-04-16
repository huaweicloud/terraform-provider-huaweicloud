package dcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceBatchInstanceNodes_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dcs_batch_instance_nodes.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceBatchInstanceNodes_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "instances.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.node_count"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.nodes.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.nodes.0.logical_node_id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.nodes.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.nodes.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.nodes.0.az_code"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.nodes.0.node_role"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.nodes.0.node_type"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.nodes.0.node_ip"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.nodes.0.node_port"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.nodes.0.node_id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.nodes.0.priority_weight"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.nodes.0.is_access"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.nodes.0.group_id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.nodes.0.group_name"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.nodes.0.is_remove_ip"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.nodes.0.replication_id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.nodes.0.dimensions.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.nodes.0.dimensions.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.nodes.0.dimensions.0.value"),
				),
			},
		},
	})
}

func testDataSourceBatchInstanceNodes_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dcs_batch_instance_nodes" "test" {
  depends_on = [huaweicloud_dcs_instance.test]
}
`, testAccDcsV1Instance_basic(name))
}
