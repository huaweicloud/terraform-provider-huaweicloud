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
				Config: testAccWafInstanceGroup_conf(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
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

func testAccWafInstanceGroup_conf(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "vpc_1" {
  name = "%s_waf"
  cidr = "192.168.0.0/24"
}

resource "huaweicloud_waf_instance_group" "group_1" {
  name   = "%s"
  vpc_id = huaweicloud_vpc.vpc_1.id
}
`, name, name)
}
