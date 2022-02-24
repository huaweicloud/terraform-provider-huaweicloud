package vpc

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/networking/v2/ports"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccNetworkingV2VIPAssociate_basic(t *testing.T) {
	var vip ports.Port
	var port ports.Port
	rName := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkingV2VIPAssociateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2VIPAssociateConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2VIPAssociated(&port, &vip),
				),
			},
		},
	})
}

func testAccCheckNetworkingV2VIPAssociateDestroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	networkingClient, err := config.NetworkingV2Client(acceptance.HW_REGION_NAME)
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
	cidr, gatewayIp := acceptance.RandomCidrAndGatewayIp()

	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

data "huaweicloud_images_image" "test" {
  name        = "Ubuntu 18.04 server 64bit"
  most_recent = true
}

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

resource "huaweicloud_networking_secgroup" "test" {
  name ="%s"
}

resource "huaweicloud_compute_instance" "test" {
  name               = "%s"
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [huaweicloud_networking_secgroup.test.id]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }
}

resource "huaweicloud_networking_vip" "test" {
  network_id = data.huaweicloud_vpc_subnet.test.id
}

resource "huaweicloud_networking_vip_associate" "test" {
  vip_id   = huaweicloud_networking_vip.test.id
  port_ids = [huaweicloud_compute_instance.test.network[0].port]
}
`, rName, cidr, rName, cidr, gatewayIp, rName, rName)
}
