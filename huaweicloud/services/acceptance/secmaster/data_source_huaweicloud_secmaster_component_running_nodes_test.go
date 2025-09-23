package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Due to testing environment limitations, this test case can only test the scenario with empty `records`.
func TestAccDataSourceComponentRunningNodes_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_secmaster_component_running_nodes.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
			acceptance.TestAccPreCheckSecMasterComponentId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceComponentRunningNodes_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "records.#"),
				),
			},
		},
	})
}

func testDataSourceComponentRunningNodes_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_component_running_nodes" "test" {
  workspace_id = "%[1]s"
  component_id = "%[2]s"
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, acceptance.HW_SECMASTER_COMPONENT_ID)
}
