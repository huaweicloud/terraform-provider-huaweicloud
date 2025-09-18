package coc

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCocScheduledTaskHistories_basic(t *testing.T) {
	dataSource := "data.huaweicloud_coc_scheduled_task_histories.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCocScheduledTaskID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCocScheduledTaskHistories_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "scheduled_task_history_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "scheduled_task_history_list.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "scheduled_task_history_list.0.task_type"),
					resource.TestCheckResourceAttrSet(dataSource, "scheduled_task_history_list.0.execution_id"),
					resource.TestCheckResourceAttrSet(dataSource, "scheduled_task_history_list.0.associated_task_name"),
					resource.TestCheckResourceAttrSet(dataSource, "scheduled_task_history_list.0.region"),
					resource.TestCheckResourceAttrSet(dataSource, "scheduled_task_history_list.0.created_by"),
					resource.TestCheckResourceAttrSet(dataSource, "scheduled_task_history_list.0.started_time"),
					resource.TestCheckResourceAttrSet(dataSource, "scheduled_task_history_list.0.status"),
					resource.TestCheckOutput("task_id_filter_is_useful", "true"),
					resource.TestCheckOutput("region_filter_is_useful", "true"),
					resource.TestCheckOutput("status_filter_is_useful", "true"),
					resource.TestCheckOutput("started_time_filter_is_useful", "true"),
					resource.TestCheckOutput("sort_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceCocScheduledTaskHistories_basic() string {
	currentTime := time.Now()
	tenMinutesLater := currentTime.Add(10 * time.Minute).UnixMilli()
	return fmt.Sprintf(`
locals {
  task_id = "%[1]s"
}

data "huaweicloud_coc_scheduled_task_histories" "test" {
  task_id = local.task_id
}

output "task_id_filter_is_useful" {
  value = length(data.huaweicloud_coc_scheduled_task_histories.test.scheduled_task_history_list) > 0
}

locals {
  region = [for v in data.huaweicloud_coc_scheduled_task_histories.test.scheduled_task_history_list[*].region :
    v if v != ""][0]
}

data "huaweicloud_coc_scheduled_task_histories" "region_filter" {
  task_id = local.task_id
  region  = local.region
}

output "region_filter_is_useful" {
  value = length(data.huaweicloud_coc_scheduled_task_histories.region_filter.scheduled_task_history_list) > 0 && alltrue(
    [for v in data.huaweicloud_coc_scheduled_task_histories.region_filter.scheduled_task_history_list[*].region :
      v == local.region]
  )
}

locals {
  status = [for v in data.huaweicloud_coc_scheduled_task_histories.test.scheduled_task_history_list[*].status :
    v if v != ""][0]
}

data "huaweicloud_coc_scheduled_task_histories" "status_filter" {
  task_id = local.task_id
  status  = local.status
}

output "status_filter_is_useful" {
  value = length(data.huaweicloud_coc_scheduled_task_histories.status_filter.scheduled_task_history_list) > 0 && alltrue(
    [for v in data.huaweicloud_coc_scheduled_task_histories.status_filter.scheduled_task_history_list[*].status :
      v == local.status]
  )
}

data "huaweicloud_coc_scheduled_task_histories" "started_time_filter" {
  task_id            = local.task_id
  started_start_time = 1
  started_end_time   = %[2]v
}

output "started_time_filter_is_useful" {
  value = length(data.huaweicloud_coc_scheduled_task_histories.started_time_filter.scheduled_task_history_list) > 0
}

data "huaweicloud_coc_scheduled_task_histories" "sort_filter" {
  task_id  = local.task_id
  sort_key = "started_time"
  sort_dir = "desc"
}

output "sort_filter_is_useful" {
  value = length(data.huaweicloud_coc_scheduled_task_histories.sort_filter.scheduled_task_history_list) > 0
}
`, acceptance.HW_COC_SCHEDULED_TASK_ID, tenMinutesLater)
}
