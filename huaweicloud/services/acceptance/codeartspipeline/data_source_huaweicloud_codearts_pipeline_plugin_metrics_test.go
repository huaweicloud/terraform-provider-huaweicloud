package codeartspipeline

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourcePipelinePluginMetrics_basic(t *testing.T) {
	dataSource := "data.huaweicloud_codearts_pipeline_plugin_metrics.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testDataSourcePipelinePluginMetrics_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "metrics.#"),
					resource.TestCheckResourceAttrSet(dataSource, "metrics.0.output_key"),
					resource.TestCheckResourceAttrSet(dataSource, "metrics.0.version"),
					resource.TestCheckResourceAttrSet(dataSource, "metrics.0.output_value"),
				),
			},
		},
	})
}

const testDataSourcePipelinePluginMetrics_basic = `
data "huaweicloud_codearts_pipeline_plugin_metrics" "test" {
  plugin_name = "official_devcloud_cloudBuild"
  version     = "0.0.5"
}
`
