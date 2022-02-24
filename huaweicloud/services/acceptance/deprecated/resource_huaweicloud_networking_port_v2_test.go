package deprecated

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/extradhcpopts"
	"github.com/chnsz/golangsdk/openstack/networking/v2/networks"
	"github.com/chnsz/golangsdk/openstack/networking/v2/ports"
	"github.com/chnsz/golangsdk/openstack/networking/v2/subnets"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/deprecated"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func TestAccNetworkingV2Port_basic(t *testing.T) {
	var network networks.Network
	var port ports.Port
	var subnet subnets.Subnet

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckDeprecated(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkingV2NetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Port_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("huaweicloud_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("huaweicloud_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("huaweicloud_networking_port_v2.port_1", &port),
				),
			},
			{
				ResourceName:      "huaweicloud_networking_subnet_v2.subnet_1",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"fixed_ip",
				},
			},
		},
	})
}

func TestAccNetworkingV2Port_noip(t *testing.T) {
	var network networks.Network
	var port ports.Port
	var subnet subnets.Subnet

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckDeprecated(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkingV2NetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Port_noip,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("huaweicloud_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("huaweicloud_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("huaweicloud_networking_port_v2.port_1", &port),
					testAccCheckNetworkingV2PortCountFixedIPs(&port, 1),
				),
			},
		},
	})
}

func TestAccNetworkingV2Port_timeout(t *testing.T) {
	var network networks.Network
	var port ports.Port
	var subnet subnets.Subnet

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckDeprecated(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkingV2NetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Port_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("huaweicloud_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("huaweicloud_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("huaweicloud_networking_port_v2.port_1", &port),
				),
			},
		},
	})
}

func TestAccNetworkingV2Port_createExtraDHCPOpts(t *testing.T) {
	var network networks.Network
	var subnet subnets.Subnet
	var port ports.Port

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckDeprecated(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkingV2NetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Port_createExtraDHCPOpts,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("huaweicloud_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("huaweicloud_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("huaweicloud_networking_port_v2.port_1", &port),
					resource.TestCheckResourceAttr(
						"huaweicloud_networking_port_v2.port_1", "extra_dhcp_option.#", "2"),
				),
			},
			{
				ResourceName:      "huaweicloud_networking_subnet_v2.subnet_1",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"fixed_ip",
				},
			},
		},
	})
}

func TestAccNetworkingV2Port_updateExtraDHCPOpts(t *testing.T) {
	var network networks.Network
	var subnet subnets.Subnet
	var port ports.Port

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckDeprecated(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkingV2NetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Port_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("huaweicloud_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("huaweicloud_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("huaweicloud_networking_port_v2.port_1", &port),
				),
			},
			{
				Config: testAccNetworkingV2Port_updateExtraDHCPOpts_1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("huaweicloud_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("huaweicloud_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("huaweicloud_networking_port_v2.port_1", &port),
					resource.TestCheckResourceAttr(
						"huaweicloud_networking_port_v2.port_1", "extra_dhcp_option.#", "1"),
				),
			},
			{
				Config: testAccNetworkingV2Port_updateExtraDHCPOpts_2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("huaweicloud_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("huaweicloud_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("huaweicloud_networking_port_v2.port_1", &port),
					resource.TestCheckResourceAttr(
						"huaweicloud_networking_port_v2.port_1", "extra_dhcp_option.#", "2"),
				),
			},
			{
				Config: testAccNetworkingV2Port_updateExtraDHCPOpts_3,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("huaweicloud_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("huaweicloud_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("huaweicloud_networking_port_v2.port_1", &port),
					resource.TestCheckResourceAttr(
						"huaweicloud_networking_port_v2.port_1", "extra_dhcp_option.#", "2"),
				),
			},
			{
				Config: testAccNetworkingV2Port_updateExtraDHCPOpts_4,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("huaweicloud_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("huaweicloud_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("huaweicloud_networking_port_v2.port_1", &port),
					resource.TestCheckNoResourceAttr(
						"huaweicloud_networking_port_v2.port_1", "extra_dhcp_option"),
				),
			},
		},
	})
}

func TestExpandNetworkingPortDHCPOptsV2Create(t *testing.T) {
	r := deprecated.ResourceNetworkingPortV2()
	d := r.TestResourceData()
	d.SetId("1")
	dhcpOpts1 := map[string]interface{}{
		"ip_version": 4,
		"name":       "A",
		"value":      "true",
	}
	dhcpOpts2 := map[string]interface{}{
		"ip_version": 6,
		"name":       "B",
		"value":      "false",
	}
	extraDHCPOpts := []map[string]interface{}{dhcpOpts1, dhcpOpts2}
	d.Set("extra_dhcp_option", extraDHCPOpts)

	expectedDHCPOptions := []extradhcpopts.CreateExtraDHCPOpt{
		{
			OptName:   "B",
			OptValue:  "false",
			IPVersion: golangsdk.IPVersion(6),
		},
		{
			OptName:   "A",
			OptValue:  "true",
			IPVersion: golangsdk.IPVersion(4),
		},
	}

	actualDHCPOptions := deprecated.ExpandNetworkingPortDHCPOptsV2Create(d.Get("extra_dhcp_option").(*schema.Set))

	assert.ElementsMatch(t, expectedDHCPOptions, actualDHCPOptions)
}

