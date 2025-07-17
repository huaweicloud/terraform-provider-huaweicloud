package dcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDcsInstanceNodes_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dcs_instance_nodes.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDcsInstanceNodes_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.#"),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.0.logical_node_id"),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.0.az_code"),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.0.node_role"),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.0.node_type"),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.0.node_ip"),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.0.node_port"),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.0.node_id"),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.0.priority_weight"),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.0.is_access"),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.0.group_id"),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.0.group_name"),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.0.is_remove_ip"),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.0.replication_id"),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.0.dimensions.#"),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.0.dimensions.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.0.dimensions.0.value"),
				),
			},
		},
	})
}

func testDataSourceDcsInstanceNodes_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dcs_instance_nodes" "test" {
  instance_id = huaweicloud_dcs_instance.test.id
}
`, testAccDcsV1Instance_basic(name))
}
