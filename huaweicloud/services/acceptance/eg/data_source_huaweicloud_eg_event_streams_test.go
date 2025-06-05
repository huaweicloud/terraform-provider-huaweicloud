package eg

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceEventStreams_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_eg_event_streams.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceEventStreamsNormalCase,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("is_stream_set_and_valid", "true"),
					resource.TestCheckOutput("is_name_set_and_valid", "true"),
					resource.TestCheckOutput("is_option_set_and_valid", "true"),
					resource.TestCheckOutput("is_rule_config_set_and_valid", "true"),
					resource.TestCheckOutput("is_source_set_and_valid", "true"),
					resource.TestCheckOutput("is_sink_set_and_valid", "true"),
				),
			},
		},
	})
}

const testAccDataSourceEventStreamsNormalCase string = `
data "huaweicloud_eg_event_streams" "test" {}

locals {
  stream_filter_result = data.huaweicloud_eg_event_streams.test.event_streams[0]
}

output "is_stream_set_and_valid" {
  value = local.stream_filter_result.id != null
}

output "is_name_set_and_valid" {
  value = local.stream_filter_result.name != null
}

output "is_option_set_and_valid" {
  value = length(local.stream_filter_result.option) > 0
}

output "is_rule_config_set_and_valid" {
  value = length(local.stream_filter_result.rule_config) > 0
}

output "is_source_set_and_valid" {
  value = try(anytrue([
    length(local.stream_filter_result.source[0].source_community_rocketmq) > 0,
    length(local.stream_filter_result.source[0].source_dms_rocketmq) > 0,
    length(local.stream_filter_result.source[0].source_kafka) > 0,
    length(local.stream_filter_result.source[0].source_mobile_rocketmq) > 0,
  ]), false)
}

output "is_sink_set_and_valid" {
  value = try(anytrue([
    length(local.stream_filter_result.sink[0].sink_fg) > 0,
    length(local.stream_filter_result.sink[0].sink_kafka) > 0,
    length(local.stream_filter_result.sink[0].sink_obs) > 0,	
  ]), false)
}
`
