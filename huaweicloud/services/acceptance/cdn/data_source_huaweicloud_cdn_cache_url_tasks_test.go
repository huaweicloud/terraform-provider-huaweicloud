package cdn

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCacheUrlTasks_basic(t *testing.T) {
	var (
		rName = "data.huaweicloud_cdn_cache_url_tasks.test"
		dc    = acceptance.InitDataSourceCheck(rName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCDNURL(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCacheUrlTasks_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "tasks.0.id"),
					resource.TestCheckResourceAttrSet(rName, "tasks.0.url"),
					resource.TestCheckResourceAttrSet(rName, "tasks.0.status"),
					resource.TestCheckResourceAttrSet(rName, "tasks.0.task_type"),
					resource.TestCheckResourceAttrSet(rName, "tasks.0.task_id"),
					resource.TestCheckResourceAttrSet(rName, "tasks.0.modify_time"),
					resource.TestCheckResourceAttrSet(rName, "tasks.0.created_at"),
					resource.TestCheckResourceAttrSet(rName, "tasks.0.file_type"),

					resource.TestCheckOutput("time_filter_is_useful", "true"),
					resource.TestCheckOutput("url_filter_is_useful", "true"),
					resource.TestCheckOutput("task_type_filter_is_useful", "true"),
					resource.TestCheckOutput("status_filter_is_useful", "true"),
					resource.TestCheckOutput("file_type_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceCacheUrlTasks_basic() string {
	now := time.Now()
	startTime := now.Add(-1 * time.Hour).UnixMilli()
	endTime := now.Add(time.Hour).UnixMilli()

	return fmt.Sprintf(`
%[1]s

# Basic test
data "huaweicloud_cdn_cache_url_tasks" "test" {
  depends_on = [huaweicloud_cdn_cache_preheat.test]
}

# Test with start time and end time
data "huaweicloud_cdn_cache_url_tasks" "time_filter" {
  start_time = "%[2]d"
  end_time   = "%[3]d"

  depends_on = [huaweicloud_cdn_cache_preheat.test]
}

output "time_filter_is_useful" {
  value = length(data.huaweicloud_cdn_cache_url_tasks.time_filter.tasks) > 0
}

# Test with URL
locals {
  url = data.huaweicloud_cdn_cache_url_tasks.test.tasks[0].url
}

data "huaweicloud_cdn_cache_url_tasks" "url_filter" {
  url = local.url
}

output "url_filter_is_useful" {
  value = length(data.huaweicloud_cdn_cache_url_tasks.url_filter.tasks) > 0 && alltrue(
    [for v in data.huaweicloud_cdn_cache_url_tasks.url_filter.tasks[*].url : v == local.url]
  )
}

# Test with task type
locals {
  task_type = data.huaweicloud_cdn_cache_url_tasks.test.tasks[0].task_type
}

data "huaweicloud_cdn_cache_url_tasks" "task_type_filter" {
  task_type = local.task_type
}

output "task_type_filter_is_useful" {
  value = length(data.huaweicloud_cdn_cache_url_tasks.task_type_filter.tasks) > 0 && alltrue(
    [for v in data.huaweicloud_cdn_cache_url_tasks.task_type_filter.tasks[*].task_type : v == local.task_type]
  )
}

# Test with status
locals {
  status = data.huaweicloud_cdn_cache_url_tasks.test.tasks[0].status
}

data "huaweicloud_cdn_cache_url_tasks" "status_filter" {
  status = local.status
}

output "status_filter_is_useful" {
  value = length(data.huaweicloud_cdn_cache_url_tasks.status_filter.tasks) > 0 && alltrue(
    [for v in data.huaweicloud_cdn_cache_url_tasks.status_filter.tasks[*].status : v == local.status]
  )
}

# Test with file type
locals {
  file_type = data.huaweicloud_cdn_cache_url_tasks.test.tasks[0].file_type
}

data "huaweicloud_cdn_cache_url_tasks" "file_type_filter" {
  file_type = local.file_type
}

output "file_type_filter_is_useful" {
  value = length(data.huaweicloud_cdn_cache_url_tasks.file_type_filter.tasks) > 0 && alltrue(
    [for v in data.huaweicloud_cdn_cache_url_tasks.file_type_filter.tasks[*].file_type : v == local.file_type]
  )
}
`, testCachePreheat_basic(), startTime, endTime)
}
