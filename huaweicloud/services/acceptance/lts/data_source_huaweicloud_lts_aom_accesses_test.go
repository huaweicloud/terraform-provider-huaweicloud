package lts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAOMAccesses_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_lts_aom_accesses.test"
		rName      = acceptance.RandomAccResourceName()
		dc         = acceptance.InitDataSourceCheck(dataSource)

		byLogGroupName   = "data.huaweicloud_lts_aom_accesses.filter_by_log_group_name"
		dcByLogGroupName = acceptance.InitDataSourceCheck(byLogGroupName)

		byLogStreamName   = "data.huaweicloud_lts_aom_accesses.filter_by_log_stream_name"
		dcByLogStreamName = acceptance.InitDataSourceCheck(byLogStreamName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckLtsAomAccess(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testDataSourceDataSourceAOMAccesses_logGroupNotExist(),
				ExpectError: regexp.MustCompile("The log group does not existed"),
			},
			{
				Config:      testDataSourceDataSourceAOMAccesses_logStreamNotExist(),
				ExpectError: regexp.MustCompile("The log stream does not existed"),
			},
			{
				Config: testDataSourceDataSourceAOMAccesses_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					dcByLogGroupName.CheckResourceExists(),
					resource.TestCheckOutput("is_log_group_name_filter_useful", "true"),
					resource.TestCheckResourceAttrPair(byLogGroupName, "accesses.0.id", "huaweicloud_lts_aom_access.test", "id"),
					resource.TestCheckResourceAttr(byLogGroupName, "accesses.0.name", rName),
					resource.TestCheckResourceAttr(byLogGroupName, "accesses.0.cluster_id", acceptance.HW_LTS_CLUSTER_ID),
					resource.TestCheckResourceAttr(byLogGroupName, "accesses.0.cluster_name", acceptance.HW_LTS_CLUSTER_NAME),
					resource.TestCheckResourceAttr(byLogGroupName, "accesses.0.namespace", "default"),
					resource.TestCheckResourceAttr(byLogGroupName, "accesses.0.workloads.0", "__ALL_DEPLOYMENTS__"),
					resource.TestCheckResourceAttr(byLogGroupName, "accesses.0.container_name", "test_container"),
					resource.TestCheckResourceAttr(byLogGroupName, "accesses.0.access_rules.0.file_name", "/test/*"),
					resource.TestCheckResourceAttrPair(byLogGroupName, "accesses.0.access_rules.0.log_group_id",
						"huaweicloud_lts_group.test", "id"),
					resource.TestCheckResourceAttrPair(byLogGroupName, "accesses.0.access_rules.0.log_group_name",
						"huaweicloud_lts_group.test", "group_name"),
					resource.TestCheckResourceAttrPair(byLogGroupName, "accesses.0.access_rules.0.log_stream_id",
						"huaweicloud_lts_stream.test", "id"),
					resource.TestCheckResourceAttrPair(byLogGroupName, "accesses.0.access_rules.0.log_stream_name",
						"huaweicloud_lts_stream.test", "stream_name"),
					dcByLogStreamName.CheckResourceExists(),
					resource.TestCheckOutput("is_log_stream_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceAOMAccesses_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_lts_aom_accesses" "test" {
  depends_on = [
    huaweicloud_lts_aom_access.test
  ]
}

# Filter by log group name
locals {
  log_group_name = huaweicloud_lts_group.test.group_name
}

data "huaweicloud_lts_aom_accesses" "filter_by_log_group_name" {
  depends_on = [
    huaweicloud_lts_aom_access.test
  ]

  log_group_name = local.log_group_name
}

locals {
  log_group_name_filter_result = [
    for v in flatten(data.huaweicloud_lts_aom_accesses.filter_by_log_group_name.accesses[*].access_rules[*].log_stream_name) :
    v == local.log_group_name
  ]
}

output "is_log_group_name_filter_useful" {
  value = length(local.log_group_name_filter_result) > 0 && alltrue(local.log_group_name_filter_result)
}

# Filter by log stream name
locals {
  log_stream_name = huaweicloud_lts_stream.test.stream_name
}

data "huaweicloud_lts_aom_accesses" "filter_by_log_stream_name" {
  depends_on = [
    huaweicloud_lts_aom_access.test
  ]

  log_stream_name = local.log_stream_name
}

locals {
  log_stream_name_filter_result = [
    for v in flatten(data.huaweicloud_lts_aom_accesses.filter_by_log_stream_name.accesses[*].access_rules[*].log_stream_name) :
   v == local.log_stream_name
  ]
}

output "is_log_stream_name_filter_useful" {
  value = length(local.log_stream_name_filter_result) > 0 && alltrue(local.log_stream_name_filter_result)
}
`, testAOMAccess_basic_step1(name))
}

func testDataSourceDataSourceAOMAccesses_logGroupNotExist() string {
	return `
data "huaweicloud_lts_aom_accesses" "filter_by_log_stream_name" {
  log_group_name = "not_found_log_group"
}
`
}

func testDataSourceDataSourceAOMAccesses_logStreamNotExist() string {
	return `
data "huaweicloud_lts_aom_accesses" "filter_by_log_stream_name" {
  log_stream_name = "not_found_log_stream"
}
`
}
