package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNetworkingSecGroupV3DataSource_basic(t *testing.T) {
	var rName = fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "data.huaweicloud_networking_secgroup.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingSecGroupV3DataSource_group(rName),
			},
			{
				Config: testAccNetworkingSecGroupV3DataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSecGroupV3DataSourceID(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "rules.#"),
				),
			},
		},
	})
}

func TestAccNetworkingSecGroupV3DataSource_byID(t *testing.T) {
	var rName = fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "data.huaweicloud_networking_secgroup.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingSecGroupV3DataSource_group(rName),
			},
			{
				Config: testAccNetworkingSecGroupV3DataSource_secGroupID(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSecGroupV3DataSourceID(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "rules.#"),
				),
			},
		},
	})
}

func testAccCheckNetworkingSecGroupV3DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find security group data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Security group data source ID not set")
		}

		return nil
	}
}

func testAccNetworkingSecGroupV3DataSource_group(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_networking_secgroup" "test" {
  name                 = "%[1]s"
  description          = "My neutron security group"
  delete_default_rules = true
}

resource "huaweicloud_networking_secgroup_rule" "ports" {
  security_group_id = huaweicloud_networking_secgroup.test.id

  direction        = "ingress"
  action           = "allow"
  ethertype        = "IPv4"
  ports            = "80-100,8080"
  protocol         = "tcp"
  remote_ip_prefix = "0.0.0.0/0"
  priority         = 5
}

resource "huaweicloud_networking_secgroup_rule" "port_range" {
  security_group_id = huaweicloud_networking_secgroup.test.id

  direction        = "ingress"
  ethertype        = "IPv4"
  port_range_min   = 101
  port_range_max   = 200
  protocol         = "tcp"
  remote_ip_prefix = "0.0.0.0/0"
}

resource "huaweicloud_networking_secgroup_rule" "remote_group" {
  security_group_id = huaweicloud_networking_secgroup.test.id

  direction       = "ingress"
  action          = "allow"
  ethertype       = "IPv4"
  remote_group_id = huaweicloud_networking_secgroup.test.id
}

resource "huaweicloud_vpc_address_group" "test" {
  name = "%[1]s"
  
  addresses = [
    "192.168.10.12",
    "192.168.11.0-192.168.11.240",
  ]
}

resource "huaweicloud_networking_secgroup_rule" "remote_address_group" {
  security_group_id = huaweicloud_networking_secgroup.test.id

  direction               = "ingress"
  action                  = "allow"
  ethertype               = "IPv4"
  ports                   = "8088"
  protocol                = "tcp"
  remote_address_group_id = huaweicloud_vpc_address_group.test.id
}  
`, rName)
}

func testAccNetworkingSecGroupV3DataSource_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_networking_secgroup" "test" {
  name = huaweicloud_networking_secgroup.test.name
}
`, testAccNetworkingSecGroupV3DataSource_group(rName))
}

func testAccNetworkingSecGroupV3DataSource_secGroupID(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_networking_secgroup" "test" {
  secgroup_id = huaweicloud_networking_secgroup.test.id
}
`, testAccNetworkingSecGroupV3DataSource_group(rName))
}
