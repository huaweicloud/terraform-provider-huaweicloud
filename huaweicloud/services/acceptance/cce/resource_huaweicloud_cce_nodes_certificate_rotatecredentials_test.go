package cce

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccCceRotateNodesCredentials_basic(t *testing.T) {
	resourceName := "huaweicloud_cce_nodes_certificate_rotatecredentials.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCceClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccCceRotateNodesCredentials_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", acceptance.HW_CCE_CLUSTER_ID),
				),
			},
		},
	})
}

func testAccCceRotateNodesCredentials_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_cce_nodes" "test" {
  cluster_id = "%[1]s"
}

locals {
  node_id = data.huaweicloud_cce_nodes.test.ids[0]
}

resource "huaweicloud_cce_nodes_certificate_rotatecredentials" "test" {
  cluster_id  = "%[1]s"
  api_version = "v3"
  kind        = "List"
  node_list {
    node_id = local.node_id
  }
}
`, acceptance.HW_CCE_CLUSTER_ID)
}
