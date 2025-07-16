package codeartspipeline

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourcePipelineModules_basic(t *testing.T) {
	dataSource := "data.huaweicloud_codearts_pipeline_modules.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testDataSourcePipelineModules_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "modules.#"),
					resource.TestCheckResourceAttrSet(dataSource, "modules.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "modules.0.base_url"),
					resource.TestCheckResourceAttrSet(dataSource, "modules.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "modules.0.location"),
					resource.TestCheckResourceAttrSet(dataSource, "modules.0.module_id"),
					resource.TestCheckResourceAttrSet(dataSource, "modules.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "modules.0.properties"),
					resource.TestCheckResourceAttrSet(dataSource, "modules.0.publisher"),
					resource.TestCheckResourceAttrSet(dataSource, "modules.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "modules.0.version"),
					resource.TestCheckResourceAttrSet(dataSource, "modules.0.tags.#"),
					resource.TestCheckResourceAttrSet(dataSource, "modules.0.url_relative"),
					resource.TestCheckResourceAttrSet(dataSource, "modules.0.manifest_version"),

					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("is_tags_filter_useful", "true"),
				),
			},
		},
	})
}

const testDataSourcePipelineModules_basic = `
data "huaweicloud_codearts_pipeline_modules" "test" {}

locals {
  name = data.huaweicloud_codearts_pipeline_modules.test.modules[0].name
  tag  = data.huaweicloud_codearts_pipeline_modules.test.modules[0].tags[0]
}

// filter by name
data "huaweicloud_codearts_pipeline_modules" "filter_by_name" {
  name = local.name
}

locals {
  filter_result_by_name = [for v in data.huaweicloud_codearts_pipeline_modules.filter_by_name.modules[*].name :
    v == local.name]
}

output "is_name_filter_useful" {
  value = length(local.filter_result_by_name) == 1 && alltrue(local.filter_result_by_name)
}

// filter by tags
data "huaweicloud_codearts_pipeline_modules" "filter_by_tags" {
  tags = [local.tag]
}

locals {
  filter_result_by_tags = [for v in data.huaweicloud_codearts_pipeline_modules.filter_by_tags.modules[*].tags :
    contains(v, local.tag)]
}

output "is_tags_filter_useful" {
  value = length(local.filter_result_by_tags) > 0 && alltrue(local.filter_result_by_name)
}
`
