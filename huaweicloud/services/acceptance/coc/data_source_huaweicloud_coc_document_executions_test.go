package coc

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCocDocumentExecutions_basic(t *testing.T) {
	dataSource := "data.huaweicloud_coc_document_executions.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCocDocumentExecutions_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.execution_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.document_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.document_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.document_version_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.document_version"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.start_time"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.end_time"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.update_time"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.creator"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.parameters.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.parameters.0.key"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.parameters.0.value"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.type"),
					resource.TestCheckOutput("creator_filter_is_useful", "true"),
					resource.TestCheckOutput("start_time_filter_is_useful", "true"),
					resource.TestCheckOutput("end_time_filter_is_useful", "true"),
					resource.TestCheckOutput("document_name_filter_is_useful", "true"),
					resource.TestCheckOutput("document_id_filter_is_useful", "true"),
					resource.TestCheckOutput("exclude_child_executions_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceCocDocumentExecutions_basic(name string) string {
	currentTime := time.Now()
	tenMinutesLater := currentTime.Add(10 * time.Minute).UnixMilli()
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_coc_document_executions" "test" {}

data "huaweicloud_coc_document_executions" "creator_filter" {
  creator = huaweicloud_coc_document_execute.test.creator
}

output "creator_filter_is_useful" {
  value = length(data.huaweicloud_coc_document_executions.creator_filter.data) > 0 && alltrue(
    [for v in data.huaweicloud_coc_document_executions.creator_filter.data[*].creator :
      v == huaweicloud_coc_document_execute.test.creator]
  )
}

data "huaweicloud_coc_document_executions" "start_time_filter" {
  start_time = huaweicloud_coc_document_execute.test.start_time - 1
}

output "start_time_filter_is_useful" {
  value = length(data.huaweicloud_coc_document_executions.start_time_filter.data) > 0 && alltrue(
    [for v in data.huaweicloud_coc_document_executions.start_time_filter.data[*].start_time :
      v >= huaweicloud_coc_document_execute.test.start_time]
  )
}

data "huaweicloud_coc_document_executions" "end_time_filter" {
  end_time = %[2]v
}

output "end_time_filter_is_useful" {
  value = length(data.huaweicloud_coc_document_executions.end_time_filter.data) > 0 && alltrue(
    [for v in data.huaweicloud_coc_document_executions.end_time_filter.data[*].end_time :
      v <= %[2]v]
  )
}

data "huaweicloud_coc_document_executions" "document_name_filter" {
  document_name = huaweicloud_coc_document_execute.test.document_name
}

output "document_name_filter_is_useful" {
  value = length(data.huaweicloud_coc_document_executions.document_name_filter.data) > 0 && alltrue(
    [for v in data.huaweicloud_coc_document_executions.document_name_filter.data[*].document_name :
      v == huaweicloud_coc_document_execute.test.document_name]
  )
}

data "huaweicloud_coc_document_executions" "document_id_filter" {
  document_id = huaweicloud_coc_document_execute.test.document_id
}

output "document_id_filter_is_useful" {
  value = length(data.huaweicloud_coc_document_executions.document_id_filter.data) > 0 && alltrue(
    [for v in data.huaweicloud_coc_document_executions.document_id_filter.data[*].document_id :
      v == huaweicloud_coc_document_execute.test.document_id]
  )
}

data "huaweicloud_coc_document_executions" "exclude_child_executions_filter" {
  exclude_child_executions = true
}

output "exclude_child_executions_filter_is_useful" {
  value = length(data.huaweicloud_coc_document_executions.exclude_child_executions_filter.data) > 0
}
`, testCocDocumentExecute_basic(name), tenMinutesLater)
}
