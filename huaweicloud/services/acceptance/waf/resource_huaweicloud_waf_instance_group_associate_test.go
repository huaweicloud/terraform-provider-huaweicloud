package waf

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/waf_hw/v1/pools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccWafInsGroupAssociate_basic(t *testing.T) {
	var group pools.Pool
	resourceName := "huaweicloud_waf_instance_group_associate.group_associate"
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
				Config: testAccWafInsGroupAssoicate_conf(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "load_balancers.#", "1"),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "load_balancers.0",
						"${huaweicloud_elb_loadbalancer.elb.id}"),
					acceptance.TestCheckResourceAttrWithVariable(resourceName,
						"group_id", "${huaweicloud_waf_instance_group.group_1.id}"),
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

func testAccWafInsGroupAssoicate_conf(name string) string {
	return fmt.Sprintf(`
%s
resource "huaweicloud_waf_instance_group" "group_1" {
  name   = "%s"
  vpc_id = huaweicloud_vpc.vpc_1.id
}

resource "huaweicloud_waf_dedicated_instance" "instance_1" {
  name               = "%s"
  available_zone     = data.huaweicloud_availability_zones.zones.names[1]
  specification_code = "waf.instance.professional"
  ecs_flavor         = data.huaweicloud_compute_flavors.flavors.ids[0]
  vpc_id             = huaweicloud_vpc.vpc_1.id
  subnet_id          = huaweicloud_vpc_subnet.vpc_subnet_1.id
  group_id           = huaweicloud_waf_instance_group.group_1.id
  
  security_group = [
    huaweicloud_networking_secgroup.secgroup.id
  ]
}


resource "huaweicloud_elb_loadbalancer" "elb" {
  name              = "%s"
  vpc_id            = huaweicloud_vpc.vpc_1.id
  cross_vpc_backend = true
  ipv4_subnet_id    = huaweicloud_vpc_subnet.vpc_subnet_1.ipv4_subnet_id

  availability_zone = [
    data.huaweicloud_availability_zones.zones.names[0],
    data.huaweicloud_availability_zones.zones.names[1]
  ]
}

resource "huaweicloud_waf_instance_group_associate" "group_associate" {
  group_id       = huaweicloud_waf_instance_group.group_1.id
  load_balancers = [huaweicloud_elb_loadbalancer.elb.id]

  depends_on = [huaweicloud_waf_dedicated_instance.instance_1]
}
`, baseDependResource(name), name, name, name)
}
