package waf

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceWafWebBasicProtectionRules_basic(t *testing.T) {
	dataSource := "data.huaweicloud_waf_web_basic_protection_rules.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceWafWebBasicProtectionRules_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.cve_number"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.risk_level"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.application_type"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.protection_type"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.effective_time"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.update_time"),

					resource.TestCheckOutput("is_level_filter_useful", "true"),
					resource.TestCheckOutput("is_rule_id_filter_useful", "true"),
					resource.TestCheckOutput("is_cve_number_filter_useful", "true"),
					resource.TestCheckOutput("is_risk_level_filter_useful", "true"),
					resource.TestCheckOutput("is_protection_type_names_filter_useful", "true"),
					resource.TestCheckOutput("is_application_type_names_filter_useful", "true"),
				),
			},
		},
	})
}

const testDataSourceDataSourceWafWebBasicProtectionRules_basic = `
data "huaweicloud_waf_web_basic_protection_rules" "test" {}

# Filter by level
locals {
  level = data.huaweicloud_waf_web_basic_protection_rules.test.items[0].risk_level
}

data "huaweicloud_waf_web_basic_protection_rules" "level_filter" {
  level = local.level
}

output "is_level_filter_useful" {
  value = length(data.huaweicloud_waf_web_basic_protection_rules.level_filter.items) > 0 && alltrue(
    [for v in data.huaweicloud_waf_web_basic_protection_rules.level_filter.items[*].risk_level : v == local.level]
  )
}

# Filter by rule_id
locals {
  rule_id = data.huaweicloud_waf_web_basic_protection_rules.test.items[0].id
}

data "huaweicloud_waf_web_basic_protection_rules" "rule_id_filter" {
  rule_id = local.rule_id
}

output "is_rule_id_filter_useful" {
  value = length(data.huaweicloud_waf_web_basic_protection_rules.rule_id_filter.items) > 0 && alltrue(
    [for v in data.huaweicloud_waf_web_basic_protection_rules.rule_id_filter.items[*].id : v == local.rule_id]
  )
}

# Filter by cve_number
locals {
  cve_number = data.huaweicloud_waf_web_basic_protection_rules.test.items[0].cve_number
}

data "huaweicloud_waf_web_basic_protection_rules" "cve_number_filter" {
  cve_number = local.cve_number
}

output "is_cve_number_filter_useful" {
  value = length(data.huaweicloud_waf_web_basic_protection_rules.cve_number_filter.items) > 0 && alltrue(
    [for v in data.huaweicloud_waf_web_basic_protection_rules.cve_number_filter.items[*].cve_number : v == local.cve_number]
  )
}

# Filter by risk_level
locals {
  risk_level = data.huaweicloud_waf_web_basic_protection_rules.test.items[0].risk_level
}

data "huaweicloud_waf_web_basic_protection_rules" "risk_level_filter" {
  risk_level = local.risk_level
}

output "is_risk_level_filter_useful" {
  value = length(data.huaweicloud_waf_web_basic_protection_rules.risk_level_filter.items) > 0 && alltrue(
    [for v in data.huaweicloud_waf_web_basic_protection_rules.risk_level_filter.items[*].risk_level : v == local.risk_level]
  )
}

# Filter by protection_type_names
locals {
  protection_type_names = data.huaweicloud_waf_web_basic_protection_rules.test.items[0].protection_type
}

data "huaweicloud_waf_web_basic_protection_rules" "protection_type_names_filter" {
  protection_type_names = local.protection_type_names
}

output "is_protection_type_names_filter_useful" {
  value = length(data.huaweicloud_waf_web_basic_protection_rules.protection_type_names_filter.items) > 0 && alltrue(
    [for v in data.huaweicloud_waf_web_basic_protection_rules.protection_type_names_filter.items[*].protection_type :
    v == local.protection_type_names]
  )
}

# Filter by application_type_names
locals {
  application_type_names = data.huaweicloud_waf_web_basic_protection_rules.test.items[0].application_type
}

data "huaweicloud_waf_web_basic_protection_rules" "application_type_names_filter" {
  application_type_names = local.application_type_names
}

output "is_application_type_names_filter_useful" {
  value = length(data.huaweicloud_waf_web_basic_protection_rules.application_type_names_filter.items) > 0 && alltrue(
    [for v in data.huaweicloud_waf_web_basic_protection_rules.application_type_names_filter.items[*].application_type :
    v == local.application_type_names]
  )
}
`
