package geminidb

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccGeminiDBScheduledTasks_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_geminidb_scheduled_tasks.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccCheckGeminidbInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGeminiDBScheduledTasks_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "schedules.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "schedules.0.job_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "schedules.0.job_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "schedules.0.job_status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "schedules.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "schedules.0.instance_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "schedules.0.instance_status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "schedules.0.datastore_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "schedules.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "schedules.0.start_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "schedules.0.end_time"),

					resource.TestCheckOutput("schedules_exist", "true"),
					resource.TestCheckOutput("job_status_filter_useful", "true"),
					resource.TestCheckOutput("job_name_filter_useful", "true"),
					resource.TestCheckOutput("time_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccGeminiDBScheduledTasks_basic() string {
	beginTime := time.Now().UTC().Add(-48 * time.Hour)
	beginTimeString := beginTime.Format("2006-01-02T15:04:05+0000")
	endTime := time.Now().UTC()
	endTimeString := endTime.Format("2006-01-02T15:04:05+0000")
	return fmt.Sprintf(`
data "huaweicloud_geminidb_scheduled_tasks" "test" {}

output "schedules_exist" {
  value = length(data.huaweicloud_geminidb_scheduled_tasks.test.schedules) > 0
}

data "huaweicloud_geminidb_scheduled_tasks" "job_status_filter" {
  job_status = "Completed"
}

output "job_status_filter_useful" {
  value = length(data.huaweicloud_geminidb_scheduled_tasks.job_status_filter.schedules) > 0 && alltrue(
    [for v in data.huaweicloud_geminidb_scheduled_tasks.job_status_filter.schedules[*].job_status : v == "Completed"]
  )
}

data "huaweicloud_geminidb_scheduled_tasks" "job_name_filter" {
  job_name = "RESIZE_FLAVOR"
}

output "job_name_filter_useful" {
  value = length(data.huaweicloud_geminidb_scheduled_tasks.job_name_filter.schedules) > 0 && alltrue(
    [for v in data.huaweicloud_geminidb_scheduled_tasks.job_name_filter.schedules[*].job_name : v == "RESIZE_FLAVOR"]
  )
}

data "huaweicloud_geminidb_scheduled_tasks" "instance_id_filter" {
  instance_id = "%[1]s"
}

output "instance_id_filter_useful" {
  value = length(data.huaweicloud_geminidb_scheduled_tasks.instance_id_filter.schedules) > 0 && alltrue(
    [for v in data.huaweicloud_geminidb_scheduled_tasks.instance_id_filter.schedules[*].instance_id : v == "%[1]s"]
  )
}

data "huaweicloud_geminidb_scheduled_tasks" "time_filter" {
  start_time = "%[2]s"
  end_time   = "%[3]s"
}

output "time_filter_is_useful" {
  value = length(data.huaweicloud_geminidb_scheduled_tasks.time_filter.schedules) > 0
}
`, acceptance.HW_GEMINIDB_INSATNCE_ID, beginTimeString, endTimeString)
}
