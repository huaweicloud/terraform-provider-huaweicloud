package eg

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataEventSubscriptions_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_eg_event_subscriptions.test"
		dc  = acceptance.InitDataSourceCheck(all)

		byChannelId   = "data.huaweicloud_eg_event_subscriptions.filter_by_channel_id"
		dcByChannelId = acceptance.InitDataSourceCheck(byChannelId)

		byName   = "data.huaweicloud_eg_event_subscriptions.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byFuzzyName   = "data.huaweicloud_eg_event_subscriptions.filter_by_fuzzy_name"
		dcByFuzzyName = acceptance.InitDataSourceCheck(byFuzzyName)

		byConnectionId   = "data.huaweicloud_eg_event_subscriptions.filter_by_connection_id"
		dcByConnectionId = acceptance.InitDataSourceCheck(byConnectionId)

		byNotFoundName   = "data.huaweicloud_eg_event_subscriptions.filter_by_not_found_name"
		dcByNotFoundName = acceptance.InitDataSourceCheck(byNotFoundName)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataEventSubscriptions_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "subscriptions.#", regexp.MustCompile(`^[0-9]+$`)),
					resource.TestCheckOutput("is_query_result_contains_subscriptions", "true"),

					// Test filter by channel_id
					dcByChannelId.CheckResourceExists(),
					resource.TestCheckOutput("is_channel_id_filter_useful", "true"),
					resource.TestCheckResourceAttrPair(byChannelId, "subscriptions.0.channel_id", "huaweicloud_eg_custom_event_channel.test", "id"),

					// Test filter by name
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckResourceAttr(byName, "subscriptions.0.name", fmt.Sprintf("%s-suffix", name)),
					resource.TestCheckResourceAttr(byName, "subscriptions.0.description", "Created by acceptance test"),
					resource.TestCheckResourceAttr(byName, "subscriptions.0.type", "EVENT"),
					resource.TestCheckResourceAttrSet(byName, "subscriptions.0.status"),
					resource.TestCheckResourceAttrSet(byName, "subscriptions.0.channel_name"),
					resource.TestCheckResourceAttrSet(byName, "subscriptions.0.created_time"),
					resource.TestCheckResourceAttrSet(byName, "subscriptions.0.updated_time"),

					// Test sources structure
					resource.TestCheckResourceAttr(byName, "subscriptions.0.sources.#", "1"),
					resource.TestCheckResourceAttr(byName, "subscriptions.0.sources.0.provider_type", "CUSTOM"),
					resource.TestCheckResourceAttrSet(byName, "subscriptions.0.sources.0.name"),
					resource.TestCheckResourceAttrSet(byName, "subscriptions.0.sources.0.filter"),
					resource.TestCheckResourceAttrSet(byName, "subscriptions.0.sources.0.created_time"),
					resource.TestCheckResourceAttrSet(byName, "subscriptions.0.sources.0.updated_time"),

					// Test targets structure
					resource.TestCheckResourceAttr(byName, "subscriptions.0.targets.#", "2"),
					resource.TestCheckResourceAttrSet(byName, "subscriptions.0.targets.0.id"),
					resource.TestCheckResourceAttrSet(byName, "subscriptions.0.targets.0.name"),
					resource.TestCheckResourceAttrSet(byName, "subscriptions.0.targets.0.provider_type"),
					resource.TestCheckResourceAttrSet(byName, "subscriptions.0.targets.0.smn_detail"),
					resource.TestCheckResourceAttrSet(byName, "subscriptions.0.targets.0.transform"),
					resource.TestCheckResourceAttrSet(byName, "subscriptions.0.targets.0.created_time"),
					resource.TestCheckResourceAttrSet(byName, "subscriptions.0.targets.0.updated_time"),
					resource.TestCheckResourceAttrSet(byName, "subscriptions.0.targets.1.id"),
					resource.TestCheckResourceAttrSet(byName, "subscriptions.0.targets.1.name"),
					resource.TestCheckResourceAttrSet(byName, "subscriptions.0.targets.1.provider_type"),
					resource.TestCheckResourceAttrSet(byName, "subscriptions.0.targets.1.connection_id"),
					resource.TestCheckResourceAttrSet(byName, "subscriptions.0.targets.1.detail"),
					resource.TestCheckResourceAttrSet(byName, "subscriptions.0.targets.1.transform"),
					resource.TestCheckResourceAttrSet(byName, "subscriptions.0.targets.1.created_time"),
					resource.TestCheckResourceAttrSet(byName, "subscriptions.0.targets.1.updated_time"),

					// Test filter by not found name
					dcByNotFoundName.CheckResourceExists(),
					resource.TestCheckOutput("is_not_found_name_filter_useful", "true"),

					// Test filter by fuzzy name
					dcByFuzzyName.CheckResourceExists(),
					resource.TestCheckOutput("is_fuzzy_name_filter_useful", "true"),

					// Test filter by connection_id
					dcByConnectionId.CheckResourceExists(),
					resource.TestCheckOutput("is_connection_id_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataEventSubscriptions_base(name string) string {
	randSmnTargetId, _ := uuid.GenerateUUID()
	randHTTPSTargetId, _ := uuid.GenerateUUID()

	return fmt.Sprintf(`
resource "huaweicloud_eg_custom_event_channel" "test" {
  name = "%[1]s"
}

resource "huaweicloud_eg_custom_event_source" "test" {
  channel_id = huaweicloud_eg_custom_event_channel.test.id
  name       = "%[1]s"
}

resource "huaweicloud_smn_topic" "test" {
  name = "%[1]s"
}

data "huaweicloud_eg_connections" "test" {
  name = "default"
}

resource "huaweicloud_eg_event_subscription" "test" {
  channel_id  = huaweicloud_eg_custom_event_channel.test.id
  name        = "%[1]s-suffix"
  description = "Created by acceptance test"

  sources {
    provider_type = "CUSTOM"
    name          = huaweicloud_eg_custom_event_source.test.name
    filter_rule   = jsonencode({
      "source": [{
        "op":"StringIn",
        "values":[huaweicloud_eg_custom_event_source.test.name]
      }]
    })
  }

  targets {
    id            = "%[2]s"
    provider_type = "OFFICIAL"
    name          = "HC.SMN"
    detail_name   = "smn_detail"
    detail        = jsonencode({
      "subject_transform": {
        "type": "CONSTANT",
        "value": "TEST_CONSTANT"
      },
      "urn": huaweicloud_smn_topic.test.topic_urn,
      "agency_name": "EG_TARGET_AGENCY",
    })
    transform = jsonencode({
      type  = "ORIGINAL"
      value = ""
    })
  }
  targets {
    id            = "%[3]s"
    provider_type = "CUSTOM"
    name          = "HTTPS"
    connection_id = try(data.huaweicloud_eg_connections.test.connections[0].id, "NOT_FOUND")
    detail_name   = "detail"
    detail        = jsonencode({
      "url": "https://test.com/example",
    })
    transform = jsonencode({
      type  = "VARIABLE"
      value = "{\"name\":\"$.data.name\"}",
      template = "My name is $${name}"
    })
  }

  lifecycle {
    ignore_changes = [sources, targets]
  }
}
`, name, randSmnTargetId, randHTTPSTargetId)
}