func TestExpandNetworkingPortDHCPOptsEmptyV2Create(t *testing.T) {
	r := deprecated.ResourceNetworkingPortV2()
	d := r.TestResourceData()
	d.SetId("1")

	expectedDHCPOptions := []extradhcpopts.CreateExtraDHCPOpt{}

	actualDHCPOptions := deprecated.ExpandNetworkingPortDHCPOptsV2Create(d.Get("extra_dhcp_option").(*schema.Set))

	assert.ElementsMatch(t, expectedDHCPOptions, actualDHCPOptions)
}

func TestExpandNetworkingPortDHCPOptsV2Update(t *testing.T) {
	r := deprecated.ResourceNetworkingPortV2()
	d := r.TestResourceData()
	d.SetId("1")
	dhcpOpts1 := map[string]interface{}{
		"ip_version": 4,
		"name":       "A",
		"value":      "true",
	}
	dhcpOpts2 := map[string]interface{}{
		"ip_version": 6,
		"name":       "B",
		"value":      "false",
	}
	extraDHCPOpts := []map[string]interface{}{dhcpOpts1, dhcpOpts2}
	d.Set("extra_dhcp_option", extraDHCPOpts)

	optsValueTrue := "true"
	optsValueFalse := "false"
	expectedDHCPOptions := []extradhcpopts.UpdateExtraDHCPOpt{
		{
			OptName:   "B",
			OptValue:  &optsValueFalse,
			IPVersion: golangsdk.IPVersion(6),
		},
		{
			OptName:   "A",
			OptValue:  &optsValueTrue,
			IPVersion: golangsdk.IPVersion(4),
		},
	}

	actualDHCPOptions := deprecated.ExpandNetworkingPortDHCPOptsV2Update(d.Get("extra_dhcp_option").(*schema.Set))

	assert.ElementsMatch(t, expectedDHCPOptions, actualDHCPOptions)
}

func TestExpandNetworkingPortDHCPOptsEmptyV2Update(t *testing.T) {
	r := deprecated.ResourceNetworkingPortV2()
	d := r.TestResourceData()
	d.SetId("1")

	expectedDHCPOptions := []extradhcpopts.UpdateExtraDHCPOpt{}

	actualDHCPOptions := deprecated.ExpandNetworkingPortDHCPOptsV2Update(d.Get("extra_dhcp_option").(*schema.Set))

	assert.ElementsMatch(t, expectedDHCPOptions, actualDHCPOptions)
}

func TestExpandNetworkingPortDHCPOptsV2Delete(t *testing.T) {
	r := deprecated.ResourceNetworkingPortV2()
	d := r.TestResourceData()
	d.SetId("1")
	dhcpOpts1 := map[string]interface{}{
		"ip_version": 4,
		"name":       "A",
		"value":      "true",
	}
	dhcpOpts2 := map[string]interface{}{
		"ip_version": 6,
		"name":       "B",
		"value":      "false",
	}
	extraDHCPOpts := []map[string]interface{}{dhcpOpts1, dhcpOpts2}
	d.Set("extra_dhcp_option", extraDHCPOpts)

	expectedDHCPOptions := []extradhcpopts.UpdateExtraDHCPOpt{
		{
			OptName: "B",
		},
		{
			OptName: "A",
		},
	}

	actualDHCPOptions := deprecated.ExpandNetworkingPortDHCPOptsV2Delete(d.Get("extra_dhcp_option").(*schema.Set))

	assert.ElementsMatch(t, expectedDHCPOptions, actualDHCPOptions)
}

func TestFlattenNetworkingPort2DHCPOptionsV2(t *testing.T) {
	dhcpOptions := extradhcpopts.ExtraDHCPOptsExt{
		ExtraDHCPOpts: []extradhcpopts.ExtraDHCPOpt{
			{
				OptName:   "A",
				OptValue:  "true",
				IPVersion: "4",
			},
			{
				OptName:   "B",
				OptValue:  "false",
				IPVersion: "6",
			},
		},
	}

	expectedDHCPOptions := []map[string]interface{}{
		{
			"ip_version": 4,
			"name":       "A",
			"value":      "true",
		},
		{
			"ip_version": 6,
			"name":       "B",
			"value":      "false",
		},
	}

	actualDHCPOptions := deprecated.FlattenNetworkingPortDHCPOptsV2(dhcpOptions)

	assert.ElementsMatch(t, expectedDHCPOptions, actualDHCPOptions)
}

