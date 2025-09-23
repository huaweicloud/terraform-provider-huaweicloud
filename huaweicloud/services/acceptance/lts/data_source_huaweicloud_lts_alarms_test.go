package lts

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAlarms_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		dataSource = "data.huaweicloud_lts_alarms.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)

		byLevel   = "data.huaweicloud_lts_alarms.filter_by_level"
		dcByLevel = acceptance.InitDataSourceCheck(byLevel)

		bySortDesc   = "data.huaweicloud_lts_alarms.filter_by_sort_desc"
		dcBySortDesc = acceptance.InitDataSourceCheck(bySortDesc)

		bySortAsc   = "data.huaweicloud_lts_alarms.filter_by_sort_asc"
		dcBySortAsc = acceptance.InitDataSourceCheck(bySortAsc)

		byTimeRange   = "data.huaweicloud_lts_alarms.filter_by_time_range"
		dcByTimeRange = acceptance.InitDataSourceCheck(byTimeRange)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckLtsAlarmActionRuleName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"null": {
				Source:            "hashicorp/null",
				VersionConstraint: "3.2.1",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAlarms_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "alarms.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dataSource, "alarms.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "alarms.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "alarms.0.arrives_at"),
					resource.TestCheckResourceAttrSet(dataSource, "alarms.0.starts_at"),
					resource.TestCheckResourceAttrSet(dataSource, "alarms.0.ends_at"),
					dcByLevel.CheckResourceExists(),
					resource.TestCheckOutput("is_alarm_level_filter_useful", "true"),
					dcBySortDesc.CheckResourceExists(),
					resource.TestCheckOutput("is_sort_filter_is_useful", "true"),
					dcBySortAsc.CheckResourceExists(),
					dcByTimeRange.CheckResourceExists(),
					resource.TestCheckOutput("is_arrives_at_set_and_valid", "true"),
					resource.TestCheckOutput("is_starts_at_set_and_valid", "true"),
					resource.TestCheckOutput("is_timeout_set_and_valid", "true"),
					resource.TestCheckOutput("is_annotations_set_and_valid", "true"),
					resource.TestCheckOutput("is_metadata_set_and_valid", "true"),
				),
			},
		},
	})
}

func testDataSourceAlarms_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_lts_group" "test" {
  group_name  = "%[1]s"
  ttl_in_days = 30
}

resource "huaweicloud_lts_stream" "test" {
  group_id    = huaweicloud_lts_group.test.id
  stream_name = "%[1]s"
}

resource "huaweicloud_lts_keywords_alarm_rule" "test" {
  name                   = "%[1]s"
  alarm_level            = "INFO"
  description            = "created by terraform"
  alarm_action_rule_name = "%[2]s"
  send_notifications     = true
  notification_frequency = 0

  keywords_requests {
    keywords               = "key_words"
    condition              = "<"
    number                 = 1
    log_group_id           = huaweicloud_lts_group.test.id
    log_stream_id          = huaweicloud_lts_stream.test.id
    search_time_range_unit = "minute"
    search_time_range      = 5
  }

  frequency {
    type            = "FIXED_RATE"
    fixed_rate      = 1
    fixed_rate_unit = "minute"
  }
}

