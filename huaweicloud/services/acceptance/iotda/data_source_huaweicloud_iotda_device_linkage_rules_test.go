package iotda

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDeviceLinkageRules_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_iotda_device_linkage_rules.test"
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
				Config: testAccDataSourceDeviceLinkageRules_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.space_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.triggers.0.type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.actions.0.type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.updated_at"),

					resource.TestCheckOutput("rule_id_filter_is_useful", "true"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("type_filter_is_useful", "true"),
					resource.TestCheckOutput("space_id_filter_is_useful", "true"),
					resource.TestCheckOutput("status_filter_is_useful", "true"),
					resource.TestCheckOutput("not_found_validation_pass", "true"),
				),
			},
		},
	})
}

func testAccDataSourceDeviceLinkageRules_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_iotda_device_linkage_rules" "test" {
  depends_on = [
    huaweicloud_iotda_device_linkage_rule.test
  ]
}

locals {
  rule_id = data.huaweicloud_iotda_device_linkage_rules.test.rules[0].id
}

data "huaweicloud_iotda_device_linkage_rules" "rule_id_filter" {
  rule_id = local.rule_id
}

output "rule_id_filter_is_useful" {
  value = length(data.huaweicloud_iotda_device_linkage_rules.rule_id_filter.rules) > 0 && alltrue(
    [for v in data.huaweicloud_iotda_device_linkage_rules.rule_id_filter.rules[*].id : v == local.rule_id]
  )
}

locals {
  name = data.huaweicloud_iotda_device_linkage_rules.test.rules[0].name
}

data "huaweicloud_iotda_device_linkage_rules" "name_filter" {
  name = local.name
}

output "name_filter_is_useful" {
  value = length(data.huaweicloud_iotda_device_linkage_rules.name_filter.rules) > 0 && alltrue(
    [for v in data.huaweicloud_iotda_device_linkage_rules.name_filter.rules[*].name : v == local.name]
  )
}

locals {
  type = data.huaweicloud_iotda_device_linkage_rules.test.rules[0].type
}

data "huaweicloud_iotda_device_linkage_rules" "type_filter" {
  type = local.type
}

output "type_filter_is_useful" {
  value = length(data.huaweicloud_iotda_device_linkage_rules.type_filter.rules) > 0 && alltrue(
    [for v in data.huaweicloud_iotda_device_linkage_rules.type_filter.rules[*].type : v == local.type]
  )
}

locals {
  space_id = data.huaweicloud_iotda_device_linkage_rules.test.rules[0].space_id
}

data "huaweicloud_iotda_device_linkage_rules" "space_id_filter" {
  space_id = local.space_id
}

output "space_id_filter_is_useful" {
  value = length(data.huaweicloud_iotda_device_linkage_rules.space_id_filter.rules) > 0 && alltrue(
    [for v in data.huaweicloud_iotda_device_linkage_rules.space_id_filter.rules[*].space_id : v == local.space_id]
  )
}

locals {
  status = data.huaweicloud_iotda_device_linkage_rules.test.rules[0].status
}

data "huaweicloud_iotda_device_linkage_rules" "status_filter" {
  status = local.status
}

output "status_filter_is_useful" {
  value = length(data.huaweicloud_iotda_device_linkage_rules.status_filter.rules) > 0 && alltrue(
    [for v in data.huaweicloud_iotda_device_linkage_rules.status_filter.rules[*].status : v == local.status]
  )
}

data "huaweicloud_iotda_device_linkage_rules" "not_found" {
  name = "not_found"
}

output "not_found_validation_pass" {
  value = length(data.huaweicloud_iotda_device_linkage_rules.not_found.rules) == 0
}
`, testDeviceLinkageRule_deviceData(name))
}
