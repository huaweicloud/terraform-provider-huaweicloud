package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCollectorModuleTemplates_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_secmaster_collector_module_templates.test"
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
				Config: testAccDataSourceCollectorModuleTemplates_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "common.#"),
					resource.TestCheckResourceAttrSet(dataSource, "common.0.template_id"),
					resource.TestCheckResourceAttrSet(dataSource, "common.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "common.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "common.0.title"),
					resource.TestCheckResourceAttrSet(dataSource, "list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "list.0.template_id"),
					resource.TestCheckResourceAttrSet(dataSource, "list.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "list.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "list.0.title"),
				),
			},
		},
	})
}

func testAccDataSourceCollectorModuleTemplates_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_collector_module_templates" "test" {
  workspace_id = "%s"
  parser_type  = "FILTER"
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID)
}