func TestExpandNetworkingPortAllowedAddressPairsV2(t *testing.T) {
	r := deprecated.ResourceNetworkingPortV2()
	d := r.TestResourceData()
	d.SetId("1")
	addressPairs1 := map[string]interface{}{
		"ip_address":  "192.0.2.1",
		"mac_address": "mac1",
	}
	addressPairs2 := map[string]interface{}{
		"ip_address":  "198.51.100.1",
		"mac_address": "mac2",
	}
	allowedAddressPairs := []map[string]interface{}{addressPairs1, addressPairs2}
	d.Set("allowed_address_pairs", allowedAddressPairs)

	expectedAllowedAddressPairs := []ports.AddressPair{
		{
			IPAddress:  "192.0.2.1",
			MACAddress: "mac1",
		},
		{
			IPAddress:  "198.51.100.1",
			MACAddress: "mac2",
		},
	}

	actualAllowedAddressPairs := deprecated.ExpandNetworkingPortAllowedAddressPairsV2(d.Get("allowed_address_pairs").(*schema.Set))

	assert.ElementsMatch(t, expectedAllowedAddressPairs, actualAllowedAddressPairs)
}

func TestFlattenNetworkingPortAllowedAddressPairsV2(t *testing.T) {
	allowedAddressPairs := []ports.AddressPair{
		{
			IPAddress:  "192.0.2.1",
			MACAddress: "mac1",
		},
		{
			IPAddress:  "198.51.100.1",
			MACAddress: "mac2",
		},
	}
	mac := "mac3"

	expectedAllowedAddressPairs := []map[string]interface{}{
		{
			"ip_address":  "192.0.2.1",
			"mac_address": "mac1",
		},
		{
			"ip_address":  "198.51.100.1",
			"mac_address": "mac2",
		},
	}

	actualAllowedAddressPairs := deprecated.FlattenNetworkingPortAllowedAddressPairsV2(mac, allowedAddressPairs)

	assert.ElementsMatch(t, expectedAllowedAddressPairs, actualAllowedAddressPairs)
}

func TestExpandNetworkingPortFixedIPV2NoFixedIPs(t *testing.T) {
	r := deprecated.ResourceNetworkingPortV2()
	d := r.TestResourceData()
	d.SetId("1")
	d.Set("no_fixed_ip", true)

	actualFixedIP := deprecated.ExpandNetworkingPortFixedIPV2(d)

	assert.Empty(t, actualFixedIP)
}

func TestExpandNetworkingPortFixedIPV2SomeFixedIPs(t *testing.T) {
	r := deprecated.ResourceNetworkingPortV2()
	d := r.TestResourceData()
	d.SetId("1")
	fixedIP1 := map[string]interface{}{
		"subnet_id":  "aaa",
		"ip_address": "192.0.201.101",
	}
	fixedIP2 := map[string]interface{}{
		"subnet_id":  "bbb",
		"ip_address": "192.0.202.102",
	}
	fixedIP := []map[string]interface{}{fixedIP1, fixedIP2}
	d.Set("fixed_ip", fixedIP)

	expectedFixedIP := []ports.IP{
		{
			SubnetID:  "aaa",
			IPAddress: "192.0.201.101",
		},
		{
			SubnetID:  "bbb",
			IPAddress: "192.0.202.102",
		},
	}

	actualFixedIP := deprecated.ExpandNetworkingPortFixedIPV2(d)

	assert.ElementsMatch(t, expectedFixedIP, actualFixedIP)
}

func testAccCheckNetworkingV2PortDestroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	networkingClient, err := config.NetworkingV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_networking_port_v2" {
			continue
		}

		_, err := ports.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmtp.Errorf("Port still exists")
		}
	}

	return nil
}

func testAccCheckNetworkingV2PortCountFixedIPs(port *ports.Port, expected int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(port.FixedIPs) != expected {
			return fmtp.Errorf("Expected %d Fixed IPs, got %d", expected, len(port.FixedIPs))
		}

		return nil
	}
}

func testAccCheckNetworkingV2PortCountAllowedAddressPairs(
	port *ports.Port, expected int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(port.AllowedAddressPairs) != expected {
			return fmtp.Errorf("Expected %d Allowed Address Pairs, got %d", expected, len(port.AllowedAddressPairs))
		}

		return nil
	}
}

const testAccNetworkingV2Port_basic = `
resource "huaweicloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "huaweicloud_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"
}

resource "huaweicloud_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"

  fixed_ip {
    subnet_id =  "${huaweicloud_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }
}
`

