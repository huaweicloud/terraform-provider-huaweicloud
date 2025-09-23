package cae

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataNotificationRules_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_cae_notification_rules.test"
		rName      = acceptance.RandomAccResourceNameWithDash()
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCaeEnvironmentIds(t, 2)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataNotificationRules_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "rules.#", regexp.MustCompile(`[1-9]\d*`)),
					resource.TestCheckOutput("is_with_env_notification_rules_returned", "true"),
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
					resource.TestCheckOutput("is_with_app_notification_rules_returned", "true"),
					resource.TestCheckOutput("is_with_app_name_set", "true"),
					resource.TestCheckOutput("is_with_app_scope_type_set", "true"),
					resource.TestCheckOutput("is_with_app_scope_applications_set", "true"),
					resource.TestCheckOutput("is_with_app_scope_env_and_com_empty", "true"),
					resource.TestCheckOutput("is_with_app_trigger_policy_type_set", "true"),
					resource.TestCheckOutput("is_with_app_enabled_set", "true"),
					resource.TestCheckOutput("is_with_com_notification_rules_returned", "true"),
					resource.TestCheckOutput("is_with_com_scope_type_set", "true"),
					resource.TestCheckOutput("is_with_com_scope_components_set", "true"),
					resource.TestCheckOutput("is_with_com_scope_env_and_app_empty", "true"),
				),
			},
		},
	})
}

func testAccDataNotificationRules_base(name string) string {
	return fmt.Sprintf(`
%[1]s

# For all components in the environment.
resource "huaweicloud_cae_notification_rule" "env" {
  count = 2

  name                  = format("%[2]s-env-%%d", count.index)
  event_name            = "Started,Healthy"
  enterprise_project_id = try(
    data.huaweicloud_cae_environments.test[count.index].environments[0].annotations.enterprise_project_id,
    null)

  scope {
    type         = "environments"
    environments = data.huaweicloud_cae_environments.test[count.index].environments[*].id
  }

  trigger_policy {
    type     = "accumulative"
    period   = 300
    count    = 50
    operator = ">="
  }

  notification {
    protocol = "email"
    endpoint = "terraform@test.com"
    template = "EN"
  }

  enabled = true
}

# For all components in the application.
resource "huaweicloud_cae_notification_rule" "app" {
  count = 2

  name                  = format("%[2]s-app-%%d", count.index)
  event_name            = "Unhealthy"
  enterprise_project_id = try(
    data.huaweicloud_cae_environments.test[count.index].environments[0].annotations.enterprise_project_id,
    null)

  scope {
    type         = "applications"
    applications = slice(huaweicloud_cae_application.test[*].id, count.index, count.index+1)
  }

  trigger_policy {
    type = "immediately"
  }

  notification {
    protocol = "sms"
    endpoint = "12345678987"
    template = "ZH"
  }
}

# For specified components.
resource "huaweicloud_cae_notification_rule" "com" {
  count = 2

  name                  = format("%[2]s-com-%%d", count.index)
  event_name            = "Started,BackOffStart,SuccessfulMountVolume"
  enterprise_project_id = try(
    data.huaweicloud_cae_environments.test[count.index].environments[0].annotations.enterprise_project_id,
    null)

  scope {
    type       = "components"
    components = slice(huaweicloud_cae_component.test[*].id, count.index, count.index+1)
  }

  trigger_policy {
    type     = "accumulative"
    period   = 300
    count    = 100
    operator = ">"
  }

  notification {
    protocol = "email"
    endpoint = "terraform@test.com"
    template = "ZH"
  }
}
`, testAccDataSourceComponents_base(name), name)
}

func testAccDataNotificationRules_basic(name string) string {
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
    for v in data.huaweicloud_cae_notification_rules.test.rules : v if strcontains(v.name, "%[2]s-env-")
  ], [])
}

output "is_with_env_notification_rules_returned" {
  value = length(local.with_env_filter_result) >= 2
}

output "is_with_env_name_set" {
  value = alltrue([for v in local.with_env_filter_result: v.name != ""])
}

