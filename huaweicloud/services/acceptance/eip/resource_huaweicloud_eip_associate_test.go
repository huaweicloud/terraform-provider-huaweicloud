package eip

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/chnsz/golangsdk/openstack/networking/v1/eips"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccEIPAssociate_basic(t *testing.T) {
	var eip eips.PublicIp
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	associateName := "huaweicloud_vpc_eip_associate.test"
	resourceName := "huaweicloud_vpc_eip.test"

	// huaweicloud_vpc_eip_associate and huaweicloud_vpc_eip have the same ID
	// and call the same API to get resource
	rc := acceptance.InitResourceCheck(
		associateName,
		&eip,
		getEipResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccEIPAssociate_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPtr(
						associateName, "port_id", &eip.PortID),
					resource.TestCheckResourceAttrPair(
						associateName, "public_ip", resourceName, "address"),
				),
			},
			{
				ResourceName:      associateName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccEIPAssociate_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "vpc_1" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "subnet_1" {
  vpc_id     = huaweicloud_vpc.vpc_1.id
  name       = "%s"
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
}

resource "huaweicloud_networking_vip" "vip_1" {
  name       = "%s"
  network_id = huaweicloud_vpc_subnet.subnet_1.id
}

resource "huaweicloud_vpc_eip" "test" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "%s"
    size        = 5
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "huaweicloud_vpc_eip_associate" "test" {
  public_ip = huaweicloud_vpc_eip.test.address
  port_id   = huaweicloud_networking_vip.vip_1.id
}
`, rName, rName, rName, rName)
}
