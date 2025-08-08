package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceComponentTemplates_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_secmaster_component_templates.test"
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
				Config: testDataSourceComponentTemplates_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "records.#"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.version"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.component_id"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.component_name"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.param"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.file_type"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.file_name"),

					resource.TestCheckOutput("is_file_type_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceComponentTemplates_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_component_templates" "test" {
  workspace_id = "%[1]s"
  component_id = "%[2]s"
}

# Filter using file_type.
locals {
  file_type = data.huaweicloud_secmaster_component_templates.test.records[0].file_type
}

data "huaweicloud_secmaster_component_templates" "file_type_filter" {
  workspace_id = "%[1]s"
  component_id = "%[2]s"
  file_type    = local.file_type
}

output "is_file_type_filter_useful" {
  value = length(data.huaweicloud_secmaster_component_templates.file_type_filter.records) > 0 && alltrue(
    [for v in data.huaweicloud_secmaster_component_templates.file_type_filter.records[*].file_type : v == local.file_type]
  )
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, acceptance.HW_SECMASTER_COMPONENT_ID)
}
