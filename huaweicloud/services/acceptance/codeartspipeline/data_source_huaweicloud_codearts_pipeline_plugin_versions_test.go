package codeartspipeline

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCodeartsPipelinePluginVersions_basic(t *testing.T) {
	dataSource := "data.huaweicloud_codearts_pipeline_plugin_versions.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCodeartsPipelinePluginVersions_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "versions.#"),
					resource.TestCheckResourceAttrSet(dataSource, "versions.0.plugin_name"),
					resource.TestCheckResourceAttrSet(dataSource, "versions.0.unique_id"),
					resource.TestCheckResourceAttrSet(dataSource, "versions.0.display_name"),
					resource.TestCheckResourceAttrSet(dataSource, "versions.0.version"),
					resource.TestCheckResourceAttrSet(dataSource, "versions.0.version_description"),
					resource.TestCheckResourceAttrSet(dataSource, "versions.0.version_attribution"),
					resource.TestCheckResourceAttrSet(dataSource, "versions.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "versions.0.plugin_attribution"),
					resource.TestCheckResourceAttrSet(dataSource, "versions.0.plugin_composition_type"),
					resource.TestCheckResourceAttrSet(dataSource, "versions.0.icon_url"),
					resource.TestCheckResourceAttrSet(dataSource, "versions.0.workspace_id"),
					resource.TestCheckResourceAttrSet(dataSource, "versions.0.refer_count"),
					resource.TestCheckResourceAttrSet(dataSource, "versions.0.active"),
					resource.TestCheckResourceAttrSet(dataSource, "versions.0.business_type"),
					resource.TestCheckResourceAttrSet(dataSource, "versions.0.business_type_display_name"),
					resource.TestCheckResourceAttrSet(dataSource, "versions.0.op_user"),
					resource.TestCheckResourceAttrSet(dataSource, "versions.0.op_time"),
					resource.TestCheckResourceAttrSet(dataSource, "versions.0.maintainers"),
					resource.TestCheckResourceAttrSet(dataSource, "versions.0.usage_count"),
					resource.TestCheckResourceAttrSet(dataSource, "versions.0.runtime_attribution"),
				),
			},
		},
	})
}

const testDataSourceCodeartsPipelinePluginVersions_basic = `
data "huaweicloud_codearts_pipeline_plugin_versions" "test" {
  plugin_name = "official_devcloud_cloudBuild"
}
`