# Wait for the alarm to be generated.
resource "null_resource" "test" {
  depends_on = [huaweicloud_lts_keywords_alarm_rule.test]

  provisioner "local-exec" {
    command = "sleep 90;"
  }
}
`, name, acceptance.HW_LTS_ALARM_ACTION_RULE_NAME)
}

func testDataSourceAlarms_basic(name string) string {
	currentTime := time.Now()
	startTime := currentTime.Add(-1 * time.Hour).UnixMilli()
	endTime := currentTime.UnixMilli()
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_lts_alarms" "test" {
  type                 = "active_alert"
  whether_custom_field = true
  start_time           = "%[2]v"
  end_time             = "%[3]v"

  depends_on = [null_resource.test]
}

# Filter by alarm level.
data "huaweicloud_lts_alarms" "filter_by_level" {
  type                 = "active_alert"
  whether_custom_field = false
  time_range           = 30
  alarm_level_ids      = ["Info"]

  depends_on = [null_resource.test]
}

locals {
  alarm_level_filter_result = [
    for v in flatten(data.huaweicloud_lts_alarms.filter_by_level.alarms[*].metadata[*].event_severity) : v == "Info"
  ]
}

output "is_alarm_level_filter_useful" {
  value = length(local.alarm_level_filter_result) > 0 && alltrue(local.alarm_level_filter_result)
}

# Filter in descending order using 'arrives_at' field as key.
data "huaweicloud_lts_alarms" "filter_by_sort_desc" {
  type                 = "active_alert"
  whether_custom_field = false
  time_range           = 30

  sort {
    order_by = ["arrives_at"]
    order    = "desc"
  }

  depends_on = [null_resource.test]
}

# Filter in ascending order using 'arrives_at' field as key.
data "huaweicloud_lts_alarms" "filter_by_sort_asc" {
  type                 = "active_alert"
  whether_custom_field = false
  time_range           = 30

  sort {
    order_by = ["arrives_at"]
    order    = "asc"
  }

  depends_on = [null_resource.test]
}

locals {
  sort_desc_filter_result   = data.huaweicloud_lts_alarms.filter_by_sort_desc.alarms
  sort_desc_last_arrives_at = data.huaweicloud_lts_alarms.filter_by_sort_desc.alarms[length(local.sort_desc_filter_result) - 1].arrives_at
  sort_asc_first_arrives_at = data.huaweicloud_lts_alarms.filter_by_sort_asc.alarms[0].arrives_at
}

output "is_sort_filter_is_useful" {
  value = length(local.sort_desc_filter_result) >=2 && local.sort_desc_last_arrives_at == local.sort_asc_first_arrives_at
}

# Filter by time range.
data "huaweicloud_lts_alarms" "filter_by_time_range" {
  type                 = "active_alert"
  whether_custom_field = false
  time_range           = 30

  depends_on = [null_resource.test]
}

locals {
  active_alarm = try([for v in data.huaweicloud_lts_alarms.filter_by_time_range.alarms :
  v if v.metadata[0].event_id == huaweicloud_lts_keywords_alarm_rule.test.id][0], {})
}

output "is_arrives_at_set_and_valid" {
  value = try(local.active_alarm.arrives_at > 0, false)
}

output "is_starts_at_set_and_valid" {
  value = try(local.active_alarm.starts_at > 0, false)
}

output "is_timeout_set_and_valid" {
  value = try(local.active_alarm.timeout > 0, false)
}

output "is_annotations_set_and_valid" {
  value = try(alltrue([
    length(local.active_alarm.annotations) > 0 &&
    local.active_alarm.annotations[0].alarm_action_rule_name != "" &&
    local.active_alarm.annotations[0].alarm_rule_alias != "" &&
    local.active_alarm.annotations[0].alarm_rule_url != "" &&
    local.active_alarm.annotations[0].alarm_status != "" &&
    local.active_alarm.annotations[0].condition_expression != "" &&
    local.active_alarm.annotations[0].condition_expression_with_value != "" &&
    local.active_alarm.annotations[0].current_value != "" &&
    local.active_alarm.annotations[0].frequency != "" &&
    local.active_alarm.annotations[0].log_info != "" &&
    local.active_alarm.annotations[0].message != "" &&
    local.active_alarm.annotations[0].notification_frequency != "" &&
    local.active_alarm.annotations[0].old_annotations != "" &&
    local.active_alarm.annotations[0].recovery_policy != "" &&
    local.active_alarm.annotations[0].type != ""
  ]), false)
}

output "is_metadata_set_and_valid" {
  value = try(alltrue([
    length(local.active_alarm.metadata) > 0 &&
    local.active_alarm.metadata[0].event_id != "" &&
    local.active_alarm.metadata[0].event_name != "" &&
    local.active_alarm.metadata[0].event_severity != "" &&
    local.active_alarm.metadata[0].event_subtype != "" &&
    local.active_alarm.metadata[0].event_type != "" &&
    local.active_alarm.metadata[0].log_group_name != "" &&
    local.active_alarm.metadata[0].log_stream_name != "" &&
    local.active_alarm.metadata[0].lts_alarm_type != "" &&
    local.active_alarm.metadata[0].resource_id != "" &&
    local.active_alarm.metadata[0].resource_provider != "" &&
    local.active_alarm.metadata[0].resource_type != ""
  ]), false)
}
`, testDataSourceAlarms_base(name), startTime, endTime)
}
