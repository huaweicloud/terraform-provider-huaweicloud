package huaweicloud

import (
	"fmt"
	"testing"
	"time"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/fwaas_v2/firewall_groups"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccFWFirewallGroupV2_basic(t *testing.T) {
	var epolicyID *string
	var ipolicyID *string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDeprecated(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFWFirewallGroupV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFWFirewallGroupV2_basic_1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWFirewallGroupV2("huaweicloud_fw_firewall_group_v2.fw_1", "", "", ipolicyID, epolicyID),
				),
			},
			{
				ResourceName:      "huaweicloud_fw_firewall_group_v2.fw_1",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccFWFirewallGroupV2_basic_2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWFirewallGroupV2(
						"huaweicloud_fw_firewall_group_v2.fw_1", "fw_1", "terraform acceptance test", ipolicyID, epolicyID),
				),
			},
		},
	})
}

func TestAccFWFirewallGroupV2_port0(t *testing.T) {
	var firewall_group FirewallGroup

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDeprecated(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFWFirewallGroupV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFWFirewallV2_port,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWFirewallGroupV2Exists("huaweicloud_fw_firewall_group_v2.fw_1", &firewall_group),
					testAccCheckFWFirewallPortCount(&firewall_group, 1),
				),
			},
		},
	})
}

func TestAccFWFirewallGroupV2_no_ports(t *testing.T) {
	var firewall_group FirewallGroup

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDeprecated(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFWFirewallGroupV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFWFirewallV2_no_ports,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWFirewallGroupV2Exists("huaweicloud_fw_firewall_group_v2.fw_1", &firewall_group),
					resource.TestCheckResourceAttr("huaweicloud_fw_firewall_group_v2.fw_1", "description", "firewall router test"),
					testAccCheckFWFirewallPortCount(&firewall_group, 0),
				),
			},
		},
	})
}

func TestAccFWFirewallGroupV2_port_update(t *testing.T) {
	var firewall_group FirewallGroup

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDeprecated(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFWFirewallGroupV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFWFirewallV2_port,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWFirewallGroupV2Exists("huaweicloud_fw_firewall_group_v2.fw_1", &firewall_group),
					testAccCheckFWFirewallPortCount(&firewall_group, 1),
				),
			},
			{
				Config: testAccFWFirewallV2_port_add,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWFirewallGroupV2Exists("huaweicloud_fw_firewall_group_v2.fw_1", &firewall_group),
					testAccCheckFWFirewallPortCount(&firewall_group, 2),
				),
			},
		},
	})
}

func TestAccFWFirewallGroupV2_port_remove(t *testing.T) {
	var firewall_group FirewallGroup

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDeprecated(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFWFirewallGroupV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFWFirewallV2_port,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWFirewallGroupV2Exists("huaweicloud_fw_firewall_group_v2.fw_1", &firewall_group),
					testAccCheckFWFirewallPortCount(&firewall_group, 1),
				),
			},
			{
				Config: testAccFWFirewallV2_port_remove,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWFirewallGroupV2Exists("huaweicloud_fw_firewall_group_v2.fw_1", &firewall_group),
					testAccCheckFWFirewallPortCount(&firewall_group, 0),
				),
			},
		},
	})
}

func testAccCheckFWFirewallGroupV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	fwClient, err := config.FwV2Client(HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud fw client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_firewall_group" {
			continue
		}

		_, err = firewall_groups.Get(fwClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmtp.Errorf("Firewall group (%s) still exists.", rs.Primary.ID)
		}
		if _, ok := err.(golangsdk.ErrDefault404); !ok {
			return err
		}
	}
	return nil
}

func testAccCheckFWFirewallGroupV2Exists(n string, firewall_group *FirewallGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		fwClient, err := config.FwV2Client(HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Exists) Error creating HuaweiCloud fw client: %s", err)
		}

		var found FirewallGroup
		err = firewall_groups.Get(fwClient, rs.Primary.ID).ExtractInto(&found)
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmtp.Errorf("Firewall group not found")
		}

		*firewall_group = found

		return nil
	}
}

func testAccCheckFWFirewallPortCount(firewall_group *FirewallGroup, expected int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(firewall_group.PortIDs) != expected {
			return fmtp.Errorf("Expected %d Ports, got %d", expected, len(firewall_group.PortIDs))
		}

		return nil
	}
}

func testAccCheckFWFirewallGroupV2(n, expectedName, expectedDescription string, ipolicyID *string, epolicyID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		fwClient, err := config.FwV2Client(HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Exists) Error creating HuaweiCloud fw client: %s", err)
		}

		var found *firewall_groups.FirewallGroup
		for i := 0; i < 5; i++ {
			// Firewall creation is asynchronous. Retry some times
			// if we get a 404 error. Fail on any other error.
			found, err = firewall_groups.Get(fwClient, rs.Primary.ID).Extract()
			if err != nil {
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					//lintignore:R018
					time.Sleep(time.Second)
					continue
				}
				return err
			}
			break
		}

		switch {
		case found.Name != expectedName:
			err = fmtp.Errorf("Expected Name to be <%s> but found <%s>", expectedName, found.Name)
		case found.Description != expectedDescription:
			err = fmtp.Errorf("Expected Description to be <%s> but found <%s>",
				expectedDescription, found.Description)
		case found.IngressPolicyID == "":
			err = fmtp.Errorf("Ingress Policy should not be empty")
		case found.EgressPolicyID == "":
			err = fmtp.Errorf("Egress Policy should not be empty")
		case ipolicyID != nil && found.IngressPolicyID == *ipolicyID:
			err = fmtp.Errorf("Ingress Policy had not been correctly updated. Went from <%s> to <%s>",
				expectedName, found.Name)
		case epolicyID != nil && found.EgressPolicyID == *epolicyID:
			err = fmtp.Errorf("Egress Policy had not been correctly updated. Went from <%s> to <%s>",
				expectedName, found.Name)
		}

		if err != nil {
			return err
		}

		ipolicyID = &found.IngressPolicyID
		epolicyID = &found.EgressPolicyID

		return nil
	}
}

