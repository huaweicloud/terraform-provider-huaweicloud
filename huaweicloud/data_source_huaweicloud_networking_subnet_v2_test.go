package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNetworkingV2SubnetDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckDeprecated(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccHuaweiCloudNetworkingSubnetV2DataSource_subnet,
			},
			{
				Config: testAccHuaweiCloudNetworkingSubnetV2DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetV2DataSourceID("data.huaweicloud_networking_subnet_v2.subnet_1"),
					testAccCheckNetworkingSubnetV2DataSourceGoodNetwork("data.huaweicloud_networking_subnet_v2.subnet_1", "huaweicloud_networking_network_v2.network_1"),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_networking_subnet_v2.subnet_1", "name", "subnet_1"),
				),
			},
		},
	})
}

func TestAccNetworkingV2SubnetDataSource_testQueries(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckDeprecated(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccHuaweiCloudNetworkingSubnetV2DataSource_subnet,
			},
			{
				Config: testAccHuaweiCloudNetworkingSubnetV2DataSource_cidr,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetV2DataSourceID("data.huaweicloud_networking_subnet_v2.subnet_1"),
				),
			},
			{
				Config: testAccHuaweiCloudNetworkingSubnetV2DataSource_dhcpEnabled,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetV2DataSourceID("data.huaweicloud_networking_subnet_v2.subnet_1"),
				),
			},
			{
				Config: testAccHuaweiCloudNetworkingSubnetV2DataSource_ipVersion,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetV2DataSourceID("data.huaweicloud_networking_subnet_v2.subnet_1"),
				),
			},
			{
				Config: testAccHuaweiCloudNetworkingSubnetV2DataSource_gatewayIP,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetV2DataSourceID("data.huaweicloud_networking_subnet_v2.subnet_1"),
				),
			},
		},
	})
}

func TestAccNetworkingV2SubnetDataSource_networkIdAttribute(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckDeprecated(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccHuaweiCloudNetworkingSubnetV2DataSource_networkIdAttribute,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetV2DataSourceID("data.huaweicloud_networking_subnet_v2.subnet_1"),
					testAccCheckNetworkingSubnetV2DataSourceGoodNetwork("data.huaweicloud_networking_subnet_v2.subnet_1", "huaweicloud_networking_network_v2.network_1"),
					testAccCheckNetworkingPortV2ID("huaweicloud_networking_port_v2.port_1"),
				),
			},
		},
	})
}

func testAccCheckNetworkingSubnetV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Can't find subnet data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("Subnet data source ID not set")
		}

		return nil
	}
}

func testAccCheckNetworkingPortV2ID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Can't find port resource: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("Port resource ID not set")
		}

		return nil
	}
}

func testAccCheckNetworkingSubnetV2DataSourceGoodNetwork(n1, n2 string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ds1, ok := s.RootModule().Resources[n1]
		if !ok {
			return fmtp.Errorf("Can't find subnet data source: %s", n1)
		}

		if ds1.Primary.ID == "" {
			return fmtp.Errorf("Subnet data source ID not set")
		}

		rs2, ok := s.RootModule().Resources[n2]
		if !ok {
			return fmtp.Errorf("Can't find network resource: %s", n2)
		}

		if rs2.Primary.ID == "" {
			return fmtp.Errorf("Network resource ID not set")
		}

		if rs2.Primary.ID != ds1.Primary.Attributes["network_id"] {
			return fmtp.Errorf("Network id and subnet network_id don't match")
		}

		return nil
	}
}

const testAccHuaweiCloudNetworkingSubnetV2DataSource_subnet = `
resource "huaweicloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "huaweicloud_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"
}
`

var testAccHuaweiCloudNetworkingSubnetV2DataSource_basic = fmt.Sprintf(`
%s

data "huaweicloud_networking_subnet_v2" "subnet_1" {
	name = "${huaweicloud_networking_subnet_v2.subnet_1.name}"
}
`, testAccHuaweiCloudNetworkingSubnetV2DataSource_subnet)

var testAccHuaweiCloudNetworkingSubnetV2DataSource_cidr = fmt.Sprintf(`
%s

data "huaweicloud_networking_subnet_v2" "subnet_1" {
	cidr = "192.168.199.0/24"
}
`, testAccHuaweiCloudNetworkingSubnetV2DataSource_subnet)

var testAccHuaweiCloudNetworkingSubnetV2DataSource_dhcpEnabled = fmt.Sprintf(`
%s

data "huaweicloud_networking_subnet_v2" "subnet_1" {
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"
	dhcp_enabled = true
}
`, testAccHuaweiCloudNetworkingSubnetV2DataSource_subnet)

var testAccHuaweiCloudNetworkingSubnetV2DataSource_ipVersion = fmt.Sprintf(`
%s

data "huaweicloud_networking_subnet_v2" "subnet_1" {
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"
  ip_version = 4
}
`, testAccHuaweiCloudNetworkingSubnetV2DataSource_subnet)

var testAccHuaweiCloudNetworkingSubnetV2DataSource_gatewayIP = fmt.Sprintf(`
%s

data "huaweicloud_networking_subnet_v2" "subnet_1" {
  gateway_ip = "${huaweicloud_networking_subnet_v2.subnet_1.gateway_ip}"
}
`, testAccHuaweiCloudNetworkingSubnetV2DataSource_subnet)

var testAccHuaweiCloudNetworkingSubnetV2DataSource_networkIdAttribute = fmt.Sprintf(`
%s

data "huaweicloud_networking_subnet_v2" "subnet_1" {
  subnet_id = "${huaweicloud_networking_subnet_v2.subnet_1.id}"
}

resource "huaweicloud_networking_port_v2" "port_1" {
  name               = "test_port"
  network_id         = "${data.huaweicloud_networking_subnet_v2.subnet_1.network_id}"
  admin_state_up  = "true"
}

`, testAccHuaweiCloudNetworkingSubnetV2DataSource_subnet)
