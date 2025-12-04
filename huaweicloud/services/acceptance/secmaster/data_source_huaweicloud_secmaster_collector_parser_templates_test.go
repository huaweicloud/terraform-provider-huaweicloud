package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSecmasterCollectorParserTemplates_basic(t *testing.T) {
	dataSource := "data.huaweicloud_secmaster_collector_parser_templates.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSecmasterCollectorParserTemplates_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.parser_id"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.title"),

					resource.TestCheckOutput("is_title_filter_useful", "true"),
					resource.TestCheckOutput("is_description_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceSecmasterCollectorParserTemplates_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_collector_parser_templates" "test" {
  workspace_id = "%[1]s"
}

# Filter by title
locals {
  title = data.huaweicloud_secmaster_collector_parser_templates.test.records[0].title
}

data "huaweicloud_secmaster_collector_parser_templates" "filter_by_title" {
  workspace_id = "%[1]s"
  title        = local.title
}

locals {
  list_by_title = data.huaweicloud_secmaster_collector_parser_templates.filter_by_title.records
}

output "is_title_filter_useful" {
  value = length(local.list_by_title) > 0 && alltrue(
    [for v in local.list_by_title[*].title : v == local.title]
  )
}

# Filter by description
locals {
  description = data.huaweicloud_secmaster_collector_parser_templates.test.records[0].description
}

data "huaweicloud_secmaster_collector_parser_templates" "filter_by_description" {
  workspace_id   = "%[1]s"
  description    = local.description
}

locals {
  list_by_description = data.huaweicloud_secmaster_collector_parser_templates.filter_by_description.records
}

output "is_description_filter_useful" {
  value = length(local.list_by_description) > 0 && alltrue(
    [for v in local.list_by_description[*].description : v == local.description]
  )
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID)
}
