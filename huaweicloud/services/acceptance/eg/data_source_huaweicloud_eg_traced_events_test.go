package eg

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceTracedEvents_basic(t *testing.T) {
	var (
		name      = acceptance.RandomAccResourceName()
		eventType = "com.example.object.created.v1"
		startTime = time.Now().UTC()
		endTime   = startTime.Add(24 * time.Hour).UTC()

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
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTracedEvents_basic(name, startTime.Format(time.RFC3339), endTime.Format(time.RFC3339), eventType),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(all, "region"),
					resource.TestCheckResourceAttrSet(all, "channel_id"),
					resource.TestCheckResourceAttrSet(all, "start_time"),
					resource.TestCheckResourceAttrSet(all, "end_time"),
					resource.TestMatchResourceAttr(all, "events.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),

					dcByEventId.CheckResourceExists(),
					resource.TestCheckOutput("is_event_id_filter_useful", "true"),

					dcBySourceName.CheckResourceExists(),
					resource.TestCheckOutput("is_source_name_filter_useful", "true"),

					dcByEventType.CheckResourceExists(),
					resource.TestCheckOutput("is_event_type_filter_useful", "true"),

					dcBySubscriptionName.CheckResourceExists(),
					resource.TestCheckOutput("is_subscription_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceTracedEvents_base(name, eventID, eventType, startTime string) string {
	targetID, _ := uuid.GenerateUUID()

	return fmt.Sprintf(`
resource "huaweicloud_eg_custom_event_channel" "test" {
  name = "%[1]s"
}

resource "huaweicloud_eg_custom_event_source" "test" {
  channel_id = huaweicloud_eg_custom_event_channel.test.id
  name       = "%[1]s"
}

resource "huaweicloud_fgs_function" "test" {
  name        = "%[1]s"
  app         = "default"
  handler     = "test.handler"
  memory_size = 128
  timeout     = 3
  runtime     = "Python2.7"
  code_type   = "inline"
  func_code   = <<EOF
# -*- coding:utf-8 -*-
import json
def handler (event, context):
    return {
        "statusCode": 200,
        "isBase64Encoded": False,
        "body": json.dumps(event),
        "headers": {
            "Content-Type": "application/json"
        }
    }
EOF
}

resource "huaweicloud_eg_event_subscription" "test" {
  depends_on = [
    huaweicloud_eg_custom_event_source.test,
    huaweicloud_eg_custom_event_channel.test,
  ]

  channel_id  = huaweicloud_eg_custom_event_channel.test.id
  name        = "%[1]s"
  description = "Created by acceptance test"

  sources {
    provider_type = "CUSTOM"
    id            = huaweicloud_eg_custom_event_source.test.id
    name          = huaweicloud_eg_custom_event_source.test.name
    detail        = jsonencode({})
    filter_rule   = jsonencode({
      "source" : [{
        "op" : "StringIn",
        "values" : [huaweicloud_eg_custom_event_source.test.name]
      }]
    })
  }

  targets {
    id            = "%[2]s"
    name          = "HC.FunctionGraph"
    provider_type = "OFFICIAL"
    detail_name   = "detail"
    detail        = jsonencode({
      "urn":huaweicloud_fgs_function.test.urn,
      "invoke_type":"ASYNC",
      "agency_name":"EG_TARGET_AGENCY"
    })
    transform     = jsonencode({
      type  = "ORIGINAL"
      value = ""
    })
  }

  lifecycle {
    ignore_changes = [sources, targets]
  }
}

resource "huaweicloud_eg_event_batch_action" "test" {
  channel_id = huaweicloud_eg_custom_event_channel.test.id

  events {
    id                = "%[3]s"
    spec_version      = "1.0"
    source            = huaweicloud_eg_custom_event_source.test.name
    type              = "%[4]s"
    data_content_type = "application/json"
    time              = "%[5]s"
    data = jsonencode({
      "terraform": "test"
    })
  }

  depends_on = [
    huaweicloud_eg_event_subscription.test
  ]
}
`, name, targetID, eventID, eventType, startTime)
}

func testAccDataSourceTracedEvents_basic(name, startTime, endTime, eventType string) string {
	eventID, _ := uuid.GenerateUUID()
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_eg_traced_events" "test" {
  depends_on = [huaweicloud_eg_event_batch_action.test]

  channel_id = huaweicloud_eg_custom_event_channel.test.id
  start_time = "%[2]s"
  end_time   = "%[3]s"
}

# Filter by event ID
data "huaweicloud_eg_traced_events" "filter_by_event_id" {
  depends_on = [huaweicloud_eg_event_batch_action.test]

  channel_id = huaweicloud_eg_custom_event_channel.test.id
  start_time = "%[2]s"
  end_time   = "%[3]s"
  event_id   = "%[5]s"
}

output "is_event_id_filter_useful" {
  value = length(data.huaweicloud_eg_traced_events.filter_by_event_id.events) >= 0
}

# Filter by source name
data "huaweicloud_eg_traced_events" "filter_by_source_name" {
  depends_on = [huaweicloud_eg_event_batch_action.test]

  channel_id  = huaweicloud_eg_custom_event_channel.test.id
  start_time  = "%[2]s"
  end_time    = "%[3]s"
  source_name = "%[4]s"
}

output "is_source_name_filter_useful" {
  value = length(data.huaweicloud_eg_traced_events.filter_by_source_name.events) >= 0
}

# Filter by subscription name
data "huaweicloud_eg_traced_events" "filter_by_subscription_name" {
  depends_on = [huaweicloud_eg_event_batch_action.test]

  channel_id        = huaweicloud_eg_custom_event_channel.test.id
  start_time        = "%[2]s"
  end_time          = "%[3]s"
  subscription_name = "%[4]s"
}

output "is_subscription_name_filter_useful" {
  value = length(data.huaweicloud_eg_traced_events.filter_by_subscription_name.events) >= 0
}

# Filter by event type
data "huaweicloud_eg_traced_events" "filter_by_event_type" {
  depends_on = [huaweicloud_eg_event_batch_action.test]

  channel_id = huaweicloud_eg_custom_event_channel.test.id
  start_time = "%[2]s"
  end_time   = "%[3]s"
  event_type = "%[6]s"
}

output "is_event_type_filter_useful" {
  value = length(data.huaweicloud_eg_traced_events.filter_by_event_type.events) >= 0
}
`, testAccDataSourceTracedEvents_base(name, eventType, eventID, startTime), startTime, endTime, name,
		eventID, eventType)
}
