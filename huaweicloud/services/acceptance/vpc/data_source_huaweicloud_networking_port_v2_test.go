package vpc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccNetworkingV2PortDataSource_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_networking_port.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      dc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2PortDataSource_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "all_fixed_ips.#", "1"),
					resource.TestCheckResourceAttrSet(dataSourceName, "mac_address"),
					resource.TestCheckResourceAttrSet(dataSourceName, "status"),
				),
			},
		},
	})
}

func testAccNetworkingV2PortDataSource_basic() string {
	rName := acceptance.RandomAccResourceName()
	cidr, gatewayIp := acceptance.RandomCidrAndGatewayIp()

	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name = "%s"
  cidr = "%s"
}

resource "huaweicloud_vpc_subnet" "test" {
  name       = "%s"
  cidr       = "%s"
  gateway_ip = "%s"
  vpc_id     = huaweicloud_vpc.test.id
}

data "huaweicloud_networking_port" "test" {
  network_id = huaweicloud_vpc_subnet.test.id
  fixed_ip   = huaweicloud_vpc_subnet.test.gateway_ip
}
`, rName, cidr, rName, cidr, gatewayIp)
}
