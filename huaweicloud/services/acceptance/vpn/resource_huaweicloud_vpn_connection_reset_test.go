package vpn

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccGatewayConnectionReset_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	ipAddress := "172.16.1.2"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testGatewayConnectionReset_basic(name, ipAddress),
			},
		},
	})
}

func testGatewayConnectionReset_basic(name, ipAddress string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_vpn_connection_reset" "test" {
  connection_id = huaweicloud_vpn_connection.test.id
}
`, testConnection_basic(name, ipAddress))
}
