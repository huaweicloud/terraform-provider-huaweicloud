package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourcePlaybookActions_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_secmaster_playbook_actions.test"
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
				Config: testDataSourcePlaybookActions_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.action_type"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.action_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.playbook_version_id"),
				),
			},
		},
	})
}

func testDataSourcePlaybookActions_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_playbook_actions" "test" {
  workspace_id = "%[1]s"
  version_id   = "%[2]s"
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, acceptance.HW_SECMASTER_VERSION_ID)
}
