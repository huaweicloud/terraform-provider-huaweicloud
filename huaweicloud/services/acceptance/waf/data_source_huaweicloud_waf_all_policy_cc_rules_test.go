package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAllPolicyCcRules_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_waf_all_policy_cc_rules.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckWafInstance(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAllPolicyCcRules_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.0.policyid"),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.0.mode"),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.0.conditions.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.0.conditions.0.category"),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.0.conditions.0.logic_operation"),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.0.conditions.0.contents.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.0.action.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.0.action.0.category"),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.0.tag_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.0.limit_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.0.limit_period"),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.0.lock_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.0.domain_aggregation"),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.0.region_aggregation"),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.0.timestamp"),

					resource.TestCheckOutput("is_policyids_filter_useful", "true"),
					resource.TestCheckOutput("is_enterprise_project_id_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceAllPolicyCcRules_base() string {
	randName := acceptance.RandomAccResourceName()

	return fmt.Sprintf(`
resource "huaweicloud_waf_policy" "test" {
  name                  = "%[1]s"
  level                 = 1
  enterprise_project_id = "0"
}

resource "huaweicloud_waf_rule_cc_protection" "test" {
  policy_id             = huaweicloud_waf_policy.test.id
  name                  = "%[1]s"
  protective_action     = "block"
  rate_limit_mode       = "cookie"
  block_page_type       = "application/json"
  page_content          = "test page content"
  user_identifier       = "test_identifier"
  limit_num             = 10
  limit_period          = 60
  lock_time             = 5
  request_aggregation   = true
  all_waf_instances     = true
  description           = "test description"
  status                = 0
  enterprise_project_id = "0"

  conditions {
    field    = "params"
    logic    = "contain"
    content  = "test content"
    subfield = "test_subfield"
  }

  conditions {
    field   = "ip"
    logic   = "equal"
    content = "192.168.0.1"
  }
}
`, randName)
}

func testAccDataSourceAllPolicyCcRules_basic() string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_waf_all_policy_cc_rules" "test" {
  depends_on = [huaweicloud_waf_rule_cc_protection.test]
}

# Filter using policyids.
locals {
  policyid = data.huaweicloud_waf_all_policy_cc_rules.test.items[0].policyid
}

data "huaweicloud_waf_all_policy_cc_rules" "policyids_filter" {
  depends_on = [huaweicloud_waf_rule_cc_protection.test]

  policyids = local.policyid
}

output "is_policyids_filter_useful" {
  value = length(data.huaweicloud_waf_all_policy_cc_rules.policyids_filter.items) > 0 && alltrue(
    [for v in data.huaweicloud_waf_all_policy_cc_rules.policyids_filter.items[*].policyid : v == local.policyid]
  )
}

# Filter using enterprise_project_id.
data "huaweicloud_waf_all_policy_cc_rules" "enterprise_project_id_filter" {
  depends_on = [huaweicloud_waf_rule_cc_protection.test]

  enterprise_project_id = "0"
}

output "is_enterprise_project_id_filter_useful" {
  value = length(data.huaweicloud_waf_all_policy_cc_rules.enterprise_project_id_filter.items) > 0 
}
`, testAccDataSourceAllPolicyCcRules_base())
}
