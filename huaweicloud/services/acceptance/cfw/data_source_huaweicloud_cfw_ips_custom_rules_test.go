package cfw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCfwIpsCustomRules_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cfw_ips_custom_rules.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCfw(t)
			acceptance.TestAccPreCheckCfwIpsCustomRule(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCfwIpsCustomRules_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.content"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.action"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.affected_os"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.ips_id"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.ips_name"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.ips_cfw_id"),
					resource.TestCheckOutput("is_default_filter_useful", "true"),
					resource.TestCheckOutput("is_action_type_filter_useful", "true"),
					resource.TestCheckOutput("is_affected_os_filter_useful", "true"),
					resource.TestCheckOutput("is_attack_type_filter_useful", "true"),
					resource.TestCheckOutput("is_protocol_filter_useful", "true"),
					resource.TestCheckOutput("is_severity_filter_useful", "true"),
					resource.TestCheckOutput("is_software_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceCfwIpsCustomRules_basic() string {
	return fmt.Sprintf(`
%[1]s

locals {
  value = 1
}

data "huaweicloud_cfw_ips_custom_rules" "test" {
  fw_instance_id = "%[2]s" 
  object_id      = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
}

output "is_default_filter_useful" {
  value = length(data.huaweicloud_cfw_ips_custom_rules.test.records) > 0
}

data "huaweicloud_cfw_ips_custom_rules" "filter_by_action_type" {
  fw_instance_id = "%[2]s" 
  object_id      = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  action_type    = local.value
}

output "is_action_type_filter_useful" {
  value = length(data.huaweicloud_cfw_ips_custom_rules.filter_by_action_type.records) > 0 && alltrue([
    for rule in data.huaweicloud_cfw_ips_custom_rules.filter_by_action_type.records : rule.action == local.value
  ])
}

data "huaweicloud_cfw_ips_custom_rules" "filter_by_affected_os" {
  fw_instance_id = "%[2]s"
  object_id      = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  affected_os    = local.value
}

output "is_affected_os_filter_useful" {
  value = length(data.huaweicloud_cfw_ips_custom_rules.filter_by_affected_os.records) > 0 && alltrue([
    for rule in data.huaweicloud_cfw_ips_custom_rules.filter_by_affected_os.records : rule.affected_os == local.value
  ])
}

data "huaweicloud_cfw_ips_custom_rules" "filter_by_attack_type" {
  fw_instance_id = "%[2]s"
  object_id      = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  attack_type    = local.value
}

output "is_attack_type_filter_useful" {
  value = length(data.huaweicloud_cfw_ips_custom_rules.filter_by_attack_type.records) > 0 && alltrue([
    for rule in data.huaweicloud_cfw_ips_custom_rules.filter_by_attack_type.records : rule.attack_type == local.value
 ])
}

data "huaweicloud_cfw_ips_custom_rules" "filter_by_protocol" {
  fw_instance_id = "%[2]s"
  object_id      = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  protocol       = local.value
}

output "is_protocol_filter_useful" {
  value = length(data.huaweicloud_cfw_ips_custom_rules.filter_by_protocol.records) > 0 && alltrue([
    for rule in data.huaweicloud_cfw_ips_custom_rules.filter_by_protocol.records : rule.protocol == local.value
  ])
}

data "huaweicloud_cfw_ips_custom_rules" "filter_by_severity" {
  fw_instance_id = "%[2]s"
  object_id      = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  severity       = local.value
}

output "is_severity_filter_useful" {
  value = length(data.huaweicloud_cfw_ips_custom_rules.filter_by_severity.records) > 0 && alltrue([
    for rule in data.huaweicloud_cfw_ips_custom_rules.filter_by_severity.records : rule.severity == local.value
  ])
}    

data "huaweicloud_cfw_ips_custom_rules" "filter_by_software" {
  fw_instance_id = "%[2]s"
  object_id      = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  software       = local.value
}

output "is_software_filter_useful" {
  value = length(data.huaweicloud_cfw_ips_custom_rules.filter_by_software.records) > 0 && alltrue([
    for rule in data.huaweicloud_cfw_ips_custom_rules.filter_by_software.records : rule.software == local.value
  ])
}
`, testAccDatasourceFirewalls_basic(), acceptance.HW_CFW_INSTANCE_ID)
}
