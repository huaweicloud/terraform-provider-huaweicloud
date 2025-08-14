package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSocComponentActionDetail_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_secmaster_soc_component_action_detail.test"
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
				Config: testDataSourceSocComponentActionDetail_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.action_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.action_desc"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.action_type"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.creator_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.can_update"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.action_version_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.action_version_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.action_version_number"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.action_enable"),
				),
			},
		},
	})
}

func testDataSourceSocComponentActionDetail_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_soc_components" "test" {
  workspace_id = "%[1]s"
}

data "huaweicloud_secmaster_soc_component_actions" "test" {
  workspace_id = "%[1]s"
  component_id = data.huaweicloud_secmaster_soc_components.test.data[0].id
  enabled      = true
}

data "huaweicloud_secmaster_soc_component_action_detail" "test" {
  workspace_id = "%[1]s"
  component_id = data.huaweicloud_secmaster_soc_components.test.data[0].id
  action_id    = data.huaweicloud_secmaster_soc_component_actions.test.data[0].id
  enabled      = true
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID)
}
