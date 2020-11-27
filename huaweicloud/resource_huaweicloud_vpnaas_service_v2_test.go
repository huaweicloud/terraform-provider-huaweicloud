package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/vpnaas/services"
)

func TestAccVpnServiceV2_basic(t *testing.T) {
	var service services.Service
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpnServiceV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpnServiceV2_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpnServiceV2Exists(
						"huaweicloud_vpnaas_service_v2.service_1", &service),
					resource.TestCheckResourceAttrPtr("huaweicloud_vpnaas_service_v2.service_1", "router_id", &service.RouterID),
					resource.TestCheckResourceAttr("huaweicloud_vpnaas_service_v2.service_1", "admin_state_up", "true"),
				),
			},
		},
	})
}

func testAccCheckVpnServiceV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.NetworkingV2Client(HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_vpnaas_service" {
			continue
		}
		_, err = services.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Service (%s) still exists.", rs.Primary.ID)
		}
		if _, ok := err.(golangsdk.ErrDefault404); !ok {
			return err
		}
	}
	return nil
}

func testAccCheckVpnServiceV2Exists(n string, serv *services.Service) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		networkingClient, err := config.NetworkingV2Client(HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud networking client: %s", err)
		}

		var found *services.Service

		found, err = services.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}
		*serv = *found

		return nil
	}
}

var testAccVpnServiceV2_basic = fmt.Sprintf(`
	resource "huaweicloud_networking_router_v2" "router_1" {
	  name = "router_1"
	  admin_state_up = "true"
	  external_network_id = "%s"
	}
	resource "huaweicloud_vpnaas_service_v2" "service_1" {
		name = "vpngw-acctest"
		router_id = "${huaweicloud_networking_router_v2.router_1.id}"
	}
	`, HW_EXTGW_ID)
