package secmaster

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSecmasterAlertRules_basic(t *testing.T) {
	dataSource := "data.huaweicloud_secmaster_alert_rules.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
			acceptance.TestAccPreCheckSecMasterPipelineID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSecmasterAlertRules_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "alert_rules.#"),
					resource.TestCheckResourceAttrSet(dataSource, "alert_rules.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "alert_rules.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "alert_rules.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "alert_rules.0.severity"),
					resource.TestCheckResourceAttrSet(dataSource, "alert_rules.0.pipeline_id"),
					resource.TestMatchResourceAttr(dataSource, "alert_rules.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(dataSource, "alert_rules.0.updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),

					resource.TestCheckOutput("is_id_filter_useful", "true"),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
					resource.TestCheckOutput("is_severity_filter_useful", "true"),
					resource.TestCheckOutput("is_pipeline_id_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceSecmasterAlertRules_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_secmaster_alert_rules" "test" {
  workspace_id = "%[2]s"
  
  depends_on = [
    huaweicloud_secmaster_alert_rule.test1,
    huaweicloud_secmaster_alert_rule.test2,
  ]
}

locals {
  id          = data.huaweicloud_secmaster_alert_rules.test.alert_rules[0].id
  name        = data.huaweicloud_secmaster_alert_rules.test.alert_rules[0].name
  status      = data.huaweicloud_secmaster_alert_rules.test.alert_rules[0].status
  severity    = data.huaweicloud_secmaster_alert_rules.test.alert_rules[0].severity
  pipeline_id = data.huaweicloud_secmaster_alert_rules.test.alert_rules[0].pipeline_id
}

data "huaweicloud_secmaster_alert_rules" "filter_by_id" {
  workspace_id = "%[2]s"
  rule_id      = local.id
}

data "huaweicloud_secmaster_alert_rules" "filter_by_name" {
  workspace_id = "%[2]s"
  name         = local.name
}

data "huaweicloud_secmaster_alert_rules" "filter_by_status" {
  workspace_id = "%[2]s"
  status       = [local.status]
}

data "huaweicloud_secmaster_alert_rules" "filter_by_severity" {
  workspace_id = "%[2]s"
  severity     = [local.severity]
}

data "huaweicloud_secmaster_alert_rules" "filter_by_pipeline_id" {
  workspace_id = "%[2]s"
  pipeline_id  = local.pipeline_id
}

output "is_id_filter_useful" {
  value = length(data.huaweicloud_secmaster_alert_rules.filter_by_id.alert_rules) > 0 && alltrue(
    [for v in data.huaweicloud_secmaster_alert_rules.filter_by_id.alert_rules[*].id : v == local.id]
  )
}

output "is_name_filter_useful" {
  value = length(data.huaweicloud_secmaster_alert_rules.filter_by_name.alert_rules) > 0 && alltrue(
    [for v in data.huaweicloud_secmaster_alert_rules.filter_by_name.alert_rules[*].name : strcontains(v, local.name)]
  )
}

output "is_status_filter_useful" {
  value = length(data.huaweicloud_secmaster_alert_rules.filter_by_status.alert_rules) > 0 && alltrue(
    [for v in data.huaweicloud_secmaster_alert_rules.filter_by_status.alert_rules[*].status : v == local.status]
  )
}

output "is_severity_filter_useful" {
  value = length(data.huaweicloud_secmaster_alert_rules.filter_by_severity.alert_rules) > 0 && alltrue(
    [for v in data.huaweicloud_secmaster_alert_rules.filter_by_severity.alert_rules[*].severity : v == local.severity]
  )
}

output "is_pipeline_id_filter_useful" {
  value = length(data.huaweicloud_secmaster_alert_rules.filter_by_pipeline_id.alert_rules) > 0 && alltrue(
    [for v in data.huaweicloud_secmaster_alert_rules.filter_by_pipeline_id.alert_rules[*].pipeline_id : v == local.pipeline_id]
  )
}

`, testAlertRule_buildData(name), acceptance.HW_SECMASTER_WORKSPACE_ID)
}

func testAlertRule_buildData(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_secmaster_alert_rule" "test1" {
  workspace_id = "%[1]s"
  pipeline_id  = "%[2]s"
  name         = "%[3]s"
  description  = "this is a test rule created by terraform"
  status       = "ENABLED"
  severity     = "TIPS"
  type = {
    "name"     = "DNS protocol attacks"
    "category" = "DDoS attacks"
  }

  triggers {
    mode              = "COUNT"
    operator          = "GT"
    expression        = 5
    severity          = "MEDIUM"
    accumulated_times = 1
  }

  query_rule = "* | select status, count(*) as count group by status"
  query_type = "SQL"
  query_plan {
    query_interval      = 1
    query_interval_unit = "HOUR"
    time_window         = 1
    time_window_unit    = "HOUR"
  }
}

resource "huaweicloud_secmaster_alert_rule" "test2" {
  workspace_id = "%[1]s"
  pipeline_id  = "%[2]s"
  name         = "%[3]s_2"
  description  = "this is a test rule created by terraform"
  status       = "DISABLED"
  severity     = "MEDIUM"
  type = {
    "name"     = "DNS protocol attacks"
    "category" = "DDoS attacks"
  }

  triggers {
    mode              = "COUNT"
    operator          = "EQ"
    expression        = 10
    severity          = "MEDIUM"
    accumulated_times = 1
  }

  query_rule = "* | select status, count(*) as count group by status"
  query_type = "SQL"
  query_plan {
    query_interval      = 5
    query_interval_unit = "MINUTE"
    time_window         = 5
    time_window_unit    = "MINUTE"
  }
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, acceptance.HW_SECMASTER_PIPELINE_ID, name)
}
