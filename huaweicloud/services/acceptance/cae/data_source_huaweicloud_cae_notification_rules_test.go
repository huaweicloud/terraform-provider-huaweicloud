package cae

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccNotificationRules_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cae_notification_rules.test"
	rName := acceptance.RandomAccResourceNameWithDash()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCaeEnvironment(t)
			acceptance.TestAccPreCheckCaeApplication(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceNotificationRules_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "rules.#", regexp.MustCompile(`[1-9]\d*`)),
					resource.TestCheckOutput("is_with_env_name_set", "true"),
					resource.TestCheckOutput("is_with_env_event_name_set", "true"),
					resource.TestCheckOutput("is_with_env_scope_environments_set", "true"),
					resource.TestCheckOutput("is_with_env_scope_app_and_com_empty", "true"),
					resource.TestCheckOutput("is_with_env_trigger_policy_type_set", "true"),
					resource.TestCheckOutput("is_with_env_trigger_policy_period_set", "true"),
					resource.TestCheckOutput("is_with_env_trigger_policy_count_set", "true"),
					resource.TestCheckOutput("is_with_env_trigger_policy_operator_set", "true"),
					resource.TestCheckOutput("is_with_env_scope_type_set", "true"),
					resource.TestCheckOutput("is_with_env_notification_endpoint_set", "true"),
					resource.TestCheckOutput("is_with_env_notification_protocol_set", "true"),
					resource.TestCheckOutput("is_with_env_notification_template_set", "true"),
					resource.TestCheckOutput("is_with_env_enabled_set", "true"),
					resource.TestCheckOutput("is_with_app_name_set", "true"),
					resource.TestCheckOutput("is_with_app_scope_type_set", "true"),
					resource.TestCheckOutput("is_with_app_scope_applications_set", "true"),
					resource.TestCheckOutput("is_with_app_scope_env_and_com_empty", "true"),
					resource.TestCheckOutput("is_with_app_trigger_policy_type_set", "true"),
					resource.TestCheckOutput("is_with_app_enabled_set", "true"),
					resource.TestCheckOutput("is_with_com_scope_type_set", "true"),
					resource.TestCheckOutput("is_with_com_scope_components_set", "true"),
					resource.TestCheckOutput("is_with_com_scope_env_and_app_empty", "true"),
				),
			},
		},
	})
}

func testDataSourceNotificationRules_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_cae_notification_rules" "test" {
  depends_on = [
    huaweicloud_cae_notification_rule.env,
    huaweicloud_cae_notification_rule.app,
    huaweicloud_cae_notification_rule.com
  ]
}

locals {
  with_env_filter_result = try([
    for v in data.huaweicloud_cae_notification_rules.test.rules : v if v.id == huaweicloud_cae_notification_rule.env.id
  ][0], [])
}

output "is_with_env_name_set" {
  value = try(local.with_env_filter_result.name == huaweicloud_cae_notification_rule.env.name, false)
}

output "is_with_env_event_name_set" {
  value = try(local.with_env_filter_result.event_name == huaweicloud_cae_notification_rule.env.event_name, false)
}

output "is_with_env_scope_environments_set" {
  value = try(local.with_env_filter_result.scope[0].environments[0] ==
  huaweicloud_cae_notification_rule.env.scope[0].environments[0], false)
}

output "is_with_env_scope_app_and_com_empty" {
  value = try(length(local.with_env_filter_result.scope[0].applications) == 0 &&
  length(local.with_env_filter_result.scope[0].components) == 0, false)
}

output "is_with_env_trigger_policy_type_set" {
  value = try(local.with_env_filter_result.trigger_policy[0].type ==
  huaweicloud_cae_notification_rule.env.trigger_policy[0].type, false)
}

output "is_with_env_trigger_policy_period_set" {
  value = try(local.with_env_filter_result.trigger_policy[0].period ==
  huaweicloud_cae_notification_rule.env.trigger_policy[0].period, false)
}

output "is_with_env_trigger_policy_count_set" {
  value = try(local.with_env_filter_result.trigger_policy[0].count ==
  huaweicloud_cae_notification_rule.env.trigger_policy[0].count, false)
}

output "is_with_env_trigger_policy_operator_set" {
  value = try(local.with_env_filter_result.trigger_policy[0].operator ==
  huaweicloud_cae_notification_rule.env.trigger_policy[0].operator, false)
}

output "is_with_env_scope_type_set" {
  value = try(local.with_env_filter_result.scope[0].type == huaweicloud_cae_notification_rule.env.scope[0].type,
  false)
}

output "is_with_env_notification_endpoint_set" {
  value = try(local.with_env_filter_result.notification[0].endpoint == huaweicloud_cae_notification_rule.env.notification[0].endpoint,
  false)
}

output "is_with_env_notification_protocol_set" {
  value = try(local.with_env_filter_result.notification[0].protocol == huaweicloud_cae_notification_rule.env.notification[0].protocol,
  false)
}

output "is_with_env_notification_template_set" {
  value = try(local.with_env_filter_result.notification[0].template == huaweicloud_cae_notification_rule.env.notification[0].template,
  false)
}

output "is_with_env_enabled_set" {
  value = try(local.with_env_filter_result.enabled == huaweicloud_cae_notification_rule.env.enabled, false)
}

locals {
  with_app_filter_result = try([
    for v in data.huaweicloud_cae_notification_rules.test.rules : v if v.id == huaweicloud_cae_notification_rule.app.id
  ][0], [])
}

output "is_with_app_name_set" {
  value = try(local.with_app_filter_result.name == huaweicloud_cae_notification_rule.app.name, false)
}

output "is_with_app_scope_type_set" {
  value = try(local.with_app_filter_result.scope[0].type == huaweicloud_cae_notification_rule.app.scope[0].type, false)
}

output "is_with_app_scope_applications_set" {
  value = try(local.with_app_filter_result.scope[0].applications[0] ==
  huaweicloud_cae_notification_rule.app.scope[0].applications[0], false)
}

output "is_with_app_scope_env_and_com_empty" {
  value = try(length(local.with_app_filter_result.scope[0].environments) == 0 &&
  length(local.with_app_filter_result.scope[0].components) == 0, false)
}

output "is_with_app_trigger_policy_type_set" {
  value = try(local.with_app_filter_result.trigger_policy[0].type ==
  huaweicloud_cae_notification_rule.app.trigger_policy[0].type, false)
}

output "is_with_app_enabled_set" {
  value = try(local.with_app_filter_result.enabled == huaweicloud_cae_notification_rule.app.enabled, false)
}

locals {
  with_com_filter_result = try([
    for v in data.huaweicloud_cae_notification_rules.test.rules : v if v.id == huaweicloud_cae_notification_rule.com.id
  ][0], [])
}

output "is_with_com_scope_type_set" {
  value = try(local.with_com_filter_result.scope[0].type == huaweicloud_cae_notification_rule.com.scope[0].type, false)
}

output "is_with_com_scope_components_set" {
  value = try(length(local.with_com_filter_result.scope[0].components) ==
  length(huaweicloud_cae_notification_rule.com.scope[0].components), false)
}

output "is_with_com_scope_env_and_app_empty" {
  value = try(length(local.with_com_filter_result.scope[0].environments) == 0 &&
  length(local.with_com_filter_result.scope[0].applications) == 0, false)
}
`, testAccResourceNotificationRule_basic_step1(name))
}
