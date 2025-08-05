package codeartspipeline

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourcePipelinePlugins_basic(t *testing.T) {
	dataSource := "data.huaweicloud_codearts_pipeline_plugins.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testDataSourcePipelinePlugins_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "plugins.#"),
					resource.TestCheckResourceAttrSet(dataSource, "plugins.0.plugin_name"),
					resource.TestCheckResourceAttrSet(dataSource, "plugins.0.version"),
					resource.TestCheckResourceAttrSet(dataSource, "plugins.0.version_description"),
					resource.TestCheckResourceAttrSet(dataSource, "plugins.0.version_attribution"),
					resource.TestCheckResourceAttrSet(dataSource, "plugins.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "plugins.0.unique_id"),
					resource.TestCheckResourceAttrSet(dataSource, "plugins.0.op_user"),
					resource.TestCheckResourceAttrSet(dataSource, "plugins.0.op_time"),
					resource.TestCheckResourceAttrSet(dataSource, "plugins.0.plugin_composition_type"),
					resource.TestCheckResourceAttrSet(dataSource, "plugins.0.plugin_attribution"),
					resource.TestCheckResourceAttrSet(dataSource, "plugins.0.workspace_id"),
					resource.TestCheckResourceAttrSet(dataSource, "plugins.0.business_type"),
					resource.TestCheckResourceAttrSet(dataSource, "plugins.0.business_type_display_name"),
					resource.TestCheckResourceAttrSet(dataSource, "plugins.0.maintainers"),
					resource.TestCheckResourceAttrSet(dataSource, "plugins.0.icon_url"),
					resource.TestCheckResourceAttrSet(dataSource, "plugins.0.refer_count"),
					resource.TestCheckResourceAttrSet(dataSource, "plugins.0.usage_count"),
					resource.TestCheckResourceAttrSet(dataSource, "plugins.0.runtime_attribution"),
					resource.TestCheckResourceAttrSet(dataSource, "plugins.0.active"),
				),
			},
		},
	})
}

const testDataSourcePipelinePlugins_basic = `
data "huaweicloud_codearts_pipeline_plugins" "test" {
  business_type      = ["Build", "Gate", "Deploy", "Test", "Normal"]
  plugin_attribution = "official"
}
`
