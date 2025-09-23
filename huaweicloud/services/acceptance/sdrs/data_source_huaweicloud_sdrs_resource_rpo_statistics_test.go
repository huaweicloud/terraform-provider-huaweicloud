package sdrs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSdrsRpoStatistics_basic(t *testing.T) {
	dataSource := "data.huaweicloud_sdrs_resource_rpo_statistics.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSdrsRpoStatistics_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("is_resource_type_filter_useful", "true"),
					resource.TestCheckOutput("is_start_time_filter_useful", "true"),
					resource.TestCheckOutput("is_end_time_filter_useful", "true"),
				),
			},
		},
	})
}

const testDataSourceSdrsRpoStatistics_basic = `
data "huaweicloud_sdrs_resource_rpo_statistics" "test" {}

# Filter by resource type
locals {
  resource_type = data.huaweicloud_sdrs_resource_rpo_statistics.test.resource_rpo_statistics[0].resource_type
}

data "huaweicloud_sdrs_resource_rpo_statistics" "filter_by_resource_type" {
  resource_type = local.resource_type
}

locals {
  resource_type_filter_result = [
    for v in data.huaweicloud_sdrs_resource_rpo_statistics.filter_by_resource_type.resource_rpo_statistics[*].resource_type :
    v == local.resource_type
  ]
}

output "is_resource_type_filter_useful" {
  value = length(local.resource_type_filter_result) > 0 && alltrue(local.resource_type_filter_result)
}

# Filter by start time
locals {
  start_time = data.huaweicloud_sdrs_resource_rpo_statistics.test.resource_rpo_statistics[0].created_at
}

data "huaweicloud_sdrs_resource_rpo_statistics" "filter_by_start_time" {
  start_time = local.start_time
}

locals {
  start_time_timestamp = format("%sZ", replace(local.start_time, " ", "T"))
  start_time_filter_result = [
    for v in data.huaweicloud_sdrs_resource_rpo_statistics.filter_by_start_time.resource_rpo_statistics :
      timecmp(format("%sZ", replace(v.created_at, " ", "T")), local.start_time_timestamp) >= 0
  ]
}

output "is_start_time_filter_useful" {
  value = length(local.start_time_filter_result) > 0 && alltrue(local.start_time_filter_result)
}

# Filter by end time
locals {
  end_time = data.huaweicloud_sdrs_resource_rpo_statistics.test.resource_rpo_statistics[0].created_at
}

data "huaweicloud_sdrs_resource_rpo_statistics" "filter_by_end_time" {
  end_time = local.end_time
}

locals {
  end_time_timestamp = format("%sZ", replace(local.end_time, " ", "T"))
  end_time_filter_result = [
    for v in data.huaweicloud_sdrs_resource_rpo_statistics.filter_by_end_time.resource_rpo_statistics :
      timecmp(format("%sZ", replace(v.created_at, " ", "T")), local.end_time_timestamp) <= 0
  ]
}

output "is_end_time_filter_useful" {
  value = length(local.end_time_filter_result) > 0 && alltrue(local.end_time_filter_result)
}
`
