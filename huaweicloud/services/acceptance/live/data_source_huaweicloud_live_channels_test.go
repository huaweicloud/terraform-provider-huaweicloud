package live

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceChannels_basic(t *testing.T) {
	var (
		rName = acceptance.RandomAccResourceName()

		dataSource = "data.huaweicloud_live_channels.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)

		byDomainName   = "data.huaweicloud_live_channels.filter_by_domain_name"
		dcByDomainName = acceptance.InitDataSourceCheck(byDomainName)

		byAppName   = "data.huaweicloud_live_channels.filter_by_app_name"
		dcByAppName = acceptance.InitDataSourceCheck(byAppName)

		byChannelID   = "data.huaweicloud_live_channels.filter_by_channel_id"
		dcByChannelID = acceptance.InitDataSourceCheck(byChannelID)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckLiveStreamingDomainName(t)
			acceptance.TestAccPreCheckLiveIngestRTMPDomainName(t)
			acceptance.TestAccPreCheckLiveTranscodingTemplateID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceChannels_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "channels.0.app_name"),
					resource.TestCheckResourceAttrSet(dataSource, "channels.0.domain_name"),
					resource.TestCheckResourceAttrSet(dataSource, "channels.0.encoder_settings.#"),
					resource.TestCheckResourceAttrSet(dataSource, "channels.0.endpoints.#"),
					resource.TestCheckResourceAttrSet(dataSource, "channels.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "channels.0.input.#"),
					resource.TestCheckResourceAttrSet(dataSource, "channels.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "channels.0.record_settings.#"),
					resource.TestCheckResourceAttrSet(dataSource, "channels.0.state"),

					dcByDomainName.CheckResourceExists(),
					resource.TestCheckOutput("domain_name_filter_is_useful", "true"),

					dcByAppName.CheckResourceExists(),
					resource.TestCheckOutput("app_name_filter_is_useful", "true"),

					dcByChannelID.CheckResourceExists(),
					resource.TestCheckOutput("channel_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceChannels_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_live_channels" "test" {
  depends_on = [
    huaweicloud_live_channel.test
  ]
}

# Filter by domain_name
locals {
  domain_name = data.huaweicloud_live_channels.test.channels[0].domain_name
}

data "huaweicloud_live_channels" "filter_by_domain_name" {
  domain_name = local.domain_name
}

locals {
  domain_name_filter_result = [
    for v in data.huaweicloud_live_channels.filter_by_domain_name.channels[*].domain_name : v == local.domain_name
  ]
}

output "domain_name_filter_is_useful" {
  value = alltrue(local.domain_name_filter_result) && length(local.domain_name_filter_result) > 0
}

# Filter by app_name
locals {
  app_name = data.huaweicloud_live_channels.test.channels[0].app_name
}

data "huaweicloud_live_channels" "filter_by_app_name" {
  app_name = local.app_name
}

locals {
  app_name_filter_result = [
    for v in data.huaweicloud_live_channels.filter_by_app_name.channels[*].app_name : v == local.app_name
  ]
}

output "app_name_filter_is_useful" {
  value = alltrue(local.app_name_filter_result) && length(local.app_name_filter_result) > 0
}

# Filter by channel_id
locals {
  channel_id = data.huaweicloud_live_channels.test.channels[0].id
}

data "huaweicloud_live_channels" "filter_by_channel_id" {
  channel_id = local.channel_id
}

locals {
  channel_id_filter_result = [
    for v in data.huaweicloud_live_channels.filter_by_channel_id.channels[*].id : v == local.channel_id
  ]
}

output "channel_id_filter_is_useful" {
  value = alltrue(local.channel_id_filter_result) && length(local.channel_id_filter_result) > 0
}
`, testLiveChannel_FLV_PULL(name))
}
