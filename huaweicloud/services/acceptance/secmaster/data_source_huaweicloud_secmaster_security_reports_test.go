package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSecurityReports_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_secmaster_security_reports.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Before running the test, prepare a security report in the workspace.
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSecurityReports_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "reports.#"),
					resource.TestCheckResourceAttrSet(dataSource, "reports.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "reports.0.report_name"),
					resource.TestCheckResourceAttrSet(dataSource, "reports.0.report_period"),
					resource.TestCheckResourceAttrSet(dataSource, "reports.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "reports.0.language"),
					resource.TestCheckResourceAttrSet(dataSource, "reports.0.report_range.#"),
					resource.TestCheckResourceAttrSet(dataSource, "reports.0.report_range.0.start"),
					resource.TestCheckResourceAttrSet(dataSource, "reports.0.report_range.0.end"),
					resource.TestCheckResourceAttrSet(dataSource, "reports.0.report_rule_infos.#"),
					resource.TestCheckResourceAttrSet(dataSource, "reports.0.report_rule_infos.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "reports.0.report_rule_infos.0.project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "reports.0.report_rule_infos.0.workspace_id"),
					resource.TestCheckResourceAttrSet(dataSource, "reports.0.report_rule_infos.0.rule"),
					resource.TestCheckResourceAttrSet(dataSource, "reports.0.report_rule_infos.0.start_time"),
					resource.TestCheckResourceAttrSet(dataSource, "reports.0.report_rule_infos.0.end_time"),
					resource.TestCheckResourceAttrSet(dataSource, "reports.0.report_rule_infos.0.email_title"),
					resource.TestCheckResourceAttrSet(dataSource, "reports.0.report_rule_infos.0.email_content"),
					resource.TestCheckResourceAttrSet(dataSource, "reports.0.report_rule_infos.0.email_to"),
				),
			},
		},
	})
}

func testDataSourceSecurityReports_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_security_reports" "test" {
  workspace_id  = "%[1]s"
  report_period = "daily"
  status        = "enable"
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID)
}
