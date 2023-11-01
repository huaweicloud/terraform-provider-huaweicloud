package vpn

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccVPNCustomerGatewaysDataSource_Basic(t *testing.T) {
	resourceName := "data.huaweicloud_vpn_customer_gateways.services"
	dc := acceptance.InitDataSourceCheck(resourceName)
	rName := acceptance.RandomAccResourceName()
	ipAddress := "172.16.1.2"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVPNCustomerGatewaysDataSourceBasic(rName, ipAddress),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "customer_gateways.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "customer_gateways.0.asn", "65000"),
					resource.TestCheckResourceAttr(resourceName, "customer_gateways.0.ip", ipAddress),
					resource.TestCheckResourceAttr(resourceName, "customer_gateways.0.route_mode", "bgp"),
					resource.TestCheckResourceAttrSet(resourceName, "customer_gateways.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "customer_gateways.0.region"),
					resource.TestCheckResourceAttrSet(resourceName, "customer_gateways.0.created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "customer_gateways.0.updated_at"),
				),
			},
		},
	})
}

func testAccVPNCustomerGatewaysDataSourceBasic(rName string, ipAddress string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpn_customer_gateway" "test" {
  name = "%s"
  ip   = "%s"
}

data "huaweicloud_vpn_customer_gateways" "services" {
  customer_gateway_id  = huaweicloud_vpn_customer_gateway.test.id
  route_mode           = "bgp"
  name       		   = huaweicloud_vpn_customer_gateway.test.name
  asn                  = 65000
  ip                   = huaweicloud_vpn_customer_gateway.test.ip
}`, rName, ipAddress)
}
