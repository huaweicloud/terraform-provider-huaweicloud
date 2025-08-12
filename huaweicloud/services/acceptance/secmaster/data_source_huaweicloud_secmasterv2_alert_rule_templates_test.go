package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAlertRuleTemplatesV2_basic(t *testing.T) {
	dataSource := "data.huaweicloud_secmasterv2_alert_rule_templates.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAlertRuleTemplatesV2_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "records.#"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.template_id"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.template_name"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.accumulated_times"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.alert_description"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.alert_name"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.alert_remediation"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.create_by"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.event_grouping"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.job_mode"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.process_status"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.query"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.query_type"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.severity"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.simulation"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.suppression"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.table_name"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.update_by"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.update_time"),

					resource.TestCheckOutput("is_template_name_filter_useful", "true"),
					resource.TestCheckOutput("is_severity_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceAlertRuleTemplatesV2_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_secmasterv2_alert_rule_templates" "test" {
  workspace_id = "%[1]s"
}

locals {
  template_name = data.huaweicloud_secmasterv2_alert_rule_templates.test.records[0].template_name
}

data "huaweicloud_secmasterv2_alert_rule_templates" "filter_by_template_name" {
  workspace_id  = "%[1]s"
  template_name = local.template_name
}

output "is_template_name_filter_useful" {
  value = length(data.huaweicloud_secmasterv2_alert_rule_templates.filter_by_template_name.records) > 0 && alltrue(
    [for v in data.huaweicloud_secmasterv2_alert_rule_templates.filter_by_template_name.records[*].template_name : v == local.template_name]
  )
}

locals {
  severity = data.huaweicloud_secmasterv2_alert_rule_templates.test.records[0].severity
}

data "huaweicloud_secmasterv2_alert_rule_templates" "filter_by_severity" {
  workspace_id = "%[1]s"
  severity     = local.severity
}

output "is_severity_filter_useful" {
  value = length(data.huaweicloud_secmasterv2_alert_rule_templates.filter_by_severity.records) > 0 && alltrue(
    [for v in data.huaweicloud_secmasterv2_alert_rule_templates.filter_by_severity.records[*].severity : v == local.severity]
  )
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID)
}
