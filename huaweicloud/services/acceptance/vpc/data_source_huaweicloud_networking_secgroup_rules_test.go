package vpc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccNetworkingSecGroupRulesDataSource_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_networking_secgroup_rules.test"
		baseConfig     = testAccNetworkingSecGroupRulesDataSource_base()
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingSecGroupRulesDataSource_basic(baseConfig),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("is_rule_id_filter_useful", "true"),
					resource.TestCheckOutput("is_direction_filter_useful", "true"),
					resource.TestCheckOutput("is_action_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccNetworkingSecGroupRulesDataSource_base() string {
	name := acceptance.RandomAccResourceNameWithDash()

	return fmt.Sprintf(`
resource "huaweicloud_networking_secgroup" "test" {
  name        = "%s-secgroup"
  description = "terraform security group rule acceptance test"
}

resource "huaweicloud_networking_secgroup_rule" "test" {
  security_group_id = huaweicloud_networking_secgroup.test.id
  description       = "description test"
  direction         = "ingress"
  ethertype         = "IPv4"
  ports             = 80
  protocol          = "tcp"
  action            = "allow"
  remote_ip_prefix  = "0.0.0.0/0"
}
`, name)
}

func testAccNetworkingSecGroupRulesDataSource_basic(baseConfig string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_networking_secgroup_rules" "test" {
  security_group_id = huaweicloud_networking_secgroup.test.id
}

locals {
  rule_id = data.huaweicloud_networking_secgroup_rules.test.rules[0].id
}
data "huaweicloud_networking_secgroup_rules" "rule_id_filter" {
  security_group_id = huaweicloud_networking_secgroup.test.id
  rule_id           = local.rule_id
}
output "is_rule_id_filter_useful" {
  value = length(data.huaweicloud_networking_secgroup_rules.rule_id_filter.rules) > 0 && alltrue( 
    [for v in data.huaweicloud_networking_secgroup_rules.rule_id_filter.rules[*].id : v == local.rule_id]
  )  
}

locals {
  direction = data.huaweicloud_networking_secgroup_rules.test.rules[0].direction
}
data "huaweicloud_networking_secgroup_rules" "direction_filter" {
  security_group_id = huaweicloud_networking_secgroup.test.id
  direction         = local.direction
}
output "is_direction_filter_useful" {
  value = length(data.huaweicloud_networking_secgroup_rules.direction_filter.rules) > 0 && alltrue(
    [for v in data.huaweicloud_networking_secgroup_rules.direction_filter.rules[*].direction : v == local.direction]
  )  
}

locals {
  action = data.huaweicloud_networking_secgroup_rules.test.rules[0].action
}
data "huaweicloud_networking_secgroup_rules" "action_filter" {
  security_group_id = huaweicloud_networking_secgroup.test.id
  action            = local.action
}
output "is_action_filter_useful" {
  value = length(data.huaweicloud_networking_secgroup_rules.action_filter.rules) > 0 && alltrue(
    [for v in data.huaweicloud_networking_secgroup_rules.action_filter.rules[*].action : v == local.action]
  )  
}
`, baseConfig)
}
