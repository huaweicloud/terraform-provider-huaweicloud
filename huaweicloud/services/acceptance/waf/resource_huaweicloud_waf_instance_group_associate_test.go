package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/chnsz/golangsdk/openstack/waf_hw/v1/pools"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
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
			acceptance.TestAccPrecheckWafInstance(t)
			// WAF group is an internal feature that does not require running test cases on a daily basis.
			acceptance.TestAccPreCheckWafGroup(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccWafInsGroupAssociate_conf(name),
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

func testAccWafInsGroupAssociate_conf(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_waf_instance_group" "group_1" {
  name   = "%[2]s"
  vpc_id = huaweicloud_vpc.test.id
}

resource "huaweicloud_waf_dedicated_instance" "instance_1" {
  name               = "%[2]s"
  available_zone     = data.huaweicloud_availability_zones.test.names[1]
  specification_code = "waf.instance.professional"
  ecs_flavor         = data.huaweicloud_compute_flavors.test.ids[0]
  vpc_id             = huaweicloud_vpc.test.id
  subnet_id          = huaweicloud_vpc_subnet.test.id
  group_id           = huaweicloud_waf_instance_group.group_1.id
  
  security_group = [
    huaweicloud_networking_secgroup.test.id
  ]
}

resource "huaweicloud_elb_loadbalancer" "elb" {
  name              = "%[2]s"
  vpc_id            = huaweicloud_vpc.test.id
  cross_vpc_backend = true
  ipv4_subnet_id    = huaweicloud_vpc_subnet.test.ipv4_subnet_id

  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0],
    data.huaweicloud_availability_zones.test.names[1]
  ]
}

resource "huaweicloud_waf_instance_group_associate" "group_associate" {
  group_id       = huaweicloud_waf_instance_group.group_1.id
  load_balancers = [huaweicloud_elb_loadbalancer.elb.id]

  depends_on = [huaweicloud_waf_dedicated_instance.instance_1]
}
`, common.TestBaseComputeResources(name), name)
}
