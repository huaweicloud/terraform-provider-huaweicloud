package waf

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/waf_hw/v1/pools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getWafInstanceGroupFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.WafDedicatedV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating HuaweiCloud WAF dedicated client: %s", err)
	}
	return pools.Get(client, state.Primary.ID)
}

func TestAccWafInstanceGroup_basic(t *testing.T) {
	var group pools.Pool
	resourceName := "huaweicloud_waf_instance_group.group_1"
	name := acceptance.RandomAccResourceName()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&group,
		getWafInstanceGroupFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccWafInstanceGroup_conf(name, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
				),
			},
			{
				Config: testAccWafInstanceGroup_conf(name, name+"_updated"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name+"_updated"),
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

func testAccWafInstanceGroup_conf(baseName, groupName string) string {
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
`, baseName, baseName, baseName, baseName, groupName)
}
