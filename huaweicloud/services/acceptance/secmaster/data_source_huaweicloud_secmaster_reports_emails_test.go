package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceReportsEmails_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_secmaster_reports_emails.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceReportsEmails_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "emails.#"),
					resource.TestCheckResourceAttrSet(dataSource, "emails.0.report_address"),
					resource.TestCheckResourceAttrSet(dataSource, "emails.0.email_status"),
					resource.TestCheckResourceAttrSet(dataSource, "emails.1.report_address"),
					resource.TestCheckResourceAttrSet(dataSource, "emails.1.email_status"),
				),
			},
		},
	})
}

func testDataSourceReportsEmails_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_reports_emails" "test" {
  workspace_id  = "%[1]s"
  email_address = "terraform@outlook.com;test@163.com"
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID)
}
