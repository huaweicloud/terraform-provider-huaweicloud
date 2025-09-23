package codeartspipeline

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourcePipelineRunDetail_basic(t *testing.T) {
	dataSource := "data.huaweicloud_codearts_pipeline_run_detail.test"
	name := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testDataSourcePipelineRunDetail_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "manifest_version"),
					resource.TestCheckResourceAttrSet(dataSource, "name"),
					resource.TestCheckResourceAttrSet(dataSource, "is_publish"),
					resource.TestCheckResourceAttrSet(dataSource, "executor_id"),
					resource.TestCheckResourceAttrSet(dataSource, "executor_name"),
					resource.TestCheckResourceAttrSet(dataSource, "status"),
					resource.TestCheckResourceAttrSet(dataSource, "trigger_type"),
					resource.TestCheckResourceAttrSet(dataSource, "run_number"),
					resource.TestCheckResourceAttrSet(dataSource, "start_time"),
					resource.TestCheckResourceAttrSet(dataSource, "end_time"),
					resource.TestCheckResourceAttrSet(dataSource, "language"),
					resource.TestCheckResourceAttrSet(dataSource, "subject_id"),
					resource.TestCheckResourceAttrSet(dataSource, "current_system_time"),
					resource.TestCheckResourceAttrSet(dataSource, "sources.#"),
					resource.TestCheckResourceAttrSet(dataSource, "sources.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "sources.0.params.#"),
					resource.TestCheckResourceAttrSet(dataSource, "sources.0.params.0.git_type"),
					resource.TestCheckResourceAttrSet(dataSource, "sources.0.params.0.git_url"),
					resource.TestCheckResourceAttrSet(dataSource, "sources.0.params.0.ssh_git_url"),
					resource.TestCheckResourceAttrSet(dataSource, "sources.0.params.0.repo_name"),
					resource.TestCheckResourceAttrSet(dataSource, "sources.0.params.0.codehub_id"),
					resource.TestCheckResourceAttrSet(dataSource, "sources.0.params.0.default_branch"),
					resource.TestCheckResourceAttrSet(dataSource, "sources.0.params.0.build_params.#"),
					resource.TestCheckResourceAttrSet(dataSource, "stages.#"),
					resource.TestCheckResourceAttrSet(dataSource, "stages.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "stages.0.identifier"),
					resource.TestCheckResourceAttrSet(dataSource, "stages.0.jobs.#"),
					resource.TestCheckResourceAttrSet(dataSource, "stages.0.jobs.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "stages.0.jobs.0.identifier"),
					resource.TestCheckResourceAttrSet(dataSource, "stages.0.jobs.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "stages.0.pre.#"),
					resource.TestCheckResourceAttrSet(dataSource, "stages.0.pre.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "stages.0.pre.0.status"),
				),
			},
		},
	})
}

func testDataSourcePipelineRunDetail_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_codearts_pipeline_run_detail" "test" {
  project_id      = huaweicloud_codearts_project.test.id
  pipeline_id     = huaweicloud_codearts_pipeline.test.id
  pipeline_run_id = huaweicloud_codearts_pipeline_action.run.pipeline_run_id
}
`, testPipelineAction_run(name))
}
