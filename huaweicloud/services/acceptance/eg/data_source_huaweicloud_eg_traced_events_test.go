package eg

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataTracedEvents_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_eg_traced_events.test"
		dc  = acceptance.InitDataSourceCheck(all)

		byEventId   = "data.huaweicloud_eg_traced_events.filter_by_event_id"
		dcByEventId = acceptance.InitDataSourceCheck(byEventId)

		bySourceName   = "data.huaweicloud_eg_traced_events.filter_by_source_name"
		dcBySourceName = acceptance.InitDataSourceCheck(bySourceName)

		byEventType   = "data.huaweicloud_eg_traced_events.filter_by_event_type"
		dcByEventType = acceptance.InitDataSourceCheck(byEventType)

		bySubscriptionName   = "data.huaweicloud_eg_traced_events.filter_by_subscription_name"
		dcBySubscriptionName = acceptance.InitDataSourceCheck(bySubscriptionName)

		byNotFoundEventId   = "data.huaweicloud_eg_traced_events.filter_by_not_found_event_id"
		dcByNotFoundEventId = acceptance.InitDataSourceCheck(byNotFoundEventId)

		endTime   = time.Now()
		startTime = endTime.Add(-24 * time.Hour)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataTracedEvents_basic(endTime, startTime),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "events.#", regexp.MustCompile(`^[0-9]*$`)),
					resource.TestCheckResourceAttrSet(all, "events.0.id"),
					resource.TestCheckResourceAttrSet(all, "events.0.type"),
					resource.TestCheckResourceAttrSet(all, "events.0.source_name"),
					resource.TestCheckResourceAttrSet(all, "events.0.source_provider"),
					resource.TestCheckResourceAttrSet(all, "events.0.subscription_name"),
					resource.TestCheckResourceAttrSet(all, "events.0.received_time"),

					dcByEventId.CheckResourceExists(),
					resource.TestCheckOutput("is_event_id_filter_useful", "true"),

					dcBySourceName.CheckResourceExists(),
					resource.TestCheckOutput("is_source_name_filter_useful", "true"),

					dcByEventType.CheckResourceExists(),
					resource.TestCheckOutput("is_event_type_filter_useful", "true"),

					dcBySubscriptionName.CheckResourceExists(),
					resource.TestCheckOutput("is_subscription_name_filter_useful", "true"),

					dcByNotFoundEventId.CheckResourceExists(),
					resource.TestCheckOutput("is_not_found_event_id_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataTracedEvents_basic(endTime, startTime time.Time) string {
	return fmt.Sprintf(`
# Query all traced events
data "huaweicloud_eg_traced_events" "test" {
  channel_id = "default"
  start_time = "%[1]s"
  end_time   = "%[2]s"
}

# Filter by event ID
data "huaweicloud_eg_traced_events" "filter_by_event_id" {
  channel_id = "default"
  start_time = "%[1]s"
  end_time   = "%[2]s"
  event_id   = "test-event-id"
}

locals {
  event_id_filter_result = [
    for v in data.huaweicloud_eg_traced_events.filter_by_event_id.events[*].id : v == "test-event-id"
  ]
}

output "is_event_id_filter_useful" {
  value = length(local.event_id_filter_result) > 0 && alltrue(local.event_id_filter_result)
}

# Filter by source name
data "huaweicloud_eg_traced_events" "filter_by_source_name" {
  channel_id = "default"
  start_time = "%[1]s"
  end_time   = "%[2]s"
  source_name = "HC.SMN"
}

locals {
  source_name_filter_result = [
    for v in data.huaweicloud_eg_traced_events.filter_by_source_name.events[*].source_name : v == "HC.SMN"
  ]
}

output "is_source_name_filter_useful" {
  value = length(local.source_name_filter_result) > 0 && alltrue(local.source_name_filter_result)
}

# Filter by event type
data "huaweicloud_eg_traced_events" "filter_by_event_type" {
  channel_id = "default"
  start_time = "%[1]s"
  end_time   = "%[2]s"
  event_type = "SMN:CloudTrace:ConsoleAction"
}

locals {
  event_type_filter_result = [
    for v in data.huaweicloud_eg_traced_events.filter_by_event_type.events[*].type : v == "SMN:CloudTrace:ConsoleAction"
  ]
}

output "is_event_type_filter_useful" {
  value = length(local.event_type_filter_result) > 0 && alltrue(local.event_type_filter_result)
}

# Filter by subscription name
data "huaweicloud_eg_traced_events" "filter_by_subscription_name" {
  channel_id        = "default"
  start_time        = "%[1]s"
  end_time          = "%[2]s"
  subscription_name = "test-subscription"
}

locals {
  subscription_name_filter_result = [
    for v in data.huaweicloud_eg_traced_events.filter_by_subscription_name.events[*].subscription_name : v == "test-subscription"
  ]
}

output "is_subscription_name_filter_useful" {
  value = length(local.subscription_name_filter_result) > 0 && alltrue(local.subscription_name_filter_result)
}

# Filter by not found event ID
data "huaweicloud_eg_traced_events" "filter_by_not_found_event_id" {
  channel_id = "default"
  start_time = "%[1]s"
  end_time   = "%[2]s"
  event_id   = "not-found-event-id"
}

output "is_not_found_event_id_filter_useful" {
  value = length(data.huaweicloud_eg_traced_events.filter_by_not_found_event_id.events) < 1
}
`, startTime.Format("2006-01-02T15:04:05Z"), endTime.Format("2006-01-02T15:04:05Z"))
}
