package lts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceStreams_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_lts_streams.test"
		rName      = acceptance.RandomAccResourceName()
		dc         = acceptance.InitDataSourceCheck(dataSource)

		byName   = "data.huaweicloud_lts_streams.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byLogGroupName   = "data.huaweicloud_lts_streams.filter_by_log_group_name"
		dcByLogGroupName = acceptance.InitDataSourceCheck(byLogGroupName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceStreams_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "streams.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestMatchResourceAttr(dataSource, "streams.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					dcByName.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(byName, "streams.0.id", "huaweicloud_lts_stream.test", "id"),
					resource.TestCheckResourceAttr(byName, "streams.0.name", rName),
					resource.TestCheckResourceAttr(byName, "streams.0.ttl_in_days", "60"),
					resource.TestCheckResourceAttr(byName, "streams.0.tags._sys_enterprise_project_id", "0"),
					resource.TestCheckResourceAttr(byName, "streams.0.tags.terraform", "test"),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					dcByLogGroupName.CheckResourceExists(),
					resource.TestCheckOutput("is_log_group_name_filter_useful", "true"),
				),
			},
			{
				Config: testDataSourceStreams_logStreamNotFoundError(),
				// The log stream does not exist.
				ExpectError: regexp.MustCompile("The log stream does not existed"),
			},
			{
				Config: testDataSourceStreams_logGroupNotFoundError(),
				// The log group does not exist.
				ExpectError: regexp.MustCompile("The log group does not existed"),
			},
		},
	})
}

func testDataSourceStreams_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_lts_group" "test" {
  group_name  = "%[1]s"
  ttl_in_days = 30
}

resource "huaweicloud_lts_stream" "test" {
  group_id              = huaweicloud_lts_group.test.id
  stream_name           = huaweicloud_lts_group.test.group_name
  ttl_in_days           = 60
  enterprise_project_id = "0"

  tags = {
    foo       = "bar"
    terraform = "test"
  }
}

data "huaweicloud_lts_streams" "test" {
  depends_on = [
    huaweicloud_lts_stream.test
  ]
}

locals {
  stream_name = huaweicloud_lts_stream.test.stream_name
}

# The name is an exact match.
data "huaweicloud_lts_streams" "filter_by_name" {
  depends_on = [
    huaweicloud_lts_stream.test
  ]

  name = local.stream_name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_lts_streams.filter_by_name.streams[*].name : v == local.stream_name
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

locals {
  log_group_name = huaweicloud_lts_group.test.group_name
}

data "huaweicloud_lts_streams" "filter_by_log_group_name" {
  depends_on = [
    huaweicloud_lts_stream.test
  ]

  log_group_name = local.log_group_name
}

# The name of the log group is not returned. Since there is only one log stream under this log group,
# the ID of the log stream is used for assertion.
locals {
  stream_ids = data.huaweicloud_lts_streams.filter_by_log_group_name.streams[*].id
}

output "is_log_group_name_filter_useful" {
  value = length(local.stream_ids) == 1 && local.stream_ids[0] == huaweicloud_lts_stream.test.id
}
`, name)
}

func testDataSourceStreams_logStreamNotFoundError() string {
	return `
data "huaweicloud_lts_streams" "not_found_log_stream" {
  name = "not_found_log_stream_name"
}
`
}

func testDataSourceStreams_logGroupNotFoundError() string {
	return `
data "huaweicloud_lts_streams" "not_found_log_group" {
  log_group_name = "not_found_log_group_name"
}
`
}
