package vpn

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccGatewayUpgrade_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config:      testGatewayUpgrade_basic(name),
				ExpectError: regexp.MustCompile(`already new version`),
			},
		},
	})
}

func testGatewayUpgrade_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_vpn_gateway_upgrade" "test" {
  vgw_id = huaweicloud_vpn_gateway.test.id
  action = "start"
}
`, testGateway_basic(name))
}
