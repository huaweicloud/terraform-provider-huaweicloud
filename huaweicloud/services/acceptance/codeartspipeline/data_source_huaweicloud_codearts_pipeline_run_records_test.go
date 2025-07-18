package codeartspipeline

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourcePipelineRunRecords_basic(t *testing.T) {
	dataSource := "data.huaweicloud_codearts_pipeline_run_records.test"
	name := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testDataSourcePipelineRunRecords_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "records.#"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.pipeline_run_id"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.executor_id"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.executor_name"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.trigger_type"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.run_number"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.start_time"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.end_time"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.detail_url"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.modify_url"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.stage_status_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.build_params.#"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.artifact_params.#"),

					resource.TestCheckOutput("is_status_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourcePipelineRunRecords_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_codearts_pipeline_run_records" "test" {
  depends_on = [huaweicloud_codearts_pipeline_action.run]

  project_id  = huaweicloud_codearts_project.test.id
  pipeline_id = huaweicloud_codearts_pipeline.test.id
}

locals {
  status = data.huaweicloud_codearts_pipeline_run_records.test.records[0].status
}

// filter by status
data "huaweicloud_codearts_pipeline_run_records" "filter_by_status" {
  depends_on = [huaweicloud_codearts_pipeline_action.run]

  project_id  = huaweicloud_codearts_project.test.id
  pipeline_id = huaweicloud_codearts_pipeline.test.id
  status      = [local.status]
}

locals {
  filter_result_by_status = [for v in data.huaweicloud_codearts_pipeline_run_records.filter_by_status.records[*].status :
    v == local.status]
}

output "is_status_filter_useful" {
  value = length(local.filter_result_by_status) > 0 && alltrue(local.filter_result_by_status)
}
`, testPipelineAction_run(name))
}
