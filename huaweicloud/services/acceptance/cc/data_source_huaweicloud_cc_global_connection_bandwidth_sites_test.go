package cc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCcGlobalConnectionBandwidthSites_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cc_global_connection_bandwidth_sites.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCcGlobalConnectionBandwidthSites_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "site_infos.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "site_infos.0.site_code"),
					resource.TestCheckResourceAttrSet(dataSource, "site_infos.0.site_type"),
					resource.TestCheckResourceAttrSet(dataSource, "site_infos.0.name_en"),
					resource.TestCheckResourceAttrSet(dataSource, "site_infos.0.name_cn"),

					resource.TestCheckOutput("id_filter_is_useful", "true"),
					resource.TestCheckOutput("site_code_filter_is_useful", "true"),
					resource.TestCheckOutput("site_type_filter_is_useful", "true"),
					resource.TestCheckOutput("name_en_filter_is_useful", "true"),
					resource.TestCheckOutput("name_cn_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceCcGlobalConnectionBandwidthSites_basic() string {
	return `
data "huaweicloud_cc_global_connection_bandwidth_sites" "test" {}

locals {
  id      = data.huaweicloud_cc_global_connection_bandwidth_sites.test.site_infos[0].id
  code    = data.huaweicloud_cc_global_connection_bandwidth_sites.test.site_infos[0].site_code
  type    = data.huaweicloud_cc_global_connection_bandwidth_sites.test.site_infos[0].site_type
  name_en = data.huaweicloud_cc_global_connection_bandwidth_sites.test.site_infos[0].name_en
  name_cn = data.huaweicloud_cc_global_connection_bandwidth_sites.test.site_infos[0].name_cn
}
	
data "huaweicloud_cc_global_connection_bandwidth_sites" "filter_by_id" {
  site_id = local.id
}
	
data "huaweicloud_cc_global_connection_bandwidth_sites" "filter_by_code" {
  site_code = local.code
}
	
data "huaweicloud_cc_global_connection_bandwidth_sites" "filter_by_type" {
  site_type = local.type
}
	
data "huaweicloud_cc_global_connection_bandwidth_sites" "filter_by_name_en" {
  name_en   = local.name_en
}
	
data "huaweicloud_cc_global_connection_bandwidth_sites" "filter_by_name_cn" {
  name_cn   = local.name_cn
}
	
locals {
  list_by_id      = data.huaweicloud_cc_global_connection_bandwidth_sites.filter_by_id.site_infos
  list_by_code    = data.huaweicloud_cc_global_connection_bandwidth_sites.filter_by_code.site_infos
  list_by_type    = data.huaweicloud_cc_global_connection_bandwidth_sites.filter_by_type.site_infos
  list_by_name_en = data.huaweicloud_cc_global_connection_bandwidth_sites.filter_by_name_en.site_infos
  list_by_name_cn = data.huaweicloud_cc_global_connection_bandwidth_sites.filter_by_name_cn.site_infos
}
output "id_filter_is_useful" {
  value = length(local.list_by_id) > 0 && alltrue([for v in local.list_by_id[*].id : v == local.id])
}
	
output "site_code_filter_is_useful" {
  value = length(local.list_by_code) > 0 && alltrue([for v in local.list_by_code[*].site_code : v == local.code])
}
	
output "site_type_filter_is_useful" {
  value = length(local.list_by_type) > 0 && alltrue([for v in local.list_by_type[*].site_type : v == local.type])
}
	
output "name_en_filter_is_useful" {
  value = length(local.list_by_name_en) > 0 && alltrue(
    [for v in local.list_by_name_en[*].name_en : v == local.name_en]
  )
}
	
output "name_cn_filter_is_useful" {
  value = length(local.list_by_name_cn) > 0 && alltrue(
    [for v in local.list_by_name_cn[*].name_cn : v == local.name_cn]
  )
}
`
}
