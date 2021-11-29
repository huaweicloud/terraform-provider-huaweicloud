package cce

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCCENodesDataSource_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_cce_nodes.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)
	rName := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCCENodesDataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "nodes.0.name", rName),
				),
			},
		},
	})
}

func testAccCCENodesDataSource_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_cce_nodes" "test" {
  cluster_id = huaweicloud_cce_cluster.test.id
  name       = huaweicloud_cce_node.test.name

  depends_on = [huaweicloud_cce_node.test]
}
`, testAccCceCluster_config(rName))
}
