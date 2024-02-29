package eip

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccGlobalEIPs_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_global_eips.all"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGlobalEIPsDataSource_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "global_eips.#"),

					resource.TestCheckOutput("is_id_filter_useful", "true"),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("is_internet_bandwidth_id_filter_useful", "true"),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
					resource.TestCheckOutput("is_ip_address_filter_useful", "true"),
					resource.TestCheckOutput("is_eps_id_filter_useful", "true"),
					resource.TestCheckOutput("is_tags_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccGlobalEIPsDataSource_basic() string {
	rNameWithDash := acceptance.RandomAccResourceNameWithDash()

	return fmt.Sprintf(`
%[1]s

data "huaweicloud_global_eips" "all" {
  depends_on = [huaweicloud_global_eip.test]
}

// filter by ID
data "huaweicloud_global_eips" "filter_by_id" {
  geip_id = huaweicloud_global_eip.test.id
}

locals {
  filter_result_by_id = [for v in data.huaweicloud_global_eips.filter_by_id.global_eips[*].id : 
    v == huaweicloud_global_eip.test.id]
}

output "is_id_filter_useful" {
  value = length(local.filter_result_by_id) == 1 && alltrue(local.filter_result_by_id) 
}

// filter by name
data "huaweicloud_global_eips" "filter_by_name" {
  name = huaweicloud_global_eip.test.name
}

locals {
  filter_result_by_name = [for v in data.huaweicloud_global_eips.filter_by_name.global_eips[*].name : 
    v == "%[2]s"]
}

output "is_name_filter_useful" {
  value = length(local.filter_result_by_name) > 0 && alltrue(local.filter_result_by_name) 
}

// filter by internet bandwidth id
data "huaweicloud_global_eips" "filter_by_internet_bandwidth_id" {
  internet_bandwidth_id = huaweicloud_global_eip.test.internet_bandwidth_id
}

output "is_internet_bandwidth_id_filter_useful" {
  value = length(data.huaweicloud_global_eips.filter_by_internet_bandwidth_id.global_eips) > 0 && alltrue(
	[for v in data.huaweicloud_global_eips.filter_by_internet_bandwidth_id.global_eips[*].internet_bandwidth_id : 
      v == huaweicloud_global_eip.test.internet_bandwidth_id]
  )
}

// filter by ip address
data "huaweicloud_global_eips" "filter_by_ip_address" {
  ip_address = huaweicloud_global_eip.test.ip_address
}

locals {
  filter_result_by_ip_address = [for v in data.huaweicloud_global_eips.filter_by_ip_address.global_eips[*].ip_address : 
    v == huaweicloud_global_eip.test.ip_address]
}

output "is_ip_address_filter_useful" {
  value = length(local.filter_result_by_ip_address) > 0 && alltrue(local.filter_result_by_ip_address) 
}

// filter by status
data "huaweicloud_global_eips" "filter_by_status" {
  status = huaweicloud_global_eip.test.status
}

locals {
  filter_result_by_status = [for v in data.huaweicloud_global_eips.filter_by_status.global_eips[*].status : 
    v == huaweicloud_global_eip.test.status]
}

output "is_status_filter_useful" {
  value = length(local.filter_result_by_status) > 0 && alltrue(local.filter_result_by_status) 
}

// filter by eps ID
data "huaweicloud_global_eips" "filter_by_eps_id" {
  enterprise_project_id = huaweicloud_global_eip.test.enterprise_project_id
}

locals {
  filter_result_by_eps_id = [for v in data.huaweicloud_global_eips.filter_by_eps_id.global_eips[*].enterprise_project_id : 
    v == huaweicloud_global_eip.test.enterprise_project_id]
}

output "is_eps_id_filter_useful" {
  value = length(local.filter_result_by_eps_id) > 0 && alltrue(local.filter_result_by_eps_id) 
}

// filter by tags
data "huaweicloud_global_eips" "filter_by_tags" {
  tags = huaweicloud_global_eip.test.tags
}

locals {
  filter_result_by_tags = [for v in data.huaweicloud_global_eips.filter_by_tags.global_eips[*].tags : 
    v == huaweicloud_global_eip.test.tags]
}

output "is_tags_filter_useful" {
  value = length(local.filter_result_by_tags) > 0 && alltrue(local.filter_result_by_tags) 
}
`, testAccGEIP_basic(rNameWithDash), rNameWithDash)
}
