package vpn

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccP2CGatewayUpgrade_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckVPNP2cGatewayId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config:      testP2CGatewayUpgrade_basic(),
				ExpectError: regexp.MustCompile(`already new version`),
			},
		},
	})
}

func testP2CGatewayUpgrade_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_vpn_p2c_gateway_upgrade" "test" {
  p2c_vgw_id = "%[1]s"
  action     = "start"
}
`, acceptance.HW_VPN_P2C_GATEWAY_ID)
}
