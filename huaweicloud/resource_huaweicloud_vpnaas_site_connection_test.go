package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/vpnaas/siteconnections"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func TestAccVpnSiteConnectionV2_basic(t *testing.T) {
	var conn siteconnections.Connection
	resourceName := "huaweicloud_vpnaas_site_connection_v2.conn_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSiteConnectionV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSiteConnectionV2_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSiteConnectionV2Exists(resourceName, &conn),
					resource.TestCheckResourceAttrPtr(resourceName, "name", &conn.Name),
					resource.TestCheckResourceAttrPtr(resourceName, "vpnservice_id", &conn.VPNServiceID),
					resource.TestCheckResourceAttrPtr(resourceName, "ikepolicy_id", &conn.IKEPolicyID),
					resource.TestCheckResourceAttrPtr(resourceName, "ipsecpolicy_id", &conn.IPSecPolicyID),
					resource.TestCheckResourceAttrPtr(resourceName, "peer_ep_group_id", &conn.PeerEPGroupID),
					resource.TestCheckResourceAttrPtr(resourceName, "local_ep_group_id", &conn.LocalEPGroupID),
					resource.TestCheckResourceAttrPtr(resourceName, "local_id", &conn.LocalID),

					resource.TestCheckResourceAttr(resourceName, "admin_state_up", "true"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
				),
			},
		},
	})
}

func testAccCheckSiteConnectionV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	networkingClient, err := config.NetworkingV2Client(HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_vpnaas_site_connection_v2" {
			continue
		}
		_, err = siteconnections.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmtp.Errorf("Site connection (%s) still exists.", rs.Primary.ID)
		}
		if _, ok := err.(golangsdk.ErrDefault404); !ok {
			return err
		}
	}
	return nil
}

func testAccCheckSiteConnectionV2Exists(n string, conn *siteconnections.Connection) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		networkingClient, err := config.NetworkingV2Client(HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating HuaweiCloud networking client: %s", err)
		}

		var found *siteconnections.Connection

		found, err = siteconnections.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}
		*conn = *found

		return nil
	}
}

var testAccSiteConnectionV2_basic = fmt.Sprintf(`
	resource "huaweicloud_networking_network_v2" "network_1" {
		name           = "tf_test_network"
  		admin_state_up = "true"
	}

	resource "huaweicloud_networking_subnet_v2" "subnet_1" {
		network_id = huaweicloud_networking_network_v2.network_1.id
  		cidr       = "192.168.199.0/24"
  		ip_version = 4
	}

	resource "huaweicloud_networking_router_v2" "router_1" {
  		name             = "my_router"
  		external_network_id = "%s"
	}

	resource "huaweicloud_networking_router_interface_v2" "router_interface_1" {
		router_id = huaweicloud_networking_router_v2.router_1.id
		subnet_id = huaweicloud_networking_subnet_v2.subnet_1.id
	}
	
	resource "huaweicloud_vpnaas_service_v2" "service_1" {
		name = "vpngw-acctest"
		router_id = huaweicloud_networking_router_v2.router_1.id
	}

	resource "huaweicloud_vpnaas_ipsec_policy_v2" "policy_1" {
	}

	resource "huaweicloud_vpnaas_ike_policy_v2" "policy_2" {
	}

	resource "huaweicloud_vpnaas_endpoint_group_v2" "group_1" {
		type = "cidr"
		endpoints = ["10.2.0.0/24", "10.3.0.0/24"]
	}
	resource "huaweicloud_vpnaas_endpoint_group_v2" "group_2" {
		type = "subnet"
		endpoints = [huaweicloud_networking_subnet_v2.subnet_1.id]
	}

	resource "huaweicloud_vpnaas_site_connection_v2" "conn_1" {
		name = "connection_1"
		ikepolicy_id = huaweicloud_vpnaas_ike_policy_v2.policy_2.id
		ipsecpolicy_id = huaweicloud_vpnaas_ipsec_policy_v2.policy_1.id
		vpnservice_id = huaweicloud_vpnaas_service_v2.service_1.id
		psk = "secret"
		peer_address = "192.168.10.1"
		peer_id = "192.168.10.1"
		local_ep_group_id = huaweicloud_vpnaas_endpoint_group_v2.group_2.id
		peer_ep_group_id = huaweicloud_vpnaas_endpoint_group_v2.group_1.id

		tags = {
			foo = "bar"
			key = "value"
		}

		depends_on = ["huaweicloud_networking_router_interface_v2.router_interface_1"]
	}
	`, HW_EXTGW_ID)
