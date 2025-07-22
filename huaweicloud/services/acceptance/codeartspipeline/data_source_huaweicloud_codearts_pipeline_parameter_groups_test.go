package codeartspipeline

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCodeartsPipelineParameterGroups_basic(t *testing.T) {
	dataSource := "data.huaweicloud_codearts_pipeline_parameter_groups.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCodeartsPipelineParameterGroups_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "groups.#"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.creator_id"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.updater_id"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.creator_name"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.updater_name"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.update_time"),

					resource.TestCheckOutput("is_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceCodeartsPipelineParameterGroups_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_codearts_pipeline_parameter_groups" "test" {
  depends_on = [huaweicloud_codearts_pipeline_parameter_group.test]

  project_id = huaweicloud_codearts_project.test.id
}

// filter by name
data "huaweicloud_codearts_pipeline_parameter_groups" "filter_by_name" {
  project_id = huaweicloud_codearts_project.test.id
  name       = huaweicloud_codearts_pipeline_parameter_group.test.name
}

locals {
  filter_result_by_name = [for v in data.huaweicloud_codearts_pipeline_parameter_groups.filter_by_name.groups[*].name :
    v == huaweicloud_codearts_pipeline_parameter_group.test.name]
}

output "is_name_filter_useful" {
  value = length(local.filter_result_by_name) > 0 && alltrue(local.filter_result_by_name)
}
`, testPipelineParameterGroup_basic(name))
}