const testAccFWFirewallGroupV2_basic_1 = `
resource "huaweicloud_fw_firewall_group_v2" "fw_1" {
  ingress_policy_id = "${huaweicloud_fw_policy_v2.policy_1.id}"
  egress_policy_id = "${huaweicloud_fw_policy_v2.policy_1.id}"

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}

resource "huaweicloud_fw_policy_v2" "policy_1" {
  name = "policy_1"
}
`

const testAccFWFirewallGroupV2_basic_2 = `
resource "huaweicloud_fw_firewall_group_v2" "fw_1" {
  name = "fw_1"
  description = "terraform acceptance test"
  ingress_policy_id = "${huaweicloud_fw_policy_v2.policy_2.id}"
  egress_policy_id = "${huaweicloud_fw_policy_v2.policy_2.id}"
  admin_state_up = true

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}

resource "huaweicloud_fw_policy_v2" "policy_2" {
  name = "policy_2"
}
`

var testAccFWFirewallV2_port = fmt.Sprintf(`
resource "huaweicloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "huaweicloud_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  enable_dhcp = true
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"
}

resource "huaweicloud_networking_router_v2" "router_1" {
  name = "router_1"
  admin_state_up = "true"
  external_network_id = "%s"
}

resource "huaweicloud_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"

  fixed_ip {
    subnet_id =  "${huaweicloud_networking_subnet_v2.subnet_1.id}"
    #ip_address = "192.168.199.23"
  }
}

resource "huaweicloud_networking_router_interface_v2" "router_interface_1" {
  router_id = "${huaweicloud_networking_router_v2.router_1.id}"
  port_id = "${huaweicloud_networking_port_v2.port_1.id}"
}

resource "huaweicloud_fw_policy_v2" "policy_1" {
  name = "policy_1"
}

resource "huaweicloud_fw_firewall_group_v2" "fw_1" {
  name = "firewall_1"
  description = "firewall router test"
  ingress_policy_id = "${huaweicloud_fw_policy_v2.policy_1.id}"
  #egress_policy_id = "${huaweicloud_fw_policy_v2.policy_1.id}"
  ports = [
	"${huaweicloud_networking_port_v2.port_1.id}"
  ]
  depends_on = ["huaweicloud_networking_router_interface_v2.router_interface_1"]
}
`, HW_EXTGW_ID)

var testAccFWFirewallV2_port_add = fmt.Sprintf(`
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

resource "huaweicloud_networking_router_v2" "router_1" {
  name = "router_1"
  admin_state_up = "true"
  external_network_id = "%s"
}

resource "huaweicloud_networking_router_v2" "router_2" {
  name = "router_2"
  admin_state_up = "true"
  external_network_id = "%s"
}

resource "huaweicloud_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"

  fixed_ip {
    subnet_id =  "${huaweicloud_networking_subnet_v2.subnet_1.id}"
    #ip_address = "192.168.199.23"
  }
}

resource "huaweicloud_networking_port_v2" "port_2" {
  name = "port_2"
  admin_state_up = "true"
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"

  fixed_ip {
    subnet_id =  "${huaweicloud_networking_subnet_v2.subnet_1.id}"
    #ip_address = "192.168.199.24"
  }
}

resource "huaweicloud_networking_router_interface_v2" "router_interface_1" {
  router_id = "${huaweicloud_networking_router_v2.router_1.id}"
  port_id = "${huaweicloud_networking_port_v2.port_1.id}"
}

resource "huaweicloud_networking_router_interface_v2" "router_interface_2" {
  router_id = "${huaweicloud_networking_router_v2.router_2.id}"
  port_id = "${huaweicloud_networking_port_v2.port_2.id}"
}

resource "huaweicloud_fw_policy_v2" "policy_1" {
  name = "policy_1"
}

resource "huaweicloud_fw_firewall_group_v2" "fw_1" {
  name = "firewall_1"
  description = "firewall router test"
  ingress_policy_id = "${huaweicloud_fw_policy_v2.policy_1.id}"
  egress_policy_id = "${huaweicloud_fw_policy_v2.policy_1.id}"
  ports = [
	"${huaweicloud_networking_port_v2.port_1.id}",
	"${huaweicloud_networking_port_v2.port_2.id}"
  ]
  depends_on = ["huaweicloud_networking_router_interface_v2.router_interface_1", "huaweicloud_networking_router_interface_v2.router_interface_2"]
}
`, HW_EXTGW_ID, HW_EXTGW_ID)

const testAccFWFirewallV2_port_remove = `
resource "huaweicloud_fw_policy_v2" "policy_1" {
  name = "policy_1"
}

resource "huaweicloud_fw_firewall_group_v2" "fw_1" {
  name = "firewall_1"
  description = "firewall router test"
  ingress_policy_id = "${huaweicloud_fw_policy_v2.policy_1.id}"
  egress_policy_id = "${huaweicloud_fw_policy_v2.policy_1.id}"
}
`

const testAccFWFirewallV2_no_ports = `
resource "huaweicloud_fw_policy_v2" "policy_1" {
  name = "policy_1"
}

resource "huaweicloud_fw_firewall_group_v2" "fw_1" {
  name = "firewall_1"
  description = "firewall router test"
  ingress_policy_id = "${huaweicloud_fw_policy_v2.policy_1.id}"
  egress_policy_id = "${huaweicloud_fw_policy_v2.policy_1.id}"
}
`
