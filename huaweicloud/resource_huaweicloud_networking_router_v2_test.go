package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/layer3/routers"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func TestAccNetworkingV2Router_basic(t *testing.T) {
	var router routers.Router

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDeprecated(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2RouterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Router_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2RouterExists("huaweicloud_networking_router_v2.router_1", &router),
				),
			},
			{
				ResourceName:      "huaweicloud_networking_router_v2.router_1",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccNetworkingV2Router_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"huaweicloud_networking_router_v2.router_1", "name", "router_2"),
				),
			},
		},
	})
}

func TestAccNetworkingV2Router_updateExternalGateway(t *testing.T) {
	var router routers.Router

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDeprecated(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2RouterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Router_updateExternalGateway1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2RouterExists("huaweicloud_networking_router_v2.router_1", &router),
				),
			},
			{
				Config: testAccNetworkingV2Router_updateExternalGateway2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"huaweicloud_networking_router_v2.router_1", "external_network_id", HW_EXTGW_ID),
				),
			},
		},
	})
}

func TestAccNetworkingV2Router_timeout(t *testing.T) {
	var router routers.Router

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDeprecated(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2RouterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Router_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2RouterExists("huaweicloud_networking_router_v2.router_1", &router),
				),
			},
		},
	})
}

func testAccCheckNetworkingV2RouterDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	networkingClient, err := config.NetworkingV2Client(HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_networking_router_v2" {
			continue
		}

		_, err := routers.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmtp.Errorf("Router still exists")
		}
	}

	return nil
}

func testAccCheckNetworkingV2RouterExists(n string, router *routers.Router) resource.TestCheckFunc {
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

		found, err := routers.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmtp.Errorf("Router not found")
		}

		*router = *found

		return nil
	}
}

const testAccNetworkingV2Router_basic = `
resource "huaweicloud_networking_router_v2" "router_1" {
	name = "router_1"
	admin_state_up = "true"
	distributed = "false"
}
`

const testAccNetworkingV2Router_update = `
resource "huaweicloud_networking_router_v2" "router_1" {
	name = "router_2"
	admin_state_up = "true"
	distributed = "false"
}
`

const testAccNetworkingV2Router_updateExternalGateway1 = `
resource "huaweicloud_networking_router_v2" "router_1" {
	name = "router"
	admin_state_up = "true"
	distributed = "false"
}
`

var testAccNetworkingV2Router_updateExternalGateway2 = fmt.Sprintf(`
resource "huaweicloud_networking_router_v2" "router_1" {
	name = "router"
	admin_state_up = "true"
	distributed = "false"
	external_network_id = "%s"
}
`, HW_EXTGW_ID)

const testAccNetworkingV2Router_timeout = `
resource "huaweicloud_networking_router_v2" "router_1" {
	name = "router_1"
	admin_state_up = "true"
	distributed = "false"

  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`
