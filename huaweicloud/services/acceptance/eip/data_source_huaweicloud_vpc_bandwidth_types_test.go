package eip

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceBandwidthTypes_basic(t *testing.T) {
	dataSource := "data.huaweicloud_vpc_bandwidth_types.all"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceBandwidthTypes_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "share_bandwidth_types.#"),
					resource.TestCheckResourceAttrSet(dataSource, "share_bandwidth_types.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "share_bandwidth_types.0.name_en"),
					resource.TestCheckResourceAttrSet(dataSource, "share_bandwidth_types.0.name_zh"),
					resource.TestCheckResourceAttrSet(dataSource, "share_bandwidth_types.0.bandwidth_type"),
					resource.TestCheckResourceAttrSet(dataSource, "share_bandwidth_types.0.public_border_group"),
					resource.TestCheckResourceAttrSet(dataSource, "share_bandwidth_types.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "share_bandwidth_types.0.updated_at"),

					resource.TestCheckOutput("bandwidth_type_filter_validation", "true"),
					resource.TestCheckOutput("name_en_filter_validation", "true"),
					resource.TestCheckOutput("name_zh_filter_validation", "true"),
					resource.TestCheckOutput("public_border_group_filter_validation", "true"),
				),
			},
		},
	})
}

const testDataSourceBandwidthTypes_basic = `
data "huaweicloud_vpc_bandwidth_types" "all" {}

data "huaweicloud_vpc_bandwidth_types" "test" {
  bandwidth_type      = local.test_refer.bandwidth_type
  name_en             = local.test_refer.name_en
  name_zh             = local.test_refer.name_zh
  public_border_group = local.test_refer.public_border_group
}

locals {
  test_refer   = data.huaweicloud_vpc_bandwidth_types.all.share_bandwidth_types[0]
  test_results = data.huaweicloud_vpc_bandwidth_types.test
}

output "bandwidth_type_filter_validation" {
  value = length(local.test_results.share_bandwidth_types) > 0 && alltrue([
    for v in local.test_results.share_bandwidth_types[*].bandwidth_type : v == local.test_refer.bandwidth_type])
}

output "name_en_filter_validation" {
  value = length(local.test_results.share_bandwidth_types) > 0 && alltrue([
    for v in local.test_results.share_bandwidth_types[*].name_en : v == local.test_refer.name_en])
}

output "name_zh_filter_validation" {
  value = length(local.test_results.share_bandwidth_types) > 0 && alltrue([
    for v in local.test_results.share_bandwidth_types[*].name_zh : v == local.test_refer.name_zh])
}

output "public_border_group_filter_validation" {
  value = length(local.test_results.share_bandwidth_types) > 0 && alltrue([
    for v in local.test_results.share_bandwidth_types[*].public_border_group : v == local.test_refer.public_border_group])
}
`
