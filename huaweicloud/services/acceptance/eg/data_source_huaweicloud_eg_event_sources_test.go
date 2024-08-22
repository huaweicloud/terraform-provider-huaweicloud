package eg

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataEventSources_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_eg_event_sources.test"
		dc  = acceptance.InitDataSourceCheck(all)

		byProviderType   = "data.huaweicloud_eg_event_sources.filter_by_provider_type"
		dcByProviderType = acceptance.InitDataSourceCheck(byProviderType)

		byChannelId   = "data.huaweicloud_eg_event_sources.filter_by_channel_id"
		dcByChannelId = acceptance.InitDataSourceCheck(byChannelId)

		byName   = "data.huaweicloud_eg_event_sources.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byNotFoundName   = "data.huaweicloud_eg_event_sources.filter_by_not_found_name"
		dcByNotFoundName = acceptance.InitDataSourceCheck(byNotFoundName)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDMSRocketMQInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataEventSources_basic_step1(name),
			},
			{
				// UpgradeMake sure the attribute 'updated_at' has been configured.
				Config: testAccDataEventSources_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "sources.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckOutput("is_query_result_contains_at_least_two_types", "true"),
					dcByProviderType.CheckResourceExists(),
					resource.TestCheckOutput("is_provider_type_filter_useful", "true"),
					resource.TestCheckResourceAttrSet(byProviderType, "sources.0.label"),
					resource.TestMatchResourceAttr(byProviderType, "sources.0.event_types.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					dcByChannelId.CheckResourceExists(),
					resource.TestCheckOutput("is_channel_id_filter_useful", "true"),
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckResourceAttrPair(byName, "sources.0.id", "huaweicloud_eg_custom_event_source.test", "id"),
					resource.TestCheckResourceAttrPair(byName, "sources.0.channel_id", "huaweicloud_eg_custom_event_channel.test", "id"),
					resource.TestCheckResourceAttrPair(byName, "sources.0.channel_name", "huaweicloud_eg_custom_event_channel.test", "name"),
					resource.TestCheckResourceAttr(byName, "sources.0.name", name),
					resource.TestCheckResourceAttr(byName, "sources.0.description", "Updated by terraform script"),
					resource.TestCheckResourceAttr(byName, "sources.0.provider_type", "CUSTOM"),
					resource.TestCheckResourceAttr(byName, "sources.0.type", "ROCKETMQ"),
					resource.TestCheckResourceAttrSet(byName, "sources.0.status"),
					resource.TestCheckResourceAttrSet(byName, "sources.0.created_at"),
					resource.TestCheckResourceAttrSet(byName, "sources.0.updated_at"),
					resource.TestCheckResourceAttrSet(byName, "sources.0.detail"),
					dcByNotFoundName.CheckResourceExists(),
					resource.TestCheckOutput("is_not_found_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataEventSources_basic_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_dms_rocketmq_instances" "test" {
  instance_id = "%[1]s"
}

resource "huaweicloud_dms_rocketmq_consumer_group" "test" {
  instance_id     = "%[1]s"
  name            = "%[2]s"
  enabled         = true
  broadcast       = true
  brokers         = ["broker-0"]
  retry_max_times = 3

  // Instances of version 5.x have no brokers returned, but versions 4.x have.
  lifecycle {
    ignore_changes = [
      brokers
    ]
  }
}

resource "huaweicloud_dms_rocketmq_topic" "test" {
  instance_id = "%[1]s"
  name        = "%[2]s"
  queue_num   = 3
  permission  = "all"

  brokers {
    name = "broker-0"
  }

  lifecycle {
    ignore_changes = [
      brokers
    ]
  }
}

data "huaweicloud_vpc_subnets" "test" {
  vpc_id = try(data.huaweicloud_dms_rocketmq_instances.test.instances[0].vpc_id, "")
}

resource "huaweicloud_eg_endpoint" "test" {
  name      = "%[2]s"
  vpc_id    = try(data.huaweicloud_dms_rocketmq_instances.test.instances[0].vpc_id, "")
  subnet_id = try(data.huaweicloud_vpc_subnets.test.subnets[0].id, "")
}

resource "huaweicloud_eg_custom_event_channel" "test" {
  depends_on = [huaweicloud_eg_endpoint.test]

  name = "%[2]s"
}
`, acceptance.HW_DMS_ROCKETMQ_INSTANCE_ID, name)
}

func testAccDataEventSources_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_eg_custom_event_source" "test" {
  channel_id  = huaweicloud_eg_custom_event_channel.test.id
  name        = "%[2]s"
  type        = "ROCKETMQ"
  description = "Created by terraform script"
  detail      = jsonencode({
    instance_id     = try(data.huaweicloud_dms_rocketmq_instances.test.instances[0].id, "")
    group           = huaweicloud_dms_rocketmq_consumer_group.test.id
    topic           = huaweicloud_dms_rocketmq_topic.test.id
    enable_acl      = false
    name            = try(data.huaweicloud_dms_rocketmq_instances.test.instances[0].name, "")
    namesrv_address = try(data.huaweicloud_dms_rocketmq_instances.test.instances[0].namesrv_address, "")
    ssl_enable      = false
  })
}
`, testAccDataEventSources_basic_base(name), name)
}

func testAccDataEventSources_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_eg_custom_event_source" "test" {
  channel_id  = huaweicloud_eg_custom_event_channel.test.id
  name        = "%[2]s"
  type        = "ROCKETMQ"
  description = "Updated by terraform script"
  detail      = jsonencode({
    instance_id     = try(data.huaweicloud_dms_rocketmq_instances.test.instances[0].id, "")
    group           = huaweicloud_dms_rocketmq_consumer_group.test.id
    topic           = huaweicloud_dms_rocketmq_topic.test.id
    enable_acl      = false
    name            = try(data.huaweicloud_dms_rocketmq_instances.test.instances[0].name, "")
    namesrv_address = try(data.huaweicloud_dms_rocketmq_instances.test.instances[0].namesrv_address, "")
    ssl_enable      = false
  })
}

