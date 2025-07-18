package workspace

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAppScheduleTaskFutureExecutions_basic(t *testing.T) {
	var (
		byDay   = "data.huaweicloud_workspace_app_schedule_task_future_executions.filter_by_day"
		dcByDay = acceptance.InitDataSourceCheck(byDay)

		byWeek   = "data.huaweicloud_workspace_app_schedule_task_future_executions.filter_by_week"
		dcByWeek = acceptance.InitDataSourceCheck(byWeek)

		byMonth   = "data.huaweicloud_workspace_app_schedule_task_future_executions.filter_by_month"
		dcByMonth = acceptance.InitDataSourceCheck(byMonth)

		byFixedTime   = "data.huaweicloud_workspace_app_schedule_task_future_executions.filter_by_fixed_time"
		dcByFixedTime = acceptance.InitDataSourceCheck(byFixedTime)

		byTimeZone   = "data.huaweicloud_workspace_app_schedule_task_future_executions.filter_by_time_zone"
		dcByTimeZone = acceptance.InitDataSourceCheck(byTimeZone)

		byExpireTime   = "data.huaweicloud_workspace_app_schedule_task_future_executions.filter_by_expire_time"
		dcByExpireTime = acceptance.InitDataSourceCheck(byExpireTime)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAppScheduleTaskFutureExecutions_basic(),
				Check: resource.ComposeTestCheckFunc(
					dcByDay.CheckResourceExists(),
					resource.TestCheckOutput("is_day_filter_useful", "true"),
					resource.TestCheckResourceAttr(byDay, "time_zone", "Asia/Shanghai"),
					resource.TestMatchResourceAttr(byDay, "future_executions.0",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2} GMT\+08:00$`)),
					dcByWeek.CheckResourceExists(),
					resource.TestCheckOutput("is_week_filter_useful", "true"),
					dcByMonth.CheckResourceExists(),
					resource.TestCheckOutput("is_month_filter_useful", "true"),
					dcByFixedTime.CheckResourceExists(),
					resource.TestCheckResourceAttr(byFixedTime, "future_executions.#", "1"),
					dcByTimeZone.CheckResourceExists(),
					resource.TestCheckResourceAttr(byTimeZone, "time_zone", "Coordinated Universal Time"),
					resource.TestMatchResourceAttr(byTimeZone, "future_executions.0",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2} GMT$`)),
					dcByExpireTime.CheckResourceExists(),
					resource.TestCheckResourceAttr(byExpireTime, "future_executions.#", "1"),
				),
			},
		},
	})
}

func testAccDataSourceAppScheduleTaskFutureExecutions_basic() string {
	scheduledTime := time.Now().Add(25 * time.Hour)
	return fmt.Sprintf(`
data "huaweicloud_workspace_app_schedule_task_future_executions" "filter_by_day" {
  scheduled_type = "DAY"
  scheduled_time = "02:00:00"
  day_interval   = "2"
}

# The data returned in the future may be greater than 5, so use ">=".
output "is_day_filter_useful" {
  value = length(data.huaweicloud_workspace_app_schedule_task_future_executions.filter_by_day.future_executions) >= 5
}

data "huaweicloud_workspace_app_schedule_task_future_executions" "filter_by_week" {
  scheduled_type = "WEEK"
  scheduled_time = "02:00:00"
  week_list      = "1,3"
}

output "is_week_filter_useful" {
  value = length(data.huaweicloud_workspace_app_schedule_task_future_executions.filter_by_week.future_executions) >= 5
}

# Filter by current month plus one month.
data "huaweicloud_workspace_app_schedule_task_future_executions" "filter_by_month" {
  scheduled_type = "MONTH"
  scheduled_time = "02:00:00"
  month_list     = "%[1]s"
  date_list      = "L"
}

output "is_month_filter_useful" {
  value = length(data.huaweicloud_workspace_app_schedule_task_future_executions.filter_by_month.future_executions) >= 5
}

# Filter by current time plus one day.
data "huaweicloud_workspace_app_schedule_task_future_executions" "filter_by_fixed_time" {
  scheduled_type = "FIXED_TIME"
  scheduled_time = "02:00:00"
  scheduled_date = "%[2]s"
}

data "huaweicloud_workspace_app_schedule_task_future_executions" "filter_by_time_zone" {
  scheduled_type = "WEEK"
  scheduled_time = "02:00:00"
  week_list      = "1,3"
  time_zone      = "Coordinated Universal Time"
}

# Filter by current time, and the expire time is one day later.
data "huaweicloud_workspace_app_schedule_task_future_executions" "filter_by_expire_time" {
  scheduled_type = "DAY"
  scheduled_time = "02:00:00"
  day_interval   = 1
  expire_time    = "%[3]s"
}
`, fmt.Sprintf("%d", scheduledTime.AddDate(0, 1, 0).Month()),
		scheduledTime.Format("2006-01-02"),
		scheduledTime.Format("2006-01-02T15:04:05+08:00"))
}
