package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/waf_hw/v1/pools"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
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
			// WAF group is an internal feature that does not require running test cases on a daily basis.
			acceptance.TestAccPreCheckWafGroup(t)
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
%s

resource "huaweicloud_waf_dedicated_instance" "instance_1" {
  name               = "%s"
  available_zone     = data.huaweicloud_availability_zones.test.names[0]
  specification_code = "waf.instance.professional"
  ecs_flavor         = data.huaweicloud_compute_flavors.test.ids[0]
  vpc_id             = huaweicloud_vpc.test.id
  subnet_id          = huaweicloud_vpc_subnet.test.id
  
  security_group = [
    huaweicloud_networking_secgroup.test.id
  ]
}

resource "huaweicloud_waf_instance_group" "group_1" {
  name   = "%s"
  vpc_id = huaweicloud_vpc.test.id

  depends_on = [huaweicloud_waf_dedicated_instance.instance_1]
}
`, common.TestBaseComputeResources(baseName), baseName, groupName)
}
