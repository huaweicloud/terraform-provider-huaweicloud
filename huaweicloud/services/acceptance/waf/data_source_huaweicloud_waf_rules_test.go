package waf

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// If you want query the rules outside of the default enterprise project,
// please configuration the corresponding enterprise project ID.
func TestAccDataSourceRules_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_waf_rules.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRules_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSourceName, "rules.#", regexp.MustCompile(`[1-9]([0-9]*)?`)),

					resource.TestCheckOutput("policy_ids_filter_is_useful", "true"),
					resource.TestCheckOutput("eps_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceRules_basic() string {
	return fmt.Sprintf(`
variable "enterprise_project_id" {
  type    = string
  default = "%s"
}

data "huaweicloud_waf_rules" "test" {
  rule_type = "cc"
}

data "huaweicloud_waf_all_policy_cc_rules" "test" {}

locals {
  policy_id = data.huaweicloud_waf_all_policy_cc_rules.test.items[0].policyid
}

data "huaweicloud_waf_rules" "filter_by_policy_ids" {
  rule_type  = "cc"
  policy_ids = local.policy_id
}

output "policy_ids_filter_is_useful" {
  value = length(data.huaweicloud_waf_rules.filter_by_policy_ids.rules) > 0
}

data "huaweicloud_waf_rules" "filter_by_eps" {
  rule_type             = "cc"
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : "0"
}

output "eps_filter_is_useful" {
  value = length(data.huaweicloud_waf_rules.filter_by_eps.rules) > 0
}
`, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