output "is_with_env_event_name_set" {
  value = alltrue([for v in local.with_env_filter_result: v.event_name != ""])
}

output "is_with_env_scope_environments_set" {
  value = alltrue([for v in local.with_env_filter_result: length(try(v.scope[0].environments, [])) > 0])
}

output "is_with_env_scope_app_and_com_empty" {
  value = alltrue([for v in local.with_env_filter_result: length(try(v.scope[0].applications, [])) == 0 &&
    length(try(v.scope[0].components, [])) == 0])
}

output "is_with_env_trigger_policy_type_set" {
  value = alltrue([for v in local.with_env_filter_result: try(v.trigger_policy[0].type, "") != ""])
}

output "is_with_env_trigger_policy_period_set" {
  value = alltrue([for v in local.with_env_filter_result: try(v.trigger_policy[0].period, 0) > 0])
}

output "is_with_env_trigger_policy_count_set" {
  value = alltrue([for v in local.with_env_filter_result: try(v.trigger_policy[0].count, 0) > 0])
}

output "is_with_env_trigger_policy_operator_set" {
  value = alltrue([for v in local.with_env_filter_result: try(v.trigger_policy[0].operator, "") != ""])
}

output "is_with_env_scope_type_set" {
  value = alltrue([for v in local.with_env_filter_result: try(v.scope[0].type, "") != ""])
}

output "is_with_env_notification_endpoint_set" {
  value = alltrue([for v in local.with_env_filter_result: try(v.notification[0].endpoint, "") != ""])
}

output "is_with_env_notification_protocol_set" {
  value = alltrue([for v in local.with_env_filter_result: try(v.notification[0].protocol, "") != ""])
}

output "is_with_env_notification_template_set" {
  value = alltrue([for v in local.with_env_filter_result: try(v.notification[0].template, "") != ""])
}

output "is_with_env_enabled_set" {
  value = alltrue([for v in local.with_env_filter_result: v.enabled != null])
}

locals {
  with_app_filter_result = try([
    for v in data.huaweicloud_cae_notification_rules.test.rules : v if strcontains(v.name, "%[2]s-app-")
  ], [])
}

output "is_with_app_notification_rules_returned" {
  value = length(local.with_app_filter_result) >= 2
}

output "is_with_app_name_set" {
  value = alltrue([for v in local.with_app_filter_result: v.name != ""])
}

output "is_with_app_scope_type_set" {
  value = alltrue([for v in local.with_app_filter_result: try(v.scope[0].type, "") != ""])
}

output "is_with_app_scope_applications_set" {
  value = alltrue([for v in local.with_app_filter_result: length(try(v.scope[0].applications, [])) > 0])
}

output "is_with_app_scope_env_and_com_empty" {
  value = alltrue([for v in local.with_app_filter_result: length(try(v.scope[0].environments, [])) == 0 &&
    length(try(v.scope[0].components, [])) == 0])
}

output "is_with_app_trigger_policy_type_set" {
  value = alltrue([for v in local.with_app_filter_result: try(v.trigger_policy[0].type, "") != ""])
}

output "is_with_app_enabled_set" {
  value = alltrue([for v in local.with_app_filter_result: v.enabled != null])
}

locals {
  with_com_filter_result = try([
    for v in data.huaweicloud_cae_notification_rules.test.rules : v if strcontains(v.name, "%[2]s-com-")
  ], [])
}

output "is_with_com_notification_rules_returned" {
  value = length(local.with_com_filter_result) >= 2
}

output "is_with_com_scope_type_set" {
  value = alltrue([for v in local.with_com_filter_result: try(v.scope[0].type, "") != ""])
}

output "is_with_com_scope_components_set" {
  value = alltrue([for v in local.with_com_filter_result: length(try(v.scope[0].components, [])) > 0])
}

output "is_with_com_scope_env_and_app_empty" {
  value = alltrue([for v in local.with_com_filter_result: length(try(v.scope[0].environments, [])) == 0 &&
    length(try(v.scope[0].applications, [])) == 0])
}
`, testAccDataNotificationRules_base(name), name)
}
