package codeartspipeline

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCodeartsPipelineAvailablePlugins_basic(t *testing.T) {
	dataSource := "data.huaweicloud_codearts_pipeline_available_plugins.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCodeartsPipelineAvailablePlugins_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.unique_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.display_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.business_type"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.removable"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.cloneable"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.disabled"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.editable"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.plugins_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.plugins_list.0.unique_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.plugins_list.0.display_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.plugins_list.0.plugin_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.plugins_list.0.disabled"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.plugins_list.0.plugin_attribution"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.plugins_list.0.plugin_composition_type"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.plugins_list.0.runtime_attribution"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.plugins_list.0.version_attribution"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.plugins_list.0.multi_step_editable"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.plugins_list.0.location"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.plugins_list.0.manifest_version"),
				),
			},
		},
	})
}

const testDataSourceCodeartsPipelineAvailablePlugins_basic = `
data "huaweicloud_codearts_pipeline_available_plugins" "test" {
  use_condition = "pipeline"
}
`
