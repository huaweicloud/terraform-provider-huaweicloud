package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSocAttachment_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_secmaster_soc_attachment.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
			acceptance.TestAccPreCheckSecMasterAttachID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSocAttachment_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "data.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data.0.file_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data.0.file_folder"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data.0.workspace_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data.0.storage_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data.0.update_time"),
				),
			},
		},
	})
}

func testDataSourceSocAttachment_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_soc_attachment" "test" {
  workspace_id = "%[1]s"
  attach_id    = "%[2]s"
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, acceptance.HW_SECMASTER_ATTACH_ID)
}
