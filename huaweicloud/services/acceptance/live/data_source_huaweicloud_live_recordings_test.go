package live

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceLiveRecordings_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_live_recordings.test"
		randInt        = acctest.RandInt()
		domainName     = fmt.Sprintf("%s.huaweicloud.com", acceptance.RandomAccResourceNameWithDash())
		dc             = acceptance.InitDataSourceCheck(dataSourceName)

		byDomainName   = "data.huaweicloud_live_recordings.filter_by_domain_name"
		dcByDomainName = acceptance.InitDataSourceCheck(byDomainName)

		byAppName   = "data.huaweicloud_live_recordings.filter_by_app_name"
		dcByAppName = acceptance.InitDataSourceCheck(byAppName)

		byStreamName   = "data.huaweicloud_live_recordings.filter_by_stream_name"
		dcByStreamName = acceptance.InitDataSourceCheck(byStreamName)

		byType   = "data.huaweicloud_live_recordings.filter_by_type"
		dcByType = acceptance.InitDataSourceCheck(byType)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceLiveRecordings_basic(domainName, randInt),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.domain_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.app_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.stream_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.default_record_config.0.obs.0.bucket"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.default_record_config.0.obs.0.region"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.default_record_config.0.hls.0.recording_length"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.default_record_config.0.hls.0.file_naming"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.default_record_config.0.hls.0.record_slice_duration"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.updated_at"),

					dcByDomainName.CheckResourceExists(),
					resource.TestCheckOutput("domain_name_filter_useful", "true"),

					dcByAppName.CheckResourceExists(),
					resource.TestCheckOutput("app_name_filter_useful", "true"),

					dcByStreamName.CheckResourceExists(),
					resource.TestCheckOutput("stream_name_filter_useful", "true"),

					dcByType.CheckResourceExists(),
					resource.TestCheckOutput("type_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceLiveRecordings_base(name string, randInt int) string {
	return fmt.Sprintf(`
resource "huaweicloud_live_domain" "test" {
  name = "%[1]s"
  type = "push"
}

resource "huaweicloud_obs_bucket" "test" {
  bucket        = "tf-bucket-%[2]d"
  storage_class = "STANDARD"
  acl           = "private"
}

resource "huaweicloud_live_bucket_authorization" "test" {
  bucket = huaweicloud_obs_bucket.test.bucket
}

resource "huaweicloud_live_recording" "test" {
  domain_name = huaweicloud_live_domain.test.name
  app_name    = "live"
  stream_name = "stream"
  type        = "CONTINUOUS_RECORD"

  obs {
    region = huaweicloud_obs_bucket.test.region
    bucket = huaweicloud_obs_bucket.test.bucket
  }

  hls {
    recording_length = 120
  }

  depends_on = [huaweicloud_live_bucket_authorization.test]
}
`, name, randInt)
}

func testDataSourceLiveRecordings_basic(name string, randInt int) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_live_recordings" "test" {
  depends_on = [
    huaweicloud_live_recording.test
  ]
}

locals {
  domain_name = data.huaweicloud_live_recordings.test.rules[0].domain_name
}

data "huaweicloud_live_recordings" "filter_by_domain_name" {
  domain_name = local.domain_name
}

output "domain_name_filter_useful" {
  value = length(data.huaweicloud_live_recordings.filter_by_domain_name.rules) > 0 && alltrue(
    [for v in data.huaweicloud_live_recordings.filter_by_domain_name.rules[*].domain_name : v == local.domain_name]
  )
}

locals {
  app_name = data.huaweicloud_live_recordings.test.rules[0].app_name
}

data "huaweicloud_live_recordings" "filter_by_app_name" {
  app_name = local.app_name
}

output "app_name_filter_useful" {
  value = length(data.huaweicloud_live_recordings.filter_by_app_name.rules) > 0 && alltrue(
    [for v in data.huaweicloud_live_recordings.filter_by_app_name.rules[*].app_name : v == local.app_name]
  )
}

locals {
  stream_name = data.huaweicloud_live_recordings.test.rules[0].stream_name
}

data "huaweicloud_live_recordings" "filter_by_stream_name" {
  stream_name = local.stream_name
}

output "stream_name_filter_useful" {
  value = length(data.huaweicloud_live_recordings.filter_by_stream_name.rules) > 0 && alltrue(
    [for v in data.huaweicloud_live_recordings.filter_by_stream_name.rules[*].stream_name : v == local.stream_name]
  )
}

locals {
  type = data.huaweicloud_live_recordings.test.rules[0].type
}

data "huaweicloud_live_recordings" "filter_by_type" {
  type = local.type
}

output "type_filter_useful" {
  value = length(data.huaweicloud_live_recordings.filter_by_type.rules) > 0 && alltrue(
    [for v in data.huaweicloud_live_recordings.filter_by_type.rules[*].type : v == local.type]
  )
}
`, testDataSourceLiveRecordings_base(name, randInt))
}
