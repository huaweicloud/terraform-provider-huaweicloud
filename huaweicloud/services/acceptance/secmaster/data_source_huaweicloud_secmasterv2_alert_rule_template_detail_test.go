package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAlertRuleTemplateDetailV2_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_secmasterv2_alert_rule_template_detail.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
			// The alert rule template ID
			acceptance.TestAccPreCheckSecMasterTemplateId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAlertRuleTemplateDetailV2_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "rule_template_id"),
					resource.TestCheckResourceAttrSet(dataSource, "template_name"),
					resource.TestCheckResourceAttrSet(dataSource, "accumulated_times"),
					resource.TestCheckResourceAttrSet(dataSource, "cu_quota_amount"),
					resource.TestCheckResourceAttrSet(dataSource, "description"),
					resource.TestCheckResourceAttrSet(dataSource, "environment"),
					resource.TestCheckResourceAttrSet(dataSource, "job_mode"),
					resource.TestCheckResourceAttrSet(dataSource, "job_mode_setting.#"),
					resource.TestCheckResourceAttrSet(dataSource, "job_mode_setting.0.search_delay_interval"),
					resource.TestCheckResourceAttrSet(dataSource, "job_mode_setting.0.search_delay_unit"),
					resource.TestCheckResourceAttrSet(dataSource, "job_mode_setting.0.search_frequency_interval"),
					resource.TestCheckResourceAttrSet(dataSource, "job_mode_setting.0.search_frequency_unit"),
					resource.TestCheckResourceAttrSet(dataSource, "job_mode_setting.0.search_period_interval"),
					resource.TestCheckResourceAttrSet(dataSource, "job_mode_setting.0.search_period_unit"),
					resource.TestCheckResourceAttrSet(dataSource, "job_output_setting.#"),
					resource.TestCheckResourceAttrSet(dataSource, "job_output_setting.0.alert_description"),
					resource.TestCheckResourceAttrSet(dataSource, "job_output_setting.0.alert_grouping"),
					resource.TestCheckResourceAttrSet(dataSource, "job_output_setting.0.alert_name"),
					resource.TestCheckResourceAttrSet(dataSource, "job_output_setting.0.alert_remediation"),
					resource.TestCheckResourceAttrSet(dataSource, "job_output_setting.0.alert_severity"),
					resource.TestCheckResourceAttrSet(dataSource, "job_output_setting.0.alert_suppression"),
					resource.TestCheckResourceAttrSet(dataSource, "process_status"),
					resource.TestCheckResourceAttrSet(dataSource, "query_type"),
					resource.TestCheckResourceAttrSet(dataSource, "script"),
					resource.TestCheckResourceAttrSet(dataSource, "status"),
					resource.TestCheckResourceAttrSet(dataSource, "table_name"),
					resource.TestCheckResourceAttrSet(dataSource, "create_by"),
					resource.TestCheckResourceAttrSet(dataSource, "create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "update_by"),
					resource.TestCheckResourceAttrSet(dataSource, "update_time"),
					resource.TestCheckResourceAttrSet(dataSource, "triggers.#"),
					resource.TestCheckResourceAttrSet(dataSource, "triggers.0.accumulated_times"),
					resource.TestCheckResourceAttrSet(dataSource, "triggers.0.expression"),
					resource.TestCheckResourceAttrSet(dataSource, "triggers.0.mode"),
					resource.TestCheckResourceAttrSet(dataSource, "triggers.0.operator"),
					resource.TestCheckResourceAttrSet(dataSource, "triggers.0.severity"),
				),
			},
		},
	})
}

func testAccDataSourceAlertRuleTemplateDetailV2_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_secmasterv2_alert_rule_template_detail" "test" {
  workspace_id = "%[1]s"
  template_id  = "%[2]s"
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, acceptance.HW_SECMASTER_TEMPLATE_ID)
}
