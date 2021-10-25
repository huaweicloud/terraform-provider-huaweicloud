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

resource "huaweicloud_waf_instance_group" "group_1" {
  name   = "%s"
  vpc_id = huaweicloud_vpc.vpc_1.id
}

data "huaweicloud_waf_instance_groups" "groups_1" {
  name = "%s"

  depends_on = [
    huaweicloud_waf_instance_group.group_1
  ]
}
`, name, name, name)
}
