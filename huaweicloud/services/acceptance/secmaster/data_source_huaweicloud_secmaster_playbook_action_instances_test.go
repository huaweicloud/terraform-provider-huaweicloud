package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSecmasterPlaybookActionInstances_basic(t *testing.T) {
	dataSource := "data.huaweicloud_secmaster_playbook_action_instances.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSecmasterPlaybookActionInstances_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "action_instances.#"),
					resource.TestCheckResourceAttrSet(dataSource, "action_instances.0.action.#"),
					resource.TestCheckResourceAttrSet(dataSource, "action_instances.0.action.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "action_instances.0.action.0.action_id"),
					resource.TestCheckResourceAttrSet(dataSource, "action_instances.0.action.0.playbook_version_id"),
					resource.TestCheckResourceAttrSet(dataSource, "action_instances.0.instance_log.#"),
				),
			},
		},
	})
}

func testDataSourceSecmasterPlaybookActionInstances_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_playbook_instances" "test" {
  workspace_id = "%[1]s"
}

data "huaweicloud_secmaster_playbook_action_instances" "test" {
  workspace_id         = "%[1]s"
  playbook_instance_id = data.huaweicloud_secmaster_playbook_instances.test.instances[0].id
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID)
}
