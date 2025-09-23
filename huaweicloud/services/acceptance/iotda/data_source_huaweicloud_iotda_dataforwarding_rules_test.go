package iotda

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDataForwardingRules_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_iotda_dataforwarding_rules.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
		name           = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckHWIOTDAAccessAddress(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDataForwardingRules_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.resource"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.trigger"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.app_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.enabled"),

					resource.TestCheckOutput("rule_id_filter_is_useful", "true"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("resource_and_trigger_filter_is_useful", "true"),
					resource.TestCheckOutput("app_type_filter_is_useful", "true"),
					resource.TestCheckOutput("enabled_filter_is_useful", "true"),
					resource.TestCheckOutput("not_found_validation_pass", "true"),
				),
			},
		},
	})
}

func testAccDataSourceDataForwardingRules_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_iotda_dataforwarding_rules" "test" {
  depends_on = [
    huaweicloud_iotda_dataforwarding_rule.test
  ]
}

locals {
  rule_id = data.huaweicloud_iotda_dataforwarding_rules.test.rules[0].id
}

data "huaweicloud_iotda_dataforwarding_rules" "rule_id_filter" {
  rule_id = local.rule_id
}

output "rule_id_filter_is_useful" {
  value = length(data.huaweicloud_iotda_dataforwarding_rules.rule_id_filter.rules) > 0 && alltrue(
    [for v in data.huaweicloud_iotda_dataforwarding_rules.rule_id_filter.rules[*].id : v == local.rule_id]
  )
}

locals {
  name = data.huaweicloud_iotda_dataforwarding_rules.test.rules[0].name
}

data "huaweicloud_iotda_dataforwarding_rules" "name_filter" {
  name = local.name
}

output "name_filter_is_useful" {
  value = length(data.huaweicloud_iotda_dataforwarding_rules.name_filter.rules) > 0 && alltrue(
    [for v in data.huaweicloud_iotda_dataforwarding_rules.name_filter.rules[*].name : v == local.name]
  )
}

locals {
  resource = data.huaweicloud_iotda_dataforwarding_rules.test.rules[0].resource
  trigger  = data.huaweicloud_iotda_dataforwarding_rules.test.rules[0].trigger
}

data "huaweicloud_iotda_dataforwarding_rules" "resource_and_trigger_filter" {
  resource = local.resource
  trigger  = local.trigger
}

locals {
  resource_and_trigger_filter_result = [
    for v in data.huaweicloud_iotda_dataforwarding_rules.resource_and_trigger_filter.rules : 
    v if(v.trigger == local.trigger && v.resource == local.resource)
  ]
}

output "resource_and_trigger_filter_is_useful" {
  value = length(local.resource_and_trigger_filter_result) > 0 
}

locals {
  app_type = data.huaweicloud_iotda_dataforwarding_rules.test.rules[0].app_type
}

data "huaweicloud_iotda_dataforwarding_rules" "app_type_filter" {
  app_type = local.app_type
}

output "app_type_filter_is_useful" {
  value = length(data.huaweicloud_iotda_dataforwarding_rules.app_type_filter.rules) > 0 && alltrue(
    [for v in data.huaweicloud_iotda_dataforwarding_rules.app_type_filter.rules[*].app_type : v == local.app_type]
  )
}

locals {
  enabled = data.huaweicloud_iotda_dataforwarding_rules.test.rules[0].enabled
}

data "huaweicloud_iotda_dataforwarding_rules" "enabled_filter" {
  enabled = local.enabled
}

output "enabled_filter_is_useful" {
  value = length(data.huaweicloud_iotda_dataforwarding_rules.enabled_filter.rules) > 0 && alltrue(
    [for v in data.huaweicloud_iotda_dataforwarding_rules.enabled_filter.rules[*].enabled : v == local.enabled]
  )
}

data "huaweicloud_iotda_dataforwarding_rules" "not_found" {
  name = "not_found"
}

output "not_found_validation_pass" {
  value = length(data.huaweicloud_iotda_dataforwarding_rules.not_found.rules) == 0
}
`, testDataForwardingRule_basic(name))
}
