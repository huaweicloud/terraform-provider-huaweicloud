package live

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceLiveTranscodings_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_live_transcodings.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
		domainName     = fmt.Sprintf("%s.huaweicloud.com", acceptance.RandomAccResourceNameWithDash())

		byAppName   = "data.huaweicloud_live_transcodings.filter_by_app_name"
		dcByAppName = acceptance.InitDataSourceCheck(byAppName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceLiveTranscodings_basic(domainName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "templates.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "templates.0.app_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "templates.0.quality_info.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "templates.0.quality_info.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "templates.0.quality_info.0.quality"),
					resource.TestCheckResourceAttrSet(dataSourceName, "templates.0.quality_info.0.video_encoding"),
					resource.TestCheckResourceAttrSet(dataSourceName, "templates.0.quality_info.0.width"),
					resource.TestCheckResourceAttrSet(dataSourceName, "templates.0.quality_info.0.height"),
					resource.TestCheckResourceAttrSet(dataSourceName, "templates.0.quality_info.0.bitrate"),
					resource.TestCheckResourceAttrSet(dataSourceName, "templates.0.quality_info.0.low_bitrate_hd"),
					resource.TestCheckResourceAttrSet(dataSourceName, "templates.0.quality_info.0.protocol"),
					resource.TestCheckResourceAttrSet(dataSourceName, "templates.0.quality_info.0.gop"),
					resource.TestCheckResourceAttrSet(dataSourceName, "templates.0.quality_info.0.bitrate_adaptive"),
					resource.TestCheckResourceAttrSet(dataSourceName, "templates.0.quality_info.0.i_frame_interval"),
					resource.TestCheckResourceAttrSet(dataSourceName, "templates.0.quality_info.0.i_frame_policy"),

					dcByAppName.CheckResourceExists(),
					resource.TestCheckOutput("app_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceLiveTranscodings_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_live_domain" "test" {
  name = "%s"
  type = "push"
}

resource "huaweicloud_live_transcoding" "test" {
  domain_name    = huaweicloud_live_domain.test.name
  app_name       = "live"
  video_encoding = "H264"
  low_bitrate_hd = true

  templates {
    name    = "temp"
    width   = 640
    height  = 360
    bitrate = 400
  }
}

data "huaweicloud_live_transcodings" "test" {
  domain_name = huaweicloud_live_domain.test.name

  depends_on = [huaweicloud_live_transcoding.test]
}

locals {
  app_name = data.huaweicloud_live_transcodings.test.templates[0].app_name
}

data "huaweicloud_live_transcodings" "filter_by_app_name" {
  domain_name = huaweicloud_live_domain.test.name
  app_name    = local.app_name
}

output "app_name_filter_useful" {
  value = length(data.huaweicloud_live_transcodings.filter_by_app_name.templates) > 0 && alltrue(
    [for v in data.huaweicloud_live_transcodings.filter_by_app_name.templates[*].app_name : v == local.app_name]
  )
}
`, name)
}
