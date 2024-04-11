package cc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCcInterRegionBandwidths_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cc_inter_region_bandwidths.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCcInterRegionBandwidths_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "inter_region_bandwidths.#"),
					resource.TestCheckResourceAttrSet(dataSource, "inter_region_bandwidths.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "inter_region_bandwidths.0.cloud_connection_id"),
					resource.TestCheckResourceAttrSet(dataSource, "inter_region_bandwidths.0.bandwidth_package_id"),
					resource.TestCheckResourceAttrSet(dataSource, "inter_region_bandwidths.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "inter_region_bandwidths.0.updated_at"),

					resource.TestCheckOutput("id_filter_is_useful", "true"),
					resource.TestCheckOutput("cloud_conn_id_filter_is_useful", "true"),
					resource.TestCheckOutput("bw_package_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceCcInterRegionBandwidths_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_cc_inter_region_bandwidths" "test" {
  depends_on = [
    huaweicloud_cc_inter_region_bandwidth.test1,
    huaweicloud_cc_inter_region_bandwidth.test2,
    huaweicloud_cc_inter_region_bandwidth.test3,
    huaweicloud_cc_inter_region_bandwidth.test4,
  ]
}

locals {
  inter_region_bandwidths   = data.huaweicloud_cc_inter_region_bandwidths.test.inter_region_bandwidths
  inter_region_bandwidth_id = local.inter_region_bandwidths[0].id
  cloud_connection_id       = local.inter_region_bandwidths[0].cloud_connection_id
  bandwidth_package_id      = local.inter_region_bandwidths[0].bandwidth_package_id
}

data "huaweicloud_cc_inter_region_bandwidths" "filter_by_id" {
  inter_region_bandwidth_id = local.inter_region_bandwidth_id
}

data "huaweicloud_cc_inter_region_bandwidths" "filter_by_cloud_conn_id" {
  cloud_connection_id = local.cloud_connection_id
}

data "huaweicloud_cc_inter_region_bandwidths" "filter_by_bw_package_id" {
  bandwidth_package_id = local.bandwidth_package_id
}

locals {
  listById          = data.huaweicloud_cc_inter_region_bandwidths.filter_by_id.inter_region_bandwidths
  listByCloudConnId = data.huaweicloud_cc_inter_region_bandwidths.filter_by_cloud_conn_id.inter_region_bandwidths
  listByBwPackageId = data.huaweicloud_cc_inter_region_bandwidths.filter_by_bw_package_id.inter_region_bandwidths
}

output "id_filter_is_useful" {
  value = length(local.listById) > 0 && alltrue(
    [for v in local.listById[*].id : v == local.inter_region_bandwidth_id]
  )
}

output "cloud_conn_id_filter_is_useful" {
  value = length(local.listByCloudConnId) > 0 && alltrue(
    [for v in local.listByCloudConnId[*].cloud_connection_id : v == local.cloud_connection_id]
  )
}

output "bw_package_id_filter_is_useful" {
  value = length(local.listByBwPackageId) > 0 && alltrue(
    [for v in local.listByBwPackageId[*].bandwidth_package_id : v == local.bandwidth_package_id]
  )
}
`, testCcInterRegionBandwidths_dataBasic(name))
}

func testCcInterRegionBandwidths_dataBasic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cc_connection" "test1" {
  name = "%[1]s_1"
}

resource "huaweicloud_cc_connection" "test2" {
  name = "%[1]s_2"
}

resource "huaweicloud_cc_bandwidth_package" "test1" {
  name           = "%[1]s_3"
  local_area_id  = "Chinese-Mainland"
  remote_area_id = "Chinese-Mainland"
  charge_mode    = "bandwidth"
  billing_mode   = 3
  bandwidth      = 5
  description    = "This is an accaptance test"
  resource_id    = huaweicloud_cc_connection.test1.id
  resource_type  = "cloud_connection"
}

resource "huaweicloud_cc_bandwidth_package" "test2" {
  name           = "%[1]s_4"
  local_area_id  = "Chinese-Mainland"
  remote_area_id = "Chinese-Mainland"
  charge_mode    = "bandwidth"
  billing_mode   = 3
  bandwidth      = 5
  description    = "This is an accaptance test"
  resource_id    = huaweicloud_cc_connection.test2.id
  resource_type  = "cloud_connection"
}

resource "huaweicloud_cc_inter_region_bandwidth" "test1" {
  cloud_connection_id  = huaweicloud_cc_connection.test1.id
  bandwidth_package_id = huaweicloud_cc_bandwidth_package.test1.id
  bandwidth            = 1
  inter_region_ids     = ["cn-north-4", "cn-south-1"]
}

resource "huaweicloud_cc_inter_region_bandwidth" "test2" {
  cloud_connection_id  = huaweicloud_cc_connection.test1.id
  bandwidth_package_id = huaweicloud_cc_bandwidth_package.test1.id
  bandwidth            = 4
  inter_region_ids     = ["cn-north-4", "cn-north-9"]
}

resource "huaweicloud_cc_inter_region_bandwidth" "test3" {
  cloud_connection_id  = huaweicloud_cc_connection.test2.id
  bandwidth_package_id = huaweicloud_cc_bandwidth_package.test2.id
  bandwidth            = 2
  inter_region_ids     = ["cn-north-4", "cn-south-1"]
}

resource "huaweicloud_cc_inter_region_bandwidth" "test4" {
  cloud_connection_id  = huaweicloud_cc_connection.test2.id
  bandwidth_package_id = huaweicloud_cc_bandwidth_package.test2.id
  bandwidth            = 3
  inter_region_ids     = ["cn-north-4", "cn-north-9"]
}
`, name)
}
