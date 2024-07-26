package vpc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/networking/v1/ports"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getNetworkInterfaceResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NetworkingV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating VPC network client: %s", err)
	}

	return ports.Get(client, state.Primary.ID)
}
func TestAccVpcNetworkInterface_basic(t *testing.T) {
	var networkInterface ports.Port
	resourceName := "huaweicloud_vpc_network_interface.test"
	rName := acceptance.RandomAccResourceName()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&networkInterface,
		getNetworkInterfaceResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkInterface_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "mac_address"),
					resource.TestCheckResourceAttr(resourceName, "security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "dhcp_lease_time", "2h"),
					resource.TestCheckResourceAttr(resourceName, "allowed_addresses.0", "192.168.1.5"),
				),
			},
			{
				Config: testAccNetworkInterface_update(rName + "-update"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", ""),
					resource.TestCheckResourceAttr(resourceName, "allowed_addresses.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "security_group_ids.#", "1"),
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

func testAccNetwork_base(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}
	
resource "huaweicloud_vpc_subnet" "test" {
  vpc_id     = huaweicloud_vpc.test.id
  name       = "%[1]s"
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
}   
`, rName)
}
func testAccNetworkInterface_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_vpc_network_interface" "test" {
  name              = "%s"
  subnet_id         = huaweicloud_vpc_subnet.test.id
  fixed_ip_v4       = "192.168.0.111"
  allowed_addresses = ["192.168.1.5"]
  dhcp_lease_time   = "2h"
}
`, testAccNetwork_base(rName), rName)
}

func testAccNetworkInterface_update(rName string) string {
	return fmt.Sprintf(`
%s

%s

resource "huaweicloud_vpc_network_interface" "test" {
  subnet_id          = huaweicloud_vpc_subnet.test.id
  security_group_ids = [huaweicloud_networking_secgroup.secgroup_1.id]
  dhcp_lease_time    = "1h"
}
`, testAccNetwork_base(rName), testAccSecGroup_basic(rName))
}
