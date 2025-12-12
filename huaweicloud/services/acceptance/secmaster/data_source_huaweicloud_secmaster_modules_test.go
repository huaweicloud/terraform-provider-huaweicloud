package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceModules_basic(t *testing.T) {
	dataSource := "data.huaweicloud_secmaster_modules.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSocModules_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.en_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.is_built_in"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.module_json"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.module_type"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.version"),
				),
			},
		},
	})
}

func testAccDataSourceSocModules_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_modules" "test" {
  workspace_id = "%[1]s"
  module_type  = "section"
  sort_key     = "creat_time"
  sort_dir     = "ASC"
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID)
}
