package elb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceL7rules_basic(t *testing.T) {
	rName := "data.huaweicloud_elb_l7rules.test"
	dc := acceptance.InitDataSourceCheck(rName)
	name := acceptance.RandomAccResourceName()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceL7rules_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "l7rules.#"),
					resource.TestCheckResourceAttrSet(rName, "l7rules.0.id"),
					resource.TestCheckResourceAttrSet(rName, "l7rules.0.type"),
					resource.TestCheckResourceAttrSet(rName, "l7rules.0.compare_type"),
					resource.TestCheckResourceAttrSet(rName, "l7rules.0.value"),
					resource.TestCheckResourceAttrSet(rName, "l7rules.0.conditions.0.value"),
					resource.TestCheckResourceAttrSet(rName, "l7rules.0.created_at"),
					resource.TestCheckResourceAttrSet(rName, "l7rules.0.updated_at"),
					resource.TestCheckOutput("l7rule_id_filter_is_useful", "true"),
					resource.TestCheckOutput("compare_type_filter_is_useful", "true"),
					resource.TestCheckOutput("value_filter_is_useful", "true"),
					resource.TestCheckOutput("type_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceL7rules_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_elb_l7rules" "test" {
  depends_on  = [huaweicloud_elb_l7rule.test]
  l7policy_id = huaweicloud_elb_l7rule.test.l7policy_id
}

locals {
  l7rule_id = huaweicloud_elb_l7rule.test.id
}
data "huaweicloud_elb_l7rules" "l7rule_id_filter" {
  l7policy_id = huaweicloud_elb_l7rule.test.l7policy_id
  l7rule_id   = huaweicloud_elb_l7rule.test.id
}
output "l7rule_id_filter_is_useful" {
  value = length(data.huaweicloud_elb_l7rules.l7rule_id_filter.l7rules) > 0 && alltrue(
  [for v in data.huaweicloud_elb_l7rules.l7rule_id_filter.l7rules[*].id : v == local.l7rule_id]
  )  
}

locals {
 compare_type = huaweicloud_elb_l7rule.test.compare_type
}
data "huaweicloud_elb_l7rules" "compare_type_filter" {
 l7policy_id  = huaweicloud_elb_l7rule.test.l7policy_id
 compare_type = huaweicloud_elb_l7rule.test.compare_type
}
output "compare_type_filter_is_useful" {
 value = length(data.huaweicloud_elb_l7rules.compare_type_filter.l7rules) > 0 && alltrue(
 [for v in data.huaweicloud_elb_l7rules.compare_type_filter.l7rules[*].compare_type : v == local.compare_type]
 )  
}

locals {
 value = huaweicloud_elb_l7rule.test.value
}
data "huaweicloud_elb_l7rules" "value_filter" {
 l7policy_id = huaweicloud_elb_l7rule.test.l7policy_id
 value       = huaweicloud_elb_l7rule.test.value
}
output "value_filter_is_useful" {
 value = length(data.huaweicloud_elb_l7rules.value_filter.l7rules) > 0 && alltrue(
 [for v in data.huaweicloud_elb_l7rules.value_filter.l7rules[*].value : v == local.value]
 )  
}

locals {
 type = huaweicloud_elb_l7rule.test.type
}
data "huaweicloud_elb_l7rules" "type_filter" {
 l7policy_id = huaweicloud_elb_l7rule.test.l7policy_id
 type        = huaweicloud_elb_l7rule.test.type
}
output "type_filter_is_useful" {
 value = length(data.huaweicloud_elb_l7rules.type_filter.l7rules) > 0 && alltrue(
 [for v in data.huaweicloud_elb_l7rules.type_filter.l7rules[*].type : v == local.type]
 )  
}

`, testAccCheckElbV3L7RuleConfig_basic(name))
}