const testAccNetworkingV2Port_noip = `
resource "huaweicloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "huaweicloud_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"
}

resource "huaweicloud_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"

  fixed_ip {
    subnet_id =  "${huaweicloud_networking_subnet_v2.subnet_1.id}"
  }
}
`
const testAccNetworkingV2Port_allowedAddressPairs = `
resource "huaweicloud_networking_network_v2" "vrrp_network" {
  name = "vrrp_network"
  admin_state_up = "true"
}

resource "huaweicloud_networking_subnet_v2" "vrrp_subnet" {
  name = "vrrp_subnet"
  cidr = "10.0.0.0/24"
  ip_version = 4
  network_id = "${huaweicloud_networking_network_v2.vrrp_network.id}"

  allocation_pools {
    start = "10.0.0.2"
    end = "10.0.0.200"
  }
}

resource "huaweicloud_networking_secgroup_v2" "secgroup_1" {
  name = "secgroup_1"
  description = "terraform security group acceptance test"
}

resource "huaweicloud_networking_router_v2" "vrrp_router" {
  name = "vrrp_router"
}

resource "huaweicloud_networking_router_interface_v2" "vrrp_interface" {
  router_id = "${huaweicloud_networking_router_v2.vrrp_router.id}"
  subnet_id = "${huaweicloud_networking_subnet_v2.vrrp_subnet.id}"
}

resource "huaweicloud_networking_port_v2" "vrrp_port_1" {
  name = "vrrp_port_1"
  admin_state_up = "true"
  network_id = "${huaweicloud_networking_network_v2.vrrp_network.id}"

  fixed_ip {
    subnet_id =  "${huaweicloud_networking_subnet_v2.vrrp_subnet.id}"
    ip_address = "10.0.0.202"
  }
}

resource "huaweicloud_networking_port_v2" "vrrp_port_2" {
  name = "vrrp_port_2"
  admin_state_up = "true"
  network_id = "${huaweicloud_networking_network_v2.vrrp_network.id}"

  fixed_ip {
    subnet_id =  "${huaweicloud_networking_subnet_v2.vrrp_subnet.id}"
    ip_address = "10.0.0.201"
  }
}

resource "huaweicloud_networking_port_v2" "instance_port" {
  name = "instance_port"
  admin_state_up = "true"
  network_id = "${huaweicloud_networking_network_v2.vrrp_network.id}"

  allowed_address_pairs {
    ip_address = "${huaweicloud_networking_port_v2.vrrp_port_1.fixed_ip.0.ip_address}"
    mac_address = "${huaweicloud_networking_port_v2.vrrp_port_1.mac_address}"
  }

  allowed_address_pairs {
    ip_address = "${huaweicloud_networking_port_v2.vrrp_port_2.fixed_ip.0.ip_address}"
    mac_address = "${huaweicloud_networking_port_v2.vrrp_port_2.mac_address}"
  }
}
`

const testAccNetworkingV2Port_timeout = `
resource "huaweicloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "huaweicloud_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"
}

resource "huaweicloud_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"

  fixed_ip {
    subnet_id =  "${huaweicloud_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }

  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`

const testAccNetworkingV2Port_createExtraDHCPOpts = `
resource "huaweicloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "huaweicloud_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"
}

resource "huaweicloud_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"
  fixed_ip {
    subnet_id =  "${huaweicloud_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }
  extra_dhcp_option {
    name = "optionA"
    value = "valueA"
  }
  extra_dhcp_option {
    name = "optionB"
    value = "valueB"
  }
}
`

const testAccNetworkingV2Port_updateExtraDHCPOpts_1 = `
resource "huaweicloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "huaweicloud_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"
}

resource "huaweicloud_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"
  fixed_ip {
    subnet_id =  "${huaweicloud_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }
  extra_dhcp_option {
    name = "optionC"
    value = "valueC"
  }
}
`

const testAccNetworkingV2Port_updateExtraDHCPOpts_2 = `
resource "huaweicloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "huaweicloud_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"
}

resource "huaweicloud_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"
  fixed_ip {
    subnet_id =  "${huaweicloud_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }
  extra_dhcp_option {
    name = "optionC"
    value = "valueC"
  }
  extra_dhcp_option {
    name = "optionD"
    value = "valueD"
  }
}
`

const testAccNetworkingV2Port_updateExtraDHCPOpts_3 = `
resource "huaweicloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "huaweicloud_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"
}

resource "huaweicloud_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"
  fixed_ip {
    subnet_id =  "${huaweicloud_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }
  extra_dhcp_option {
    name = "optionD"
    value = "valueD"
  }
  extra_dhcp_option {
    name = "optionE"
    value = "valueE"
  }
}
`

const testAccNetworkingV2Port_updateExtraDHCPOpts_4 = `
resource "huaweicloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "huaweicloud_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"
}

resource "huaweicloud_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"
  fixed_ip {
    subnet_id =  "${huaweicloud_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }
}
`
