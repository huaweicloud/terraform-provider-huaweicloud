package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceLayoutWizards_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_secmaster_layout_wizards.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
			acceptance.TestAccPreCheckSecMasterLayoutID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceLayoutWizards_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.wizard_json"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.creator_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.update_time"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.workspace_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.is_binding"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.binding_button.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.binding_button.0.button_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.binding_button.0.button_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.version"),
				),
			},
		},
	})
}

func testDataSourceLayoutWizards_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_layout_wizards" "test" {
  workspace_id = "%[1]s"
  layout_id    = "%[2]s"
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, acceptance.HW_SECMASTER_LAYOUT_ID)
}
