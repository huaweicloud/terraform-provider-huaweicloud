package live

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDisablePushStreams_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_live_disable_push_streams.test"
		domainName     = fmt.Sprintf("%s.huaweicloud.com", acceptance.RandomAccResourceNameWithDash())
		createTime     = time.Now().UTC().Add(24 * time.Hour).Format("2006-01-02T15:04:05Z")
		dc             = acceptance.InitDataSourceCheck(dataSourceName)

		byAppName   = "data.huaweicloud_live_disable_push_streams.filter_by_app_name"
		dcByAppName = acceptance.InitDataSourceCheck(byAppName)

		byStreamName   = "data.huaweicloud_live_disable_push_streams.filter_by_stream_name"
		dcByStreamName = acceptance.InitDataSourceCheck(byStreamName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDisablePushStreams_basic(domainName, createTime),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "blocks.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "blocks.0.app_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "blocks.0.stream_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "blocks.0.resume_time"),

					dcByAppName.CheckResourceExists(),
					resource.TestCheckOutput("app_name_filter_useful", "true"),

					dcByStreamName.CheckResourceExists(),
					resource.TestCheckOutput("stream_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDisablePushStreams_basic(name, nowTime string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_live_disable_push_streams" "test" {
  depends_on = [
    huaweicloud_live_disable_push_stream.test
  ]

  domain_name = huaweicloud_live_domain.test.name
}

locals {
  app_name = data.huaweicloud_live_disable_push_streams.test.blocks[0].app_name
}

data "huaweicloud_live_disable_push_streams" "filter_by_app_name" {
  domain_name = huaweicloud_live_domain.test.name
  app_name    = local.app_name
}

output "app_name_filter_useful" {
  value = length(data.huaweicloud_live_disable_push_streams.filter_by_app_name.blocks) > 0 && alltrue(
    [for v in data.huaweicloud_live_disable_push_streams.filter_by_app_name.blocks[*].app_name : v == local.app_name]
  )
}

locals {
  stream_name = data.huaweicloud_live_disable_push_streams.test.blocks[0].stream_name
}

data "huaweicloud_live_disable_push_streams" "filter_by_stream_name" {
  domain_name = huaweicloud_live_domain.test.name
  stream_name = local.stream_name
}

output "stream_name_filter_useful" {
  value = length(data.huaweicloud_live_disable_push_streams.filter_by_stream_name.blocks) > 0 && alltrue(
    [for v in data.huaweicloud_live_disable_push_streams.filter_by_stream_name.blocks[*].stream_name : v == local.stream_name]
  )
}
`, testAccDisablePushStream_basic(name, nowTime))
}
