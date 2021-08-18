package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/compute/v2/servers"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/ports"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccNetworkingV2VIPAssociate_basic(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	var instance servers.Server
	var vip ports.Port
	var port ports.Port

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2VIPAssociateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2VIPAssociateConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("huaweicloud_compute_instance.test", &instance),
					testAccCheckNetworkingV2VIPExists("data.huaweicloud_networking_port.port", &port),
					testAccCheckNetworkingV2VIPExists("huaweicloud_networking_vip.vip_1", &vip),
					testAccCheckNetworkingV2VIPAssociated(&port, &vip),
				),
			},
		},
	})
}

func testAccCheckNetworkingV2VIPAssociateDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	networkingClient, err := config.NetworkingV2Client(HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_networking_vip_associate" {
			continue
		}

		vipID := rs.Primary.Attributes["vip_id"]
		_, err = ports.Get(networkingClient, vipID).Extract()
		if err != nil {
			// If the error is a 404, then the vip port does not exist,
			// and therefore the floating IP cannot be associated to it.
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return nil
			}
			return err
		}
	}

	logp.Printf("[DEBUG] Destroy NetworkingVIPAssociated success!")
	return nil
}

func testAccCheckNetworkingV2VIPAssociated(p *ports.Port, vip *ports.Port) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, ip := range p.FixedIPs {
			for _, addresspair := range vip.AllowedAddressPairs {
				if ip.IPAddress == addresspair.IPAddress {
					logp.Printf("[DEBUG] Check NetworkingVIPAssociated success!")
					return nil
				}
			}
		}

		return fmtp.Errorf("VIP %s was not attached to port %s", vip.ID, p.ID)
	}
}

func testAccNetworkingV2VIPAssociateConfig_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_networking_port" "port" {
  port_id = huaweicloud_compute_instance.test.network[0].port
}

resource "huaweicloud_networking_vip" "vip_1" {
  network_id = data.huaweicloud_vpc_subnet.test.id
}

resource "huaweicloud_networking_vip_associate" "vip_associate_1" {
  vip_id   = huaweicloud_networking_vip.vip_1.id
  port_ids = [huaweicloud_compute_instance.test.network[0].port]
}
`, testAccComputeV2Instance_basic(rName))
}
