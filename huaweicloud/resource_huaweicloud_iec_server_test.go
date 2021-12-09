package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/chnsz/golangsdk/openstack/ecs/v1/cloudservers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccIecServerResource_basic(t *testing.T) {
	var cloudserver cloudservers.CloudServer
	rName := fmt.Sprintf("iec-%s", acctest.RandString(5))
	resourceName := "huaweicloud_iec_server.server_test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIecServerDestory,
		Steps: []resource.TestStep{
			{
				Config: testAccIecServer_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIecServerExists(resourceName, &cloudserver),
					resource.TestCheckResourceAttr(resourceName, "name", "server-"+rName),
					resource.TestCheckResourceAttr(resourceName, "image_name", "Ubuntu 16.04 server 64bit"),
					resource.TestCheckResourceAttr(resourceName, "nics.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "subnet_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "security_groups.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "volume_attached.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "system_disk_type", "SAS"),
					resource.TestCheckResourceAttr(resourceName, "system_disk_size", "40"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttrSet(resourceName, "system_disk_id"),
					resource.TestCheckResourceAttrSet(resourceName, "public_ip"),
				),
			},
		},
	})
}

func testAccCheckIecServerExists(n string, cloudserver *cloudservers.CloudServer) resource.TestCheckFunc {

	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID has been seted")
		}

		config := testAccProvider.Meta().(*config.Config)
		iecClient, err := config.IECV1Client(HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating Huaweicloud IEC client: %s", err)
		}

		found, err := cloudservers.Get(iecClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmtp.Errorf("IEC Server not found")
		}
		*cloudserver = *found

		return nil
	}
}

func testAccCheckIecServerDestory(s *terraform.State) error {

	config := testAccProvider.Meta().(*config.Config)
	iecClient, err := config.IECV1Client(HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating Huaweicloud IEC client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_iec_security_group" {
			continue
		}
		_, err := cloudservers.Get(iecClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmtp.Errorf("IEC Server still exists")
		}
	}

	return nil
}

func testAccIecServer_basic(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_iec_flavors" "flavors_test" {}

data "huaweicloud_iec_images" "images_test" {
  name = "Ubuntu 16.04 server 64bit"
}

data "huaweicloud_iec_sites" "sites_test" {}

resource "huaweicloud_iec_vpc" "vpc_test" {
  name = "vpc-%s"
  cidr = "192.168.0.0/16"
  mode = "CUSTOMER"
}

resource "huaweicloud_iec_vpc_subnet" "subnet_test" {
  name = "subnet-%s"
  cidr = "192.168.0.0/16"
  gateway_ip = "192.168.0.1"
  vpc_id  = huaweicloud_iec_vpc.vpc_test.id
  site_id = data.huaweicloud_iec_sites.sites_test.sites[0].id
}

resource "huaweicloud_iec_keypair" "keypair_test" {
  name = "keypair-%s"
}

resource "huaweicloud_iec_security_group" "secgroup_test" {
  name        = "secgroup-%s"
  description = "this is a test group"
}

resource "huaweicloud_iec_security_group_rule" "rule_test" {
  direction      = "ingress"
  port_range_min = 445
  port_range_max = 445
  protocol       = "tcp"
  security_group_id = huaweicloud_iec_security_group.secgroup_test.id
  remote_ip_prefix = "0.0.0.0/0"
}

resource "huaweicloud_iec_server" "server_test" {
  name            = "server-%s"
  image_id        = data.huaweicloud_iec_images.images_test.images[0].id
  flavor_id       = data.huaweicloud_iec_flavors.flavors_test.flavors[1].id
  vpc_id          = huaweicloud_iec_vpc.vpc_test.id
  subnet_ids      = [huaweicloud_iec_vpc_subnet.subnet_test.id]
  security_groups = [huaweicloud_iec_security_group.secgroup_test.id]
  
  key_pair         = huaweicloud_iec_keypair.keypair_test.name
  bind_eip         = true
  system_disk_type = "SAS"
  system_disk_size = 40
  
  coverage_sites {
    site_id  = data.huaweicloud_iec_sites.sites_test.sites[0].id
    operator = data.huaweicloud_iec_sites.sites_test.sites[0].lines[0].operator
  }
}
`, rName, rName, rName, rName, rName)
}
