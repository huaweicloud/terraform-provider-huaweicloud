package css

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccCssVpcepConnectionsUpdate_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testCssVpcepConnectionsUpdate_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckOutput("updated_status", "rejected"),
				),
			},
		},
	})
}

func testCssVpcepConnectionsUpdate_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_css_vpcep_connections" "test" {
  cluster_id = "%s"
}

locals {
  connection_ids = [for v in data.huaweicloud_css_vpcep_connections.test.connections : v.id]
}

resource "huaweicloud_css_vpcep_connections_update" "test" {
  cluster_id       = "%s"
  action           = "reject"
  endpoint_id_list = local.connection_ids
}

data "huaweicloud_css_vpcep_connections" "after" {
  cluster_id = "%s"

  depends_on = [huaweicloud_css_vpcep_connections_update.test]
}

output "updated_status" {
  value = data.huaweicloud_css_vpcep_connections.after.connections[0].status
}
`, acceptance.HW_CSS_CLUSTER_ID, acceptance.HW_CSS_CLUSTER_ID, acceptance.HW_CSS_CLUSTER_ID)
}
