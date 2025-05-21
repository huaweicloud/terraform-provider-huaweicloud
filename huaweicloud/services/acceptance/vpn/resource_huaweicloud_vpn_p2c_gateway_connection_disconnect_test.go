package vpn

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccP2CGatewaydisconnect_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckVPNP2cGatewayId(t)
			acceptance.TestAccPreCheckVPNP2cGatewayConnectionId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testP2CGatewaydisconnect_basic(),
			},
		},
	})
}

func testP2CGatewaydisconnect_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_vpn_p2c_gateway_connection_disconnect" "test" {
  p2c_vgw_id    = "%s"
  connection_id = "%s"
}
`, acceptance.HW_VPN_P2C_GATEWAY_ID, acceptance.HW_VPN_P2C_GATEWAY_CONNECTION_ID)
}
