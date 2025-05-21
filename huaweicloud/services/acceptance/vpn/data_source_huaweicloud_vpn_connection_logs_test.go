package vpn

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceVpnConnectionLogs_basic(t *testing.T) {
	dataSource := "data.huaweicloud_vpn_connection_logs.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceVpnConnectionLogs_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "logs.#"),
					resource.TestCheckResourceAttrSet(dataSource, "logs.0.raw_message"),
					resource.TestCheckResourceAttrSet(dataSource, "logs.0.time"),
				),
			},
		},
	})
}

func testConnectionLog_Base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test1" {
  name = "%[1]s-1"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test1" {
  name       = "%[1]s-1"
  vpc_id     = huaweicloud_vpc.test1.id
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
}

resource "huaweicloud_vpc" "test2" {
  name = "%[1]s-2"
  cidr = "172.16.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test2" {
  name       = "%[1]s-2"
  vpc_id     = huaweicloud_vpc.test2.id
  cidr       = "172.16.0.0/24"
  gateway_ip = "172.16.0.1"
}

resource "huaweicloud_vpc_eip" "test" {
  count = 4

  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "%[1]s-${count.index}"
    size        = 8
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

data "huaweicloud_vpn_gateway_availability_zones" "test" {
  flavor          = "professional1"
  attachment_type = "vpc"
}

resource "huaweicloud_vpn_gateway" "test1" {
  name               = "%[1]s-1"
  vpc_id             = huaweicloud_vpc.test1.id
  local_subnets      = [huaweicloud_vpc_subnet.test1.cidr]
  connect_subnet     = huaweicloud_vpc_subnet.test1.id
  availability_zones = [
    data.huaweicloud_vpn_gateway_availability_zones.test.names[0],
    data.huaweicloud_vpn_gateway_availability_zones.test.names[1]
  ]

  eip1 {
    id = huaweicloud_vpc_eip.test[0].id
  }

  eip2 {
    id = huaweicloud_vpc_eip.test[1].id
  }
}

resource "huaweicloud_vpn_gateway" "test2" {
  name               = "%[1]s-2"
  vpc_id             = huaweicloud_vpc.test2.id
  local_subnets      = [huaweicloud_vpc_subnet.test2.cidr]
  connect_subnet     = huaweicloud_vpc_subnet.test2.id
  availability_zones = [
    data.huaweicloud_vpn_gateway_availability_zones.test.names[0],
    data.huaweicloud_vpn_gateway_availability_zones.test.names[1]
  ]

  eip1 {
    id = huaweicloud_vpc_eip.test[2].id
  }

  eip2 {
    id = huaweicloud_vpc_eip.test[3].id
  }
}

resource "huaweicloud_vpn_customer_gateway" "test1" {
  name = "%[1]s-1"
  ip   = huaweicloud_vpc_eip.test[2].address
}

resource "huaweicloud_vpn_customer_gateway" "test2" {
  name = "%[1]s-2"
  ip   = huaweicloud_vpc_eip.test[0].address
}

resource "huaweicloud_vpn_connection" "test1" {
  name                = "%[1]s-1"
  gateway_id          = huaweicloud_vpn_gateway.test1.id
  gateway_ip          = huaweicloud_vpn_gateway.test1.master_eip[0].id
  customer_gateway_id = huaweicloud_vpn_customer_gateway.test1.id
  peer_subnets        = ["172.16.0.0/24"]
  vpn_type            = "static"
  psk                 = "Test@123"
  enable_nqa          = true
}

resource "huaweicloud_vpn_connection" "test2" {
  name                = "%[1]s-2"
  gateway_id          = huaweicloud_vpn_gateway.test2.id
  gateway_ip          = huaweicloud_vpn_gateway.test2.master_eip[0].id
  customer_gateway_id = huaweicloud_vpn_customer_gateway.test2.id
  peer_subnets        = ["192.168.0.0/24"]
  vpn_type            = "static"
  psk                 = "Test@123"
  enable_nqa          = true
}
`, name)
}

func testDataSourceVpnConnectionLogs_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_vpn_connection_logs" "test" {
  depends_on = [huaweicloud_vpn_connection.test2]

  vpn_connection_id = huaweicloud_vpn_connection.test1.id
}
`, testConnectionLog_Base(name))
}
