package eip

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccGlobalEIPAccessSites_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_global_eip_access_sites.all"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGlobalEIPAccessSitesDataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "access_sites.#"),
				),
			},
		},
	})
}

const testAccGlobalEIPAccessSitesDataSource_basic string = `data "huaweicloud_global_eip_access_sites" "all" {}`

func TestAccGlobalEIPAccessSites_filter(t *testing.T) {
	byName := "data.huaweicloud_global_eip_access_sites.filter_by_name"
	byNameNotFound := "data.huaweicloud_global_eip_access_sites.filter_by_name_not_found"
	dcbyName := acceptance.InitDataSourceCheck(byName)
	dcbyNameNotFound := acceptance.InitDataSourceCheck(byNameNotFound)

	byProxyRegion := "data.huaweicloud_global_eip_access_sites.filter_by_proxy_region"
	byProxyRegionNotFound := "data.huaweicloud_global_eip_access_sites.filter_by_proxy_region_not_found"
	dcbyProxyRegion := acceptance.InitDataSourceCheck(byProxyRegion)
	dcbyProxyRegionNotFound := acceptance.InitDataSourceCheck(byProxyRegionNotFound)

	byIecAZCode := "data.huaweicloud_global_eip_access_sites.filter_by_iec_az_code"
	byIecAZCodeNotFound := "data.huaweicloud_global_eip_access_sites.filter_by_iec_az_code_not_found"
	dcbyIecAZCode := acceptance.InitDataSourceCheck(byIecAZCode)
	dcbyIecAZCodeNotFound := acceptance.InitDataSourceCheck(byIecAZCodeNotFound)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGlobalEIPAccessSites_filter(),
				Check: resource.ComposeTestCheckFunc(

					dcbyName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					dcbyNameNotFound.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful_not_found", "true"),

					dcbyProxyRegion.CheckResourceExists(),
					resource.TestCheckOutput("is_proxy_region_filter_useful", "true"),
					dcbyProxyRegionNotFound.CheckResourceExists(),
					resource.TestCheckOutput("is_proxy_region_filter_useful_not_found", "true"),

					dcbyIecAZCode.CheckResourceExists(),
					resource.TestCheckOutput("is_iec_az_code_filter_useful", "true"),
					dcbyIecAZCodeNotFound.CheckResourceExists(),
					resource.TestCheckOutput("is_iec_az_code_filter_useful_not_found", "true"),
				),
			},
		},
	})
}

func testAccGlobalEIPAccessSites_filter() string {
	rNameWithDash := acceptance.RandomAccResourceNameWithDash()

	return fmt.Sprintf(`
// filter by name
data "huaweicloud_global_eip_access_sites" "filter_by_name" {
  name = "cn-south-guangzhou"
}

data "huaweicloud_global_eip_access_sites" "filter_by_name_not_found" {
  name = "%[1]s"
}

locals {
  filter_result_by_name = [for v in data.huaweicloud_global_eip_access_sites.filter_by_name.access_sites[*].name : 
    v == "cn-south-guangzhou"]
}

output "is_name_filter_useful" {
  value = length(local.filter_result_by_name) > 0 && alltrue(local.filter_result_by_name) 
}

output "is_name_filter_useful_not_found" {
  value = length(data.huaweicloud_global_eip_access_sites.filter_by_name_not_found.access_sites) == 0
}

// filter by proxy_region
data "huaweicloud_global_eip_access_sites" "filter_by_proxy_region" {
  proxy_region = "cn-south-1"
}

data "huaweicloud_global_eip_access_sites" "filter_by_proxy_region_not_found" {
  proxy_region = "%[1]s"
}

locals {
  filter_result_by_proxy_region = [for v in data.huaweicloud_global_eip_access_sites.filter_by_proxy_region.access_sites[*].proxy_region : 
    v == "cn-south-1"]
}

output "is_proxy_region_filter_useful" {
  value = length(local.filter_result_by_proxy_region) > 0 && alltrue(local.filter_result_by_proxy_region) 
}

output "is_proxy_region_filter_useful_not_found" {
  value = length(data.huaweicloud_global_eip_access_sites.filter_by_proxy_region_not_found.access_sites) == 0
}

// filter by iec_az_code
data "huaweicloud_global_eip_access_sites" "filter_by_iec_az_code" {
  iec_az_code = "cn-north-900d"
}

data "huaweicloud_global_eip_access_sites" "filter_by_iec_az_code_not_found" {
  iec_az_code = "%[1]s"
}

locals {
  filter_result_by_iec_az_code = [for v in data.huaweicloud_global_eip_access_sites.filter_by_iec_az_code.access_sites[*].iec_az_code : 
    v == "cn-north-900d"]
}

output "is_iec_az_code_filter_useful" {
  value = length(local.filter_result_by_iec_az_code) > 0 && alltrue(local.filter_result_by_iec_az_code) 
}

output "is_iec_az_code_filter_useful_not_found" {
  value = length(data.huaweicloud_global_eip_access_sites.filter_by_iec_az_code_not_found.access_sites) == 0
}
`, rNameWithDash)
}
