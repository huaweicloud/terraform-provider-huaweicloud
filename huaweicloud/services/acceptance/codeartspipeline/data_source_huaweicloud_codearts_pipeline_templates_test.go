package codeartspipeline

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourcePipelineTemplates_basic(t *testing.T) {
	dataSource := "data.huaweicloud_codearts_pipeline_templates.test"
	name := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testDataSourcePipelineTemplates_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "templates.#"),
					resource.TestCheckResourceAttrSet(dataSource, "templates.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "templates.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "templates.0.manifest_version"),
					resource.TestCheckResourceAttrSet(dataSource, "templates.0.language"),
					resource.TestCheckResourceAttrSet(dataSource, "templates.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "templates.0.is_system"),
					resource.TestCheckResourceAttrSet(dataSource, "templates.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "templates.0.update_time"),
					resource.TestCheckResourceAttrSet(dataSource, "templates.0.creator_id"),
					resource.TestCheckResourceAttrSet(dataSource, "templates.0.creator_name"),
					resource.TestCheckResourceAttrSet(dataSource, "templates.0.updater_id"),
					resource.TestCheckResourceAttrSet(dataSource, "templates.0.is_favorite"),
					resource.TestCheckResourceAttrSet(dataSource, "templates.0.is_show_source"),
					resource.TestCheckResourceAttrSet(dataSource, "templates.0.stages.#"),

					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("is_language_filter_useful", "true"),
					resource.TestCheckOutput("is_is_system_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourcePipelineTemplates_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_codearts_pipeline_templates" "test" {
  depends_on = [huaweicloud_codearts_pipeline_template.test]
}

// filter by name
data "huaweicloud_codearts_pipeline_templates" "filter_by_name" {
  name = huaweicloud_codearts_pipeline_template.test.name
}

locals {
  filter_result_by_name = [for v in data.huaweicloud_codearts_pipeline_templates.filter_by_name.templates[*].name :
    v == huaweicloud_codearts_pipeline_template.test.name]
}

output "is_name_filter_useful" {
  value = length(local.filter_result_by_name) > 0 && alltrue(local.filter_result_by_name)
}

// filter by language
data "huaweicloud_codearts_pipeline_templates" "filter_by_language" {
  language = huaweicloud_codearts_pipeline_template.test.language
}

locals {
  filter_result_by_language = [for v in data.huaweicloud_codearts_pipeline_templates.filter_by_language.templates[*].language :
    v == huaweicloud_codearts_pipeline_template.test.language]
}

output "is_language_filter_useful" {
  value = length(local.filter_result_by_language) > 0 && alltrue(local.filter_result_by_language)
}

// filter by is_system
data "huaweicloud_codearts_pipeline_templates" "filter_by_is_system" {
  is_system = huaweicloud_codearts_pipeline_template.test.is_system
}

locals {
  filter_result_by_is_system = [for v in data.huaweicloud_codearts_pipeline_templates.filter_by_is_system.templates[*].is_system :
    v == huaweicloud_codearts_pipeline_template.test.is_system]
}

output "is_is_system_filter_useful" {
  value = length(local.filter_result_by_is_system) > 0 && alltrue(local.filter_result_by_is_system)
}
`, testPipelineTemplate_basic(name))
}
