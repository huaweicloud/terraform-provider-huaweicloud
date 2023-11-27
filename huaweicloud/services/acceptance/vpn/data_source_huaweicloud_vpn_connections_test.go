package vpn

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccVPNConnectionsDataSource_Basic(t *testing.T) {
	resourceName := "data.huaweicloud_vpn_connections.services"
	dc := acceptance.InitDataSourceCheck(resourceName)
	rName := acceptance.RandomAccResourceName()
	ipAddress := "172.16.1.4"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVPNConnectionsDataSourceBasic(rName, ipAddress),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "connections.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "connections.0.status", "DOWN"),
					resource.TestCheckResourceAttr(resourceName, "connections.0.vpn_type", "static"),
					resource.TestCheckResourceAttrSet(resourceName, "connections.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "connections.0.created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "connections.0.updated_at"),
				),
			},
		},
	})
}

func testAccVPNConnectionsDataSourceBasic(rName, ipAddress string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_vpn_connections" "services" {
  status        = "DOWN"
  vpn_type      = huaweicloud_vpn_connection.test.vpn_type
  gateway_id    = huaweicloud_vpn_gateway.test.id
  gateway_ip    = huaweicloud_vpn_gateway.test.master_eip[0].id
  name          = huaweicloud_vpn_connection.test.name
  connection_id = huaweicloud_vpn_connection.test.id
}`, testConnection_basic(rName, ipAddress))
}
