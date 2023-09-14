package vpc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccNetworkingSecGroupV3DataSource_basic(t *testing.T) {
	var rName = acceptance.RandomAccResourceNameWithDash()
	dataSourceName := "data.huaweicloud_networking_secgroup.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingSecGroupV3DataSource_group(rName),
			},
			{
				Config: testAccNetworkingSecGroupV3DataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", rName),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.#"),
				),
			},
		},
	})
}

func TestAccNetworkingSecGroupV3DataSource_byID(t *testing.T) {
	var rName = acceptance.RandomAccResourceNameWithDash()
	dataSourceName := "data.huaweicloud_networking_secgroup.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingSecGroupV3DataSource_group(rName),
			},
			{
				Config: testAccNetworkingSecGroupV3DataSource_secGroupID(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", rName),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.#"),
				),
			},
		},
	})
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
