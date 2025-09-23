package codeartspipeline

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCodeartsPipelineTriggerFailedRecords_basic(t *testing.T) {
	dataSource := "data.huaweicloud_codearts_pipeline_trigger_failed_records.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCodeartsPipelineTriggerFailedRecords_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "records.#"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.trigger_type"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.trigger_time"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.executor_id"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.executor_name"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.reason"),
				),
			},
		},
	})
}

func testDataSourceCodeartsPipelineTriggerFailedRecords_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_codearts_pipeline_trigger_failed_records" "test" {
  depends_on = [huaweicloud_codearts_pipeline_action.run[0], huaweicloud_codearts_pipeline_action.run[1]]

  project_id  = huaweicloud_codearts_project.test.id
  pipeline_id = huaweicloud_codearts_pipeline.test.id
}
`, testDataSourceCodeartsPipelineQueueingRecords_queueingBase(name))
}
