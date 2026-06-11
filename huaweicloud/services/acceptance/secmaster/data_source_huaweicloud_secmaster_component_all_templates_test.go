package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceComponentAllTemplates_basic(t *testing.T) {
	dataSource := "data.huaweicloud_secmaster_component_all_templates.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Before running test, prepare component templates in SecMaster.
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceComponentAllTemplates_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.component_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.template_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.task_config"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.update_time"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.project_id"),

					resource.TestCheckOutput("is_search_value_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceComponentAllTemplates_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_component_all_templates" "test" {
  workspace_id = "%[1]s"
}

# Filter by search_value
locals {
  search_value = data.huaweicloud_secmaster_component_all_templates.test.data[0].template_name
}

data "huaweicloud_secmaster_component_all_templates" "filter_by_search_value" {
  workspace_id  = "%[1]s"
  search_value  = local.search_value
}

output "is_search_value_filter_useful" {
  value = length(data.huaweicloud_secmaster_component_all_templates.filter_by_search_value.data) > 0 && alltrue(
    [for v in data.huaweicloud_secmaster_component_all_templates.filter_by_search_value.data[*].template_name :
    v == local.search_value]
  )
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID)
}
