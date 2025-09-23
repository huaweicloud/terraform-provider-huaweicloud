package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSiemdirectories_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_secmaster_siem_directories.test"
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
				Config: testAccDataSourceSiemdirectories_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "workspaceid"),
					resource.TestCheckResourceAttrSet(dataSource, "project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "directories.#"),
					resource.TestCheckResourceAttrSet(dataSource, "directory_i18ns.#"),
				),
			},
		},
	})
}

func testAccDataSourceSiemdirectories_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_siem_directories" "test" {
  workspace_id = "%[1]s"
  category     = "ALERT_RULE"
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID)
}
