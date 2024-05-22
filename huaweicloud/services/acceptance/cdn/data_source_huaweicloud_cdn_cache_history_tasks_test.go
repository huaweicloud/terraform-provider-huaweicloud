package cdn

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceCacheHistoryTasks_basic(t *testing.T) {
	rName := "data.huaweicloud_cdn_cache_history_tasks.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCDNURL(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceCacheHistoryTasks_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "tasks.0.id"),
					resource.TestCheckResourceAttrSet(rName, "tasks.0.status"),
					resource.TestCheckResourceAttrSet(rName, "tasks.0.processing"),
					resource.TestCheckResourceAttrSet(rName, "tasks.0.succeed"),
					resource.TestCheckResourceAttrSet(rName, "tasks.0.failed"),
					resource.TestCheckResourceAttrSet(rName, "tasks.0.total"),
					resource.TestCheckResourceAttrSet(rName, "tasks.0.task_type"),
					resource.TestCheckResourceAttrSet(rName, "tasks.0.created_at"),
					resource.TestCheckResourceAttrSet(rName, "tasks.0.file_type"),

					resource.TestCheckOutput("enterprise_project_id_filter_is_useful", "true"),

					resource.TestCheckOutput("status_filter_is_useful", "true"),

					resource.TestCheckOutput("date_filter_is_useful", "true"),

					resource.TestCheckOutput("order_filter_is_useful", "true"),

					resource.TestCheckOutput("file_type_filter_is_useful", "true"),

					resource.TestCheckOutput("task_type_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceCacheHistoryTasks_basic() string {
	now := time.Now()
	startDate := now.Add(-1 * time.Hour).UnixMilli()
	endDate := now.Add(time.Hour).UnixMilli()

	return fmt.Sprintf(`
%[1]s

# Basic test
data "huaweicloud_cdn_cache_history_tasks" "test" {
  depends_on = [huaweicloud_cdn_cache_preheat.test]
}

# Test with enterprise project ID
data "huaweicloud_cdn_cache_history_tasks" "enterprise_project_id_filter" {
  enterprise_project_id = "all"

  depends_on = [huaweicloud_cdn_cache_preheat.test]
}

output "enterprise_project_id_filter_is_useful" {
  value = length(data.huaweicloud_cdn_cache_history_tasks.enterprise_project_id_filter.tasks) > 0
}

# Test with status
locals {
  status = data.huaweicloud_cdn_cache_history_tasks.test.tasks[0].status
}

data "huaweicloud_cdn_cache_history_tasks" "status_filter" {
  status = local.status
}

output "status_filter_is_useful" {
  value = length(data.huaweicloud_cdn_cache_history_tasks.status_filter.tasks) > 0 && alltrue(
    [for v in data.huaweicloud_cdn_cache_history_tasks.status_filter.tasks[*].status : v == local.status]
  )  
}

# Test with start date and end date
data "huaweicloud_cdn_cache_history_tasks" "date_filter" {
  start_date = "%[2]d"
  end_date   = "%[3]d"

  depends_on = [huaweicloud_cdn_cache_preheat.test]
}

output "date_filter_is_useful" {
  value = length(data.huaweicloud_cdn_cache_history_tasks.date_filter.tasks) > 0
}

# Test with order field and order type
data "huaweicloud_cdn_cache_history_tasks" "order_filter" {
  order_field = "succeed"
  order_type  = "asc"

  depends_on = [huaweicloud_cdn_cache_preheat.test]
}

output "order_filter_is_useful" {
  value = length(data.huaweicloud_cdn_cache_history_tasks.order_filter.tasks) > 0
}

# Test with file type
locals {
  file_type = data.huaweicloud_cdn_cache_history_tasks.test.tasks[0].file_type
}

data "huaweicloud_cdn_cache_history_tasks" "file_type_filter" {
  file_type = local.file_type
}

output "file_type_filter_is_useful" {
  value = length(data.huaweicloud_cdn_cache_history_tasks.file_type_filter.tasks) > 0 && alltrue(
    [for v in data.huaweicloud_cdn_cache_history_tasks.file_type_filter.tasks[*].file_type : v == local.file_type]
  )  
}

# Test with task type
locals {
  task_type = data.huaweicloud_cdn_cache_history_tasks.test.tasks[0].task_type
}

data "huaweicloud_cdn_cache_history_tasks" "task_type_filter" {
  task_type = local.task_type
}

output "task_type_filter_is_useful" {
  value = length(data.huaweicloud_cdn_cache_history_tasks.task_type_filter.tasks) > 0 && alltrue(
    [for v in data.huaweicloud_cdn_cache_history_tasks.task_type_filter.tasks[*].task_type : v == local.task_type]
  )  
}
`, testCachePreheat_basic(), startDate, endDate)
}
