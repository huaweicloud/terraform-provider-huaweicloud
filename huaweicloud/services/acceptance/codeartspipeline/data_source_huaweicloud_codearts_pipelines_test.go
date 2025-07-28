package codeartspipeline

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourcePipelines_basic(t *testing.T) {
	dataSource := "data.huaweicloud_codearts_pipelines.test"
	name := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testDataSourcePipelines_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "pipelines.#"),
					resource.TestCheckResourceAttrSet(dataSource, "pipelines.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "pipelines.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "pipelines.0.is_publish"),
					resource.TestCheckResourceAttrSet(dataSource, "pipelines.0.is_collect"),
					resource.TestCheckResourceAttrSet(dataSource, "pipelines.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "pipelines.0.manifest_version"),
					resource.TestCheckResourceAttrSet(dataSource, "pipelines.0.latest_run.#"),
					resource.TestCheckResourceAttrSet(dataSource, "pipelines.0.latest_run.0.pipeline_run_id"),
					resource.TestCheckResourceAttrSet(dataSource, "pipelines.0.latest_run.0.executor_id"),
					resource.TestCheckResourceAttrSet(dataSource, "pipelines.0.latest_run.0.executor_name"),
					resource.TestCheckResourceAttrSet(dataSource, "pipelines.0.latest_run.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "pipelines.0.latest_run.0.trigger_type"),
					resource.TestCheckResourceAttrSet(dataSource, "pipelines.0.latest_run.0.run_number"),
					resource.TestCheckResourceAttrSet(dataSource, "pipelines.0.latest_run.0.start_time"),
					resource.TestCheckResourceAttrSet(dataSource, "pipelines.0.latest_run.0.end_time"),
					resource.TestCheckResourceAttrSet(dataSource, "pipelines.0.latest_run.0.stage_status_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "pipelines.0.latest_run.0.build_params.#"),
					resource.TestCheckResourceAttrSet(dataSource, "pipelines.0.latest_run.0.artifact_params.#"),
					resource.TestCheckResourceAttrSet(dataSource, "pipelines.0.convert_sign"),

					resource.TestCheckOutput("is_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourcePipelines_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_codearts_pipelines" "test" {
  depends_on = [huaweicloud_codearts_pipeline_action.run]

  project_id = huaweicloud_codearts_project.test.id
}

// filter by name
data "huaweicloud_codearts_pipelines" "filter_by_name" {
  project_id = huaweicloud_codearts_project.test.id
  name       = huaweicloud_codearts_pipeline.test.name
}

locals {
  filter_result_by_name = [for v in data.huaweicloud_codearts_pipelines.filter_by_name.pipelines[*].name :
    v == huaweicloud_codearts_pipeline.test.name]
}

output "is_name_filter_useful" {
  value = length(local.filter_result_by_name) > 0 && alltrue(local.filter_result_by_name)
}
`, testPipelineAction_run(name))
}
