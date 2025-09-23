package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAlertRuleTemplateDetail_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_secmaster_alert_rule_template_detail.test"
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
				Config: testAccDataSourceAlertRuleTemplateDetail_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "rule_template_id"),
					resource.TestCheckResourceAttrSet(dataSource, "template_name"),
					resource.TestCheckResourceAttrSet(dataSource, "query"),
					resource.TestCheckResourceAttrSet(dataSource, "query_type"),
					resource.TestCheckResourceAttrSet(dataSource, "severity"),
					resource.TestCheckResourceAttrSet(dataSource, "event_grouping"),
					resource.TestCheckResourceAttrSet(dataSource, "schedule.#"),
					resource.TestCheckResourceAttrSet(dataSource, "schedule.0.frequency_interval"),
					resource.TestCheckResourceAttrSet(dataSource, "schedule.0.frequency_unit"),
					resource.TestCheckResourceAttrSet(dataSource, "schedule.0.period_interval"),
					resource.TestCheckResourceAttrSet(dataSource, "schedule.0.period_unit"),
					resource.TestCheckResourceAttrSet(dataSource, "schedule.0.delay_interval"),
					resource.TestCheckResourceAttrSet(dataSource, "schedule.0.overtime_interval"),
					resource.TestCheckResourceAttrSet(dataSource, "triggers.#"),
					resource.TestCheckResourceAttrSet(dataSource, "triggers.0.mode"),
					resource.TestCheckResourceAttrSet(dataSource, "triggers.0.operator"),
					resource.TestCheckResourceAttrSet(dataSource, "triggers.0.expression"),
					resource.TestCheckResourceAttrSet(dataSource, "triggers.0.severity"),
					resource.TestCheckResourceAttrSet(dataSource, "triggers.0.accumulated_times"),
					resource.TestCheckResourceAttrSet(dataSource, "update_time"),
				),
			},
		},
	})
}

func testAccDataSourceAlertRuleTemplateDetail_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_alert_rule_template_detail" "test" {
  workspace_id = "%[1]s"
  template_id  = "%[2]s"
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, acceptance.HW_SECMASTER_TEMPLATE_ID)
}
