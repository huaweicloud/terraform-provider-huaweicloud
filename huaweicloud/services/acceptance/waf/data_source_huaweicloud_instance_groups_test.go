package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataSourceWafInstanceGroups_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	dataSourceName := "data.huaweicloud_waf_instance_groups.groups_1"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// WAF group is an internal feature that does not require running test cases on a daily basis.
			acceptance.TestAccPreCheckWafGroup(t)
		},
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
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_compute_flavors" "flavors" {
  availability_zone = data.huaweicloud_availability_zones.test.names[1]
  performance_type  = "normal"
  cpu_core_count    = 2
}

resource "huaweicloud_waf_dedicated_instance" "instance_1" {
  name               = "%[2]s"
  available_zone     = data.huaweicloud_availability_zones.test.names[1]
  specification_code = "waf.instance.professional"
  ecs_flavor         = data.huaweicloud_compute_flavors.flavors.ids[0]
  vpc_id             = huaweicloud_vpc.test.id
  subnet_id          = huaweicloud_vpc_subnet.test.id
  
  security_group = [
    huaweicloud_networking_secgroup.test.id
  ]
}

resource "huaweicloud_waf_instance_group" "group_1" {
  name   = "%[2]s"
  vpc_id = huaweicloud_vpc.test.id

  depends_on = [huaweicloud_waf_dedicated_instance.instance_1]
}

data "huaweicloud_waf_instance_groups" "groups_1" {
  name = "%[2]s"

  depends_on = [
    huaweicloud_waf_instance_group.group_1
  ]
}
`, common.TestBaseNetwork(name), name)
}