func testAccDataEventSubscriptions_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

# Query all subscriptions
data "huaweicloud_eg_event_subscriptions" "test" {
  depends_on = [huaweicloud_eg_event_subscription.test]
}

output "is_query_result_contains_subscriptions" {
  value = length(data.huaweicloud_eg_event_subscriptions.test.subscriptions) > 0
}

# Filter by channel ID
locals {
  channel_id = huaweicloud_eg_custom_event_channel.test.id
}

data "huaweicloud_eg_event_subscriptions" "filter_by_channel_id" {
  depends_on = [huaweicloud_eg_event_subscription.test]

  channel_id = local.channel_id
}

locals {
  channel_id_filter_result = [
    for v in data.huaweicloud_eg_event_subscriptions.filter_by_channel_id.subscriptions[*].channel_id : v == local.channel_id
  ]
}

output "is_channel_id_filter_useful" {
  value = length(local.channel_id_filter_result) > 0 && alltrue(local.channel_id_filter_result)
}

# Filter by name (exact match)
locals {
  subscription_name = "%[2]s-suffix"
}

data "huaweicloud_eg_event_subscriptions" "filter_by_name" {
  depends_on = [huaweicloud_eg_event_subscription.test]

  name = local.subscription_name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_eg_event_subscriptions.filter_by_name.subscriptions[*].name : v == local.subscription_name
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by fuzzy name (fuzzy match)
locals {
  subscription_name_prefix = "%[2]s"
}

data "huaweicloud_eg_event_subscriptions" "filter_by_fuzzy_name" {
  depends_on = [huaweicloud_eg_event_subscription.test]

  fuzzy_name = local.subscription_name_prefix
}

locals {
  fuzzy_name_filter_result = [
    for v in data.huaweicloud_eg_event_subscriptions.filter_by_fuzzy_name.subscriptions[*].name : strcontains(v, local.subscription_name_prefix)
  ]
}

output "is_fuzzy_name_filter_useful" {
  value = length(local.fuzzy_name_filter_result) > 0 && alltrue(local.fuzzy_name_filter_result)
}

# Filter by not found name
data "huaweicloud_eg_event_subscriptions" "filter_by_not_found_name" {
  name = "not_found"
}

output "is_not_found_name_filter_useful" {
  value = length(data.huaweicloud_eg_event_subscriptions.filter_by_not_found_name.subscriptions) < 1
}

# Filter by connection ID
locals {
  connection_id = try(data.huaweicloud_eg_connections.test.connections[0].id, "NOT_FOUND")
}

data "huaweicloud_eg_event_subscriptions" "filter_by_connection_id" {
  depends_on = [huaweicloud_eg_event_subscription.test]

  connection_id = local.connection_id
}

locals {
  connection_id_filter_result = [
    for v in data.huaweicloud_eg_event_subscriptions.filter_by_connection_id.subscriptions[*].targets : 
      length([for t in v : t.connection_id == local.connection_id]) > 0
  ]
}

output "is_connection_id_filter_useful" {
  value = length(local.connection_id_filter_result) > 0 && alltrue(local.connection_id_filter_result)
}
`, testAccDataEventSubscriptions_base(name), name)
}
