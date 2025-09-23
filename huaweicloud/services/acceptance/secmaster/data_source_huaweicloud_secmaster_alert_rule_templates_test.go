package secmaster

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSecmasterAlertRuleTemplates_basic(t *testing.T) {
	dataSource := "data.huaweicloud_secmaster_alert_rule_templates.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSecmasterAlertRuleTemplates_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "templates.#"),
					resource.TestCheckResourceAttrSet(dataSource, "templates.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "templates.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "templates.0.severity"),
					resource.TestMatchResourceAttr(dataSource, "templates.0.updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),

					resource.TestCheckOutput("is_severity_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceSecmasterAlertRuleTemplates_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_alert_rule_templates" "test" {
  workspace_id = "%[1]s"
}

locals {
  severity = data.huaweicloud_secmaster_alert_rule_templates.test.templates[0].severity
}

data "huaweicloud_secmaster_alert_rule_templates" "filter_by_severity" {
  workspace_id = "%[1]s"
  severity     = [local.severity]
}

output "is_severity_filter_useful" {
  value = length(data.huaweicloud_secmaster_alert_rule_templates.filter_by_severity.templates) > 0 && alltrue(
    [for v in data.huaweicloud_secmaster_alert_rule_templates.filter_by_severity.templates[*].severity : v == local.severity]
  )
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID)
}