data "huaweicloud_eg_event_sources" "test" {
  depends_on = [huaweicloud_eg_custom_event_source.test]
}

// At least one of official event source exist, e.g. official HC.CCE names Default.
output "is_query_result_contains_at_least_two_types" {
  value = length(distinct(data.huaweicloud_eg_event_sources.test.sources[*].provider_type)) >= 2
}

# Filter by provider type
data "huaweicloud_eg_event_sources" "filter_by_provider_type" {
  depends_on = [huaweicloud_eg_custom_event_source.test]

  provider_type = "OFFICIAL"
}

locals {
  provider_type_filter_result = [
    for v in data.huaweicloud_eg_event_sources.filter_by_provider_type.sources[*].provider_type : v == "OFFICIAL"
  ]
}

output "is_provider_type_filter_useful" {
  value = length(local.provider_type_filter_result) > 0 && alltrue(local.provider_type_filter_result)
}

# Filter by channel ID
locals {
  channel_id = huaweicloud_eg_custom_event_channel.test.id
}

data "huaweicloud_eg_event_sources" "filter_by_channel_id" {
  depends_on = [huaweicloud_eg_custom_event_source.test]

  channel_id = local.channel_id
}

locals {
  channel_id_filter_result = [
    for v in data.huaweicloud_eg_event_sources.filter_by_channel_id.sources[*].channel_id : v == local.channel_id
  ]
}

output "is_channel_id_filter_useful" {
  value = length(local.channel_id_filter_result) > 0 && alltrue(local.channel_id_filter_result)
}

# Filter by name
data "huaweicloud_eg_event_sources" "filter_by_name" {
  depends_on = [huaweicloud_eg_custom_event_source.test]

  name = "%[2]s"
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_eg_event_sources.filter_by_name.sources[*].name : v == "%[2]s"
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by not found name
data "huaweicloud_eg_event_sources" "filter_by_not_found_name" {
  name = "not_found"
}

output "is_not_found_name_filter_useful" {
  value = length(data.huaweicloud_eg_event_sources.filter_by_not_found_name.sources) < 1
}
`, testAccDataEventSources_basic_base(name), name)
}
