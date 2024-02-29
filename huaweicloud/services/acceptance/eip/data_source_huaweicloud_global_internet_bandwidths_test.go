package eip

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccGlobalInternetBandwidths_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_global_internet_bandwidths.all"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGlobalInternetBandwidthsDataSource_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "internet_bandwidths.#"),

					resource.TestCheckOutput("is_id_filter_useful", "true"),
					resource.TestCheckOutput("is_size_filter_useful", "true"),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("is_access_site_filter_useful", "true"),
					resource.TestCheckOutput("is_type_filter_useful", "true"),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
					resource.TestCheckOutput("is_eps_id_filter_useful", "true"),
					resource.TestCheckOutput("is_tags_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccGlobalInternetBandwidthsDataSource_basic() string {
	rNameWithDash := acceptance.RandomAccResourceNameWithDash()

	return fmt.Sprintf(`
%[1]s

data "huaweicloud_global_internet_bandwidths" "all" {
  depends_on = [huaweicloud_global_internet_bandwidth.test]
}

// filter by ID
data "huaweicloud_global_internet_bandwidths" "filter_by_id" {
  bandwidth_id = huaweicloud_global_internet_bandwidth.test.id
}

locals {
  filter_result_by_id = [for v in data.huaweicloud_global_internet_bandwidths.filter_by_id.internet_bandwidths[*].id : 
    v == huaweicloud_global_internet_bandwidth.test.id]
}

output "is_id_filter_useful" {
  value = length(local.filter_result_by_id) == 1 && alltrue(local.filter_result_by_id) 
}

// filter by name
data "huaweicloud_global_internet_bandwidths" "filter_by_name" {
  name = huaweicloud_global_internet_bandwidth.test.name
}

locals {
  filter_result_by_name = [for v in data.huaweicloud_global_internet_bandwidths.filter_by_name.internet_bandwidths[*].name : 
    v == "%[2]s"]
}

output "is_name_filter_useful" {
  value = length(local.filter_result_by_name) > 0 && alltrue(local.filter_result_by_name) 
}

// filter by size
data "huaweicloud_global_internet_bandwidths" "filter_by_size" {
  size = huaweicloud_global_internet_bandwidth.test.size
}

locals {
  filter_result_by_size = [for v in data.huaweicloud_global_internet_bandwidths.filter_by_size.internet_bandwidths[*].size : 
    v == huaweicloud_global_internet_bandwidth.test.size]
}

output "is_size_filter_useful" {
  value = length(local.filter_result_by_size) > 0 && alltrue(local.filter_result_by_size) 
}

// filter by access site
data "huaweicloud_global_internet_bandwidths" "filter_by_access_site" {
  access_site = huaweicloud_global_internet_bandwidth.test.access_site
}

locals {
  filter_result_by_access_site = [for v in data.huaweicloud_global_internet_bandwidths.filter_by_access_site.internet_bandwidths[*].access_site : 
    v == huaweicloud_global_internet_bandwidth.test.access_site]
}

output "is_access_site_filter_useful" {
  value = length(local.filter_result_by_access_site) > 0 && alltrue(local.filter_result_by_access_site) 
}

// filter by type
data "huaweicloud_global_internet_bandwidths" "filter_by_type" {
  type = huaweicloud_global_internet_bandwidth.test.type
}

locals {
  filter_result_by_type = [for v in data.huaweicloud_global_internet_bandwidths.filter_by_type.internet_bandwidths[*].type : 
    v == huaweicloud_global_internet_bandwidth.test.type]
}

output "is_type_filter_useful" {
  value = length(local.filter_result_by_type) > 0 && alltrue(local.filter_result_by_type) 
}

// filter by status
data "huaweicloud_global_internet_bandwidths" "filter_by_status" {
  status = huaweicloud_global_internet_bandwidth.test.status
}

locals {
  filter_result_by_status = [for v in data.huaweicloud_global_internet_bandwidths.filter_by_status.internet_bandwidths[*].status : 
    v == huaweicloud_global_internet_bandwidth.test.status]
}

output "is_status_filter_useful" {
  value = length(local.filter_result_by_status) > 0 && alltrue(local.filter_result_by_status) 
}

// filter by eps ID
data "huaweicloud_global_internet_bandwidths" "filter_by_eps_id" {
  enterprise_project_id = huaweicloud_global_internet_bandwidth.test.enterprise_project_id
}

locals {
  filter_result_by_eps_id = [for v in data.huaweicloud_global_internet_bandwidths.filter_by_eps_id.internet_bandwidths[*].enterprise_project_id : 
    v == huaweicloud_global_internet_bandwidth.test.enterprise_project_id]
}

output "is_eps_id_filter_useful" {
  value = length(local.filter_result_by_eps_id) > 0 && alltrue(local.filter_result_by_eps_id) 
}

// filter by tags
data "huaweicloud_global_internet_bandwidths" "filter_by_tags" {
  tags = huaweicloud_global_internet_bandwidth.test.tags
}

locals {
  filter_result_by_tags = [for v in data.huaweicloud_global_internet_bandwidths.filter_by_tags.internet_bandwidths[*].tags : 
    v == huaweicloud_global_internet_bandwidth.test.tags]
}

output "is_tags_filter_useful" {
  value = length(local.filter_result_by_tags) > 0 && alltrue(local.filter_result_by_tags) 
}
`, testAccInternetBandwidth_basic(rNameWithDash), rNameWithDash)
}
