package vpc

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/networking/v2/ports"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getVipResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.NetworkingV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating HuaweiCloud Network client: %s", err)
	}
	return ports.Get(c, state.Primary.ID).Extract()
}

func TestAccNetworkingV2VIP_basic(t *testing.T) {
	var vip ports.Port
	rName := acceptance.RandomAccResourceNameWithDash()
	cidr, gatewayIp := acceptance.RandomCidrAndGatewayIp()
	resourceName := "huaweicloud_networking_vip.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&vip,
		getVipResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2VIPConfig_ipv4(rName, cidr, gatewayIp),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "ip_version", "4"),
					resource.TestCheckResourceAttrSet(resourceName, "mac_address"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccNetworkingV2VIPConfig_ipv4(rName+"-update", cidr, gatewayIp),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName+"-update"),
				),
			},
		},
	})
}

func TestAccNetworkingV2VIP_ipv6(t *testing.T) {
	var vip ports.Port
	rName := acceptance.RandomAccResourceNameWithDash()
	cidr, gatewayIp := acceptance.RandomCidrAndGatewayIp()
	resourceName := "huaweicloud_networking_vip.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&vip,
		getVipResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2VIPConfig_ipv6(rName, cidr, gatewayIp),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "ip_version", "6"),
					resource.TestCheckResourceAttrSet(resourceName, "mac_address"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNetworkingV2VIPConfig_ipv4(rName, cidr, gatewayIp string) string {
	return fmt.Sprintf(`
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

resource "huaweicloud_networking_vip" "test" {
  name       = "%s"
  network_id = huaweicloud_vpc_subnet.test.id
}
`, rName, cidr, rName, cidr, gatewayIp, rName)
}

func testAccNetworkingV2VIPConfig_ipv6(rName, cidr, gatewayIp string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name = "%s"
  cidr = "%s"
}

resource "huaweicloud_vpc_subnet" "test" {
  name        = "%s"
  cidr        = "%s"
  gateway_ip  = "%s"
  vpc_id      = huaweicloud_vpc.test.id
  ipv6_enable = true
}

resource "huaweicloud_networking_vip" "test" {
  name       = "%s"
  network_id = huaweicloud_vpc_subnet.test.id
  ip_version = 6
}
`, rName, cidr, rName, cidr, gatewayIp, rName)
}
