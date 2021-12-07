package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceWafInstanceGroups_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	dataSourceName := "data.huaweicloud_waf_instance_groups.groups_1"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWafInstanceGroups_conf(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "groups.0.name", name),
				),
			},
		},
	})
}

func testAccWafInstanceGroups_conf(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "vpc_1" {
  name = "%s_waf"
  cidr = "192.168.0.0/24"
}

resource "huaweicloud_vpc_subnet" "vpc_subnet_1" {
  name       = "%s_waf"
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
  vpc_id     = huaweicloud_vpc.vpc_1.id
}

resource "huaweicloud_networking_secgroup" "secgroup" {
  name        = "%s_waf"
  description = "terraform security group acceptance test"
}

data "huaweicloud_availability_zones" "zones" {}

data "huaweicloud_compute_flavors" "flavors" {
  availability_zone = data.huaweicloud_availability_zones.zones.names[1]
  performance_type  = "normal"
  cpu_core_count    = 2
}

resource "huaweicloud_waf_dedicated_instance" "instance_1" {
  name               = "%s"
  available_zone     = data.huaweicloud_availability_zones.zones.names[1]
  specification_code = "waf.instance.professional"
  ecs_flavor         = data.huaweicloud_compute_flavors.flavors.ids[0]
  vpc_id             = huaweicloud_vpc.vpc_1.id
  subnet_id          = huaweicloud_vpc_subnet.vpc_subnet_1.id
  
  security_group = [
    huaweicloud_networking_secgroup.secgroup.id
  ]
}

resource "huaweicloud_waf_instance_group" "group_1" {
  name   = "%s"
  vpc_id = huaweicloud_vpc.vpc_1.id

  depends_on = [huaweicloud_waf_dedicated_instance.instance_1]
}

data "huaweicloud_waf_instance_groups" "groups_1" {
  name = "%s"

  depends_on = [
    huaweicloud_waf_instance_group.group_1
  ]
}
`, name, name, name, name, name, name)
}
