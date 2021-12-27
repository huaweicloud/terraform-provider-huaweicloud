package deprecated

import (
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/vpnaas/services"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccVpnServiceV2_basic(t *testing.T) {
	var service services.Service
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckDeprecated(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckVpnServiceV2Destroy,
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
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	networkingClient, err := config.NetworkingV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_vpnaas_service" {
			continue
		}
		_, err = services.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmtp.Errorf("Service (%s) still exists.", rs.Primary.ID)
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
			return fmtp.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		networkingClient, err := config.NetworkingV2Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating HuaweiCloud networking client: %s", err)
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

const testAccVpnServiceV2_basic = `
data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

resource "huaweicloud_vpnaas_service_v2" "service_1" {
  name      = "vpngw-acctest"
  router_id = data.huaweicloud_vpc.test.id
}
`
