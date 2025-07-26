package codeartspipeline

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourcePipelineRuntimeVariables_basic(t *testing.T) {
	dataSource := "data.huaweicloud_codearts_pipeline_runtime_variables.test"
	name := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testDataSourcePipelineRuntimeVariables_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "variables.#"),
					resource.TestCheckResourceAttrSet(dataSource, "variables.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "variables.0.sequence"),
					resource.TestCheckResourceAttrSet(dataSource, "variables.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "variables.0.value"),
					resource.TestCheckResourceAttrSet(dataSource, "variables.0.is_secret"),
					resource.TestCheckResourceAttrSet(dataSource, "variables.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "variables.0.is_runtime"),
					resource.TestCheckResourceAttrSet(dataSource, "variables.0.is_reset"),
				),
			},
		},
	})
}

func testDataSourcePipelineRuntimeVariables_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_codearts_pipeline_runtime_variables" "test" {
  project_id  = huaweicloud_codearts_project.test.id
  pipeline_id = huaweicloud_codearts_pipeline.test.id
}
`, testPipeline_basic(name))
}
