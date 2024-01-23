package eip

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccGlobalEIPPools_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_global_eip_pools.all"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGlobalEIPPoolsDataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "geip_pools.#"),
				),
			},
		},
	})
}

const testAccGlobalEIPPoolsDataSource_basic string = `data "huaweicloud_global_eip_pools" "all" {}`

func TestAccGlobalEIPPools_filter(t *testing.T) {
	byName := "data.huaweicloud_global_eip_pools.filter_by_name"
	byNameNotFound := "data.huaweicloud_global_eip_pools.filter_by_name_not_found"
	dcbyName := acceptance.InitDataSourceCheck(byName)
	dcbyNameNotFound := acceptance.InitDataSourceCheck(byNameNotFound)

	byAccessSite := "data.huaweicloud_global_eip_pools.filter_by_access_site"
	byAccessSiteNotFound := "data.huaweicloud_global_eip_pools.filter_by_access_site_not_found"
	dcbyAccessSite := acceptance.InitDataSourceCheck(byAccessSite)
	dcbyAccessSiteNotFound := acceptance.InitDataSourceCheck(byAccessSiteNotFound)

	byISP := "data.huaweicloud_global_eip_pools.filter_by_isp"
	byISPNotFound := "data.huaweicloud_global_eip_pools.filter_by_isp_not_found"
	dcbyISP := acceptance.InitDataSourceCheck(byISP)
	dcbyISPNotFound := acceptance.InitDataSourceCheck(byISPNotFound)

	byIPVersion := "data.huaweicloud_global_eip_pools.filter_by_ip_version"
	byIPVersionNotFound := "data.huaweicloud_global_eip_pools.filter_by_ip_version_not_found"
	dcbyIPVersion := acceptance.InitDataSourceCheck(byIPVersion)
	dcbyIPVersionNotFound := acceptance.InitDataSourceCheck(byIPVersionNotFound)

	byType := "data.huaweicloud_global_eip_pools.filter_by_type"
	byTypeNotFound := "data.huaweicloud_global_eip_pools.filter_by_type_not_found"
	dcbyType := acceptance.InitDataSourceCheck(byType)
	dcbyTypeNotFound := acceptance.InitDataSourceCheck(byTypeNotFound)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGlobalEIPPools_filter(),
				Check: resource.ComposeTestCheckFunc(

					dcbyName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					dcbyNameNotFound.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful_not_found", "true"),

					dcbyAccessSite.CheckResourceExists(),
					resource.TestCheckOutput("is_access_site_filter_useful", "true"),
					dcbyAccessSiteNotFound.CheckResourceExists(),
					resource.TestCheckOutput("is_access_site_filter_useful_not_found", "true"),

					dcbyISP.CheckResourceExists(),
					resource.TestCheckOutput("is_isp_filter_useful", "true"),
					dcbyISPNotFound.CheckResourceExists(),
					resource.TestCheckOutput("is_isp_filter_useful_not_found", "true"),

					dcbyIPVersion.CheckResourceExists(),
					resource.TestCheckOutput("is_ip_version_filter_useful", "true"),
					dcbyIPVersionNotFound.CheckResourceExists(),
					resource.TestCheckOutput("is_ip_version_filter_useful_not_found", "true"),

					dcbyType.CheckResourceExists(),
					resource.TestCheckOutput("is_type_filter_useful", "true"),
					dcbyTypeNotFound.CheckResourceExists(),
					resource.TestCheckOutput("is_type_filter_useful_not_found", "true"),
				),
			},
		},
	})
}

func testAccGlobalEIPPools_filter() string {
	rName := acceptance.RandomAccResourceName()
	rNameWithDash := acceptance.RandomAccResourceNameWithDash()
	return fmt.Sprintf(`
// filter by name
data "huaweicloud_global_eip_pools" "filter_by_name" {
  name = "bgp_default"
}

data "huaweicloud_global_eip_pools" "filter_by_name_not_found" {
  name = "%[1]s"
}

locals {
  filter_result_by_name = [for v in data.huaweicloud_global_eip_pools.filter_by_name.geip_pools[*].name : 
    v == "bgp_default"]
}

output "is_name_filter_useful" {
  value = length(local.filter_result_by_name) > 0 && alltrue(local.filter_result_by_name) 
}

output "is_name_filter_useful_not_found" {
  value = length(data.huaweicloud_global_eip_pools.filter_by_name_not_found.geip_pools) == 0
}

// filter by access site
data "huaweicloud_global_eip_pools" "filter_by_access_site" {
  access_site = "cn-south-guangzhou"
}

data "huaweicloud_global_eip_pools" "filter_by_access_site_not_found" {
  access_site = "%[2]s"
}

locals {
  filter_result_by_access_site = [for v in data.huaweicloud_global_eip_pools.filter_by_access_site.geip_pools[*].access_site : 
    v == "cn-south-guangzhou"]
}

output "is_access_site_filter_useful" {
  value = length(local.filter_result_by_access_site) > 0 && alltrue(local.filter_result_by_access_site) 
}

output "is_access_site_filter_useful_not_found" {
  value = length(data.huaweicloud_global_eip_pools.filter_by_access_site_not_found.geip_pools) == 0
}

// filter by isp
data "huaweicloud_global_eip_pools" "filter_by_isp" {
  isp = "Dyn_BGP"
}

data "huaweicloud_global_eip_pools" "filter_by_isp_not_found" {
  isp = "%[1]s"
}

locals {
  filter_result_by_isp = [for v in data.huaweicloud_global_eip_pools.filter_by_isp.geip_pools[*].isp : 
    v == "Dyn_BGP"]
}

output "is_isp_filter_useful" {
  value = length(local.filter_result_by_isp) > 0 && alltrue(local.filter_result_by_isp) 
}

output "is_isp_filter_useful_not_found" {
  value = length(data.huaweicloud_global_eip_pools.filter_by_isp_not_found.geip_pools) == 0
}

// filter by ip version
data "huaweicloud_global_eip_pools" "filter_by_ip_version" {
  ip_version = 4
}

data "huaweicloud_global_eip_pools" "filter_by_ip_version_not_found" {
  ip_version = 6
}

locals {
  filter_result_by_ip_version = [for v in data.huaweicloud_global_eip_pools.filter_by_ip_version.geip_pools[*].ip_version : 
    v == 4]
}

output "is_ip_version_filter_useful" {
  value = length(local.filter_result_by_ip_version) > 0 && alltrue(local.filter_result_by_ip_version) 
}

output "is_ip_version_filter_useful_not_found" {
  value = length(data.huaweicloud_global_eip_pools.filter_by_ip_version_not_found.geip_pools) == 0
}

// filter by type
data "huaweicloud_global_eip_pools" "filter_by_type" {
  type = "GEIP"
}

data "huaweicloud_global_eip_pools" "filter_by_type_not_found" {
  type = "GEIP_SEGMENT"
}

locals {
  filter_result_by_type = [for v in data.huaweicloud_global_eip_pools.filter_by_type.geip_pools[*].type : 
    v == "GEIP"]
}

output "is_type_filter_useful" {
  value = length(local.filter_result_by_type) > 0 && alltrue(local.filter_result_by_type) 
}

output "is_type_filter_useful_not_found" {
  value = length(data.huaweicloud_global_eip_pools.filter_by_type_not_found.geip_pools) == 0
}
`, rName, rNameWithDash)
}
