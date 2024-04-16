package cc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCcBandwidthPackages_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cc_bandwidth_packages.filter_by_id"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)
	baseConfig := testDataSourceCcBandwidthPackages_base(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCcBandwidthPackages_basic(baseConfig, rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_packages.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_packages.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_packages.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_packages.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_packages.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_packages.0.updated_at"),
					resource.TestCheckOutput("is_id_useful", "true"),
					resource.TestCheckOutput("is_bandwidth_filter_useful", "true"),
					resource.TestCheckOutput("is_tags_filter_useful", "true"),
					resource.TestCheckOutput("is_name_and_tags_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceCcBandwidthPackages_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cc_bandwidth_package" "test1" {
  name           = "%[1]s1"
  local_area_id  = "Chinese-Mainland"
  remote_area_id = "Chinese-Mainland"
  charge_mode    = "bandwidth"
  billing_mode   = 3
  bandwidth      = 6
  description    = "desc 1"
	  
  tags = {
    foo   = "bar"
    owner = "terraform_test"
  }
}

resource "huaweicloud_cc_bandwidth_package" "test2" {
  name           = "%[1]s2"
  local_area_id  = "Chinese-Mainland"
  remote_area_id = "Chinese-Mainland"
  charge_mode    = "bandwidth"
  billing_mode   = 3
  bandwidth      = 12
  description    = "desc 2"
		
  tags = {
    foo = "bar"
  }
}
`, name)
}

func testDataSourceCcBandwidthPackages_basic(baseConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

locals {
  id = huaweicloud_cc_bandwidth_package.test1.id
}
  
data "huaweicloud_cc_bandwidth_packages" "filter_by_id" {
  bandwidth_package_id = local.id
	  
  depends_on = [
    huaweicloud_cc_bandwidth_package.test1,
    huaweicloud_cc_bandwidth_package.test2,
  ]
}
	
output "is_id_useful" {
  value = length(data.huaweicloud_cc_bandwidth_packages.filter_by_id.bandwidth_packages) >= 1 && alltrue(
    [for v in data.huaweicloud_cc_bandwidth_packages.filter_by_id.bandwidth_packages[*].id : v == local.id]
  )
}
	
locals {
  bandwidth = huaweicloud_cc_bandwidth_package.test1.bandwidth
}
	  
data "huaweicloud_cc_bandwidth_packages" "filter_by_bandwidth" {
  bandwidth = local.bandwidth
	
  depends_on = [
    huaweicloud_cc_bandwidth_package.test1,
    huaweicloud_cc_bandwidth_package.test2,
  ]
}
		
output "is_bandwidth_filter_useful" {
  value = length(data.huaweicloud_cc_bandwidth_packages.filter_by_bandwidth.bandwidth_packages) >= 1 && alltrue([
    for v in data.huaweicloud_cc_bandwidth_packages.filter_by_bandwidth.bandwidth_packages[*].bandwidth : 
      v == local.bandwidth
  ])
}
	
locals {
  tags = huaweicloud_cc_bandwidth_package.test1.tags
}
	  
data "huaweicloud_cc_bandwidth_packages" "filter_by_tags" {
  tags = local.tags
	
  depends_on = [
    huaweicloud_cc_bandwidth_package.test1,
    huaweicloud_cc_bandwidth_package.test2,
  ]
}
		
output "is_tags_filter_useful" {
  value = length(data.huaweicloud_cc_bandwidth_packages.filter_by_tags.bandwidth_packages) >= 1 && alltrue([
    for bp in data.huaweicloud_cc_bandwidth_packages.filter_by_tags.bandwidth_packages : alltrue([
      for k, v in local.tags : bp.tags[k] == v
    ])
  ])
}

locals {
  tags2 = huaweicloud_cc_bandwidth_package.test2.tags
  name  = huaweicloud_cc_bandwidth_package.test2.name
}
		
data "huaweicloud_cc_bandwidth_packages" "filter_by_name_and_tags" {
  tags = local.tags2
  name = local.name
	  
  depends_on = [
    huaweicloud_cc_bandwidth_package.test1,
    huaweicloud_cc_bandwidth_package.test2,
  ]
}
		  
output "is_name_and_tags_filter_useful" {
  value = length(data.huaweicloud_cc_bandwidth_packages.filter_by_name_and_tags.bandwidth_packages) >= 1 && alltrue([
    for bp in data.huaweicloud_cc_bandwidth_packages.filter_by_name_and_tags.bandwidth_packages : alltrue([
      for k, v in local.tags2 : bp.tags[k] == v
    ]) && bp.name == local.name
  ])
}
`, baseConfig, name)
}
