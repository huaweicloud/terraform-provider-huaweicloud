package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourcePlaybookApprovals_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_secmaster_playbook_approvals.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
			acceptance.TestAccPreCheckSecMasterVersionId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourcePlaybookApprovals_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.result"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.content"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.resource_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.user_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.update_time"),
				),
			},
		},
	})
}

func testDataSourcePlaybookApprovals_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_playbook_approvals" "test" {
  workspace_id  = "%[1]s"
  resource_id   = "%[2]s"
  approve_type  = "PLAYBOOK"
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, acceptance.HW_SECMASTER_VERSION_ID)
}
