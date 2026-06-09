package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceComponentHistoryConfiguration_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_secmaster_component_history_configuration.test"
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
				Config: testDataSourceComponentHistoryConfiguration_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "records.#"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.node_id"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.node_name"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.specification"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.list.0.configuration_id"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.list.0.file_name"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.list.0.file_type"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.list.0.node_id"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.list.0.param"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.list.0.type"),
				),
			},
		},
	})
}

func testDataSourceComponentHistoryConfiguration_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_component_history_configuration" "test" {
  workspace_id = "%[1]s"
  component_id = "%[2]s"
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, acceptance.HW_SECMASTER_COMPONENT_ID)
}
