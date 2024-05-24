package rms

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRmsResources_basic(t *testing.T) {
	dataSource1 := "data.huaweicloud_rms_resources.basic"
	dataSource2 := "data.huaweicloud_rms_resources.filter_by_tracked"
	dataSource3 := "data.huaweicloud_rms_resources.filter_by_type"
	dc1 := acceptance.InitDataSourceCheck(dataSource1)
	dc2 := acceptance.InitDataSourceCheck(dataSource2)
	dc3 := acceptance.InitDataSourceCheck(dataSource3)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceRmsResources_basic,
				Check: resource.ComposeTestCheckFunc(
					dc1.CheckResourceExists(),
					dc2.CheckResourceExists(),
					dc3.CheckResourceExists(),
					resource.TestCheckOutput("is_results_not_empty", "true"),
					resource.TestCheckOutput("is_tracked_filter_useful", "true"),
					resource.TestCheckOutput("is_type_filter_useful", "true"),
				),
			},
		},
	})
}

var testDataSourceDataSourceRmsResources_basic = `
data "huaweicloud_rms_resources" "basic" {}

data "huaweicloud_rms_resources" "filter_by_tracked" {
  tracked = true
}

data "huaweicloud_rms_resources" "filter_by_type" {
  type = "vpc.vpcs"
}

locals {
  service_result = [for v in data.huaweicloud_rms_resources.filter_by_type.resources[*].service : v == "vpc"]
  type_result    = [for v in data.huaweicloud_rms_resources.filter_by_type.resources[*].type : v == "vpcs"]
}

output "is_results_not_empty" {
  value = length(data.huaweicloud_rms_resources.basic.resources) > 0
}

output "is_tracked_filter_useful" {
  value = length(data.huaweicloud_rms_resources.filter_by_tracked.resources) > 0
}

output "is_type_filter_useful" {
  value = alltrue(local.service_result) && length(local.service_result) > 0 && alltrue(local.type_result) && length(local.type_result) > 0
}
`
