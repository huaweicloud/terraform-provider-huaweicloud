package codeartspipeline

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCodeartsPipelineQueueingRecords_basic(t *testing.T) {
	dataSource := "data.huaweicloud_codearts_pipeline_queueing_records.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCodeartsPipelineQueueingRecords_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "records.#"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.pipeline_run_id"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.enqueue_time"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.trigger_type"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.creator_id"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.creator_name"),
				),
			},
		},
	})
}

func testDataSourceCodeartsPipelineQueueingRecords_queueingBase(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_codearts_pipeline_action" "run" {
  count = 3

  action      = "run"
  project_id  = huaweicloud_codearts_project.test.id
  pipeline_id = huaweicloud_codearts_pipeline.test.id

  sources {
    type = "code"

    params {
      codehub_id     = huaweicloud_codearts_repository.test.id
      git_type       = "codehub"
      git_url        = huaweicloud_codearts_repository.test.https_url
      default_branch = "master"

      build_params {
        build_type    = "branch"
        event_type    = "Manual"
        target_branch = "master"
      }
    }
  }

  choose_jobs   = ["JOB_ijsHS"]
  choose_stages = ["17501613644926b1155c7-cb8a-4af8-87b7-80c73ddb266a"]
  description   = "demo"
}
`, testPipelineAction_run_basic(name))
}

func testDataSourceCodeartsPipelineQueueingRecords_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_codearts_pipeline_queueing_records" "test" {
  depends_on = [huaweicloud_codearts_pipeline_action.run[0], huaweicloud_codearts_pipeline_action.run[1]]

  project_id  = huaweicloud_codearts_project.test.id
  pipeline_id = huaweicloud_codearts_pipeline.test.id
}
`, testDataSourceCodeartsPipelineQueueingRecords_queueingBase(name))
}
