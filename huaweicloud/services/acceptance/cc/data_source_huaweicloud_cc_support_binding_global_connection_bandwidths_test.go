package cc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCcSupportBindingGlobalConnectionBandwidths_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cc_support_binding_global_connection_bandwidths.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCcSupportBindingGlobalConnectionBandwidths_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "globalconnection_bandwidths.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "globalconnection_bandwidths.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "globalconnection_bandwidths.0.local_site_code"),
					resource.TestCheckResourceAttrSet(dataSource, "globalconnection_bandwidths.0.remote_site_code"),
					resource.TestCheckResourceAttrSet(dataSource, "globalconnection_bandwidths.0.size"),
					resource.TestCheckResourceAttrSet(dataSource, "globalconnection_bandwidths.0.charge_mode"),
					resource.TestCheckResourceAttrSet(dataSource, "globalconnection_bandwidths.0.local_area"),
					resource.TestCheckResourceAttrSet(dataSource, "globalconnection_bandwidths.0.remote_area"),
					resource.TestCheckResourceAttrSet(dataSource, "globalconnection_bandwidths.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "globalconnection_bandwidths.0.sla_level"),
					resource.TestCheckResourceAttrSet(dataSource, "globalconnection_bandwidths.0.domain_id"),
					resource.TestCheckOutput("is_default_filter_useful", "true"),
					resource.TestCheckOutput("is_area_filter_useful", "true"),
					resource.TestCheckOutput("is_id_filter_useful", "true"),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("is_size_filter_useful", "true"),
					resource.TestCheckOutput("is_charge_mode_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceCcSupportBindingGlobalConnectionBandwidths_basic(name string) string {
	return fmt.Sprintf(`
%s

locals {
  id              = huaweicloud_cc_global_connection_bandwidth.gcb2.id
  name            = huaweicloud_cc_global_connection_bandwidth.gcb2.name
  size            = huaweicloud_cc_global_connection_bandwidth.gcb2.size
  charge_mode     = huaweicloud_cc_global_connection_bandwidth.gcb1.charge_mode
  local_area      = huaweicloud_cc_global_connection_bandwidth.gcb1.local_area
  remote_area     = huaweicloud_cc_global_connection_bandwidth.gcb1.remote_area
  binding_service = "GCN"
}

data "huaweicloud_cc_support_binding_global_connection_bandwidths" "test" {
  binding_service = local.binding_service
	
  depends_on = [
    huaweicloud_cc_global_connection_bandwidth.gcb1,
    huaweicloud_cc_global_connection_bandwidth.gcb2,
  ]
}
  
output "is_default_filter_useful" {
  value = length(data.huaweicloud_cc_support_binding_global_connection_bandwidths.test.globalconnection_bandwidths) >= 2
}

data "huaweicloud_cc_support_binding_global_connection_bandwidths" "filter_by_area" {
  binding_service = local.binding_service
  local_area      = local.local_area
  remote_area     = local.remote_area

  depends_on = [
    huaweicloud_cc_global_connection_bandwidth.gcb1,
    huaweicloud_cc_global_connection_bandwidth.gcb2,
  ]
}

output "is_area_filter_useful" {
  value = length(data.huaweicloud_cc_support_binding_global_connection_bandwidths.filter_by_area.globalconnection_bandwidths) >= 1 && alltrue([
    for v in data.huaweicloud_cc_support_binding_global_connection_bandwidths.filter_by_area.globalconnection_bandwidths[*] : 
      v.local_site_code == local.local_area && v.remote_site_code == local.remote_area
  ])
}

data "huaweicloud_cc_support_binding_global_connection_bandwidths" "filter_by_id" {
  binding_service = local.binding_service
  gcb_id          = local.id

  depends_on = [
    huaweicloud_cc_global_connection_bandwidth.gcb1,
    huaweicloud_cc_global_connection_bandwidth.gcb2,
  ]
}

output "is_id_filter_useful" {
  value = length(data.huaweicloud_cc_support_binding_global_connection_bandwidths.filter_by_id.globalconnection_bandwidths) >= 1 && alltrue(
    [for v in data.huaweicloud_cc_support_binding_global_connection_bandwidths.filter_by_id.globalconnection_bandwidths[*] : v.id == local.id]
  )
}

data "huaweicloud_cc_support_binding_global_connection_bandwidths" "filter_by_name" {
  binding_service = local.binding_service
  name            = local.name

  depends_on = [
    huaweicloud_cc_global_connection_bandwidth.gcb1,
    huaweicloud_cc_global_connection_bandwidth.gcb2,
  ]
}

output "is_name_filter_useful" {
  value = length(data.huaweicloud_cc_support_binding_global_connection_bandwidths.filter_by_name.globalconnection_bandwidths) >= 1 && alltrue(
    [for v in data.huaweicloud_cc_support_binding_global_connection_bandwidths.filter_by_name.globalconnection_bandwidths[*] : v.name == local.name]
  )
}

data "huaweicloud_cc_support_binding_global_connection_bandwidths" "filter_by_size" {
  binding_service = local.binding_service
  size            = local.size

  depends_on = [
    huaweicloud_cc_global_connection_bandwidth.gcb1,
    huaweicloud_cc_global_connection_bandwidth.gcb2,
  ]
}

output "is_size_filter_useful" {
  value = length(data.huaweicloud_cc_support_binding_global_connection_bandwidths.filter_by_size.globalconnection_bandwidths) >= 1 && alltrue(
    [for v in data.huaweicloud_cc_support_binding_global_connection_bandwidths.filter_by_size.globalconnection_bandwidths[*] : v.size == local.size]
  )
}

data "huaweicloud_cc_support_binding_global_connection_bandwidths" "filter_by_charge_mode" {
  binding_service = local.binding_service
  charge_mode     = local.charge_mode

  depends_on = [
    huaweicloud_cc_global_connection_bandwidth.gcb1,
    huaweicloud_cc_global_connection_bandwidth.gcb2,
  ]
}

output "is_charge_mode_filter_useful" {
  value = length(data.huaweicloud_cc_support_binding_global_connection_bandwidths.filter_by_charge_mode.globalconnection_bandwidths) >= 1 && alltrue([
    for v in data.huaweicloud_cc_support_binding_global_connection_bandwidths.filter_by_charge_mode.globalconnection_bandwidths[*] : 
      v.charge_mode == local.charge_mode
  ])
}
`, testDataSourceCcSupportBindingGlobalConnectionBandwidths_base(name))
}

func testDataSourceCcSupportBindingGlobalConnectionBandwidths_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cc_global_connection_bandwidth" "gcb1" {
  name        = "%[1]s1"
  type        = "Area"  
  bordercross = false
  charge_mode = "95"
  size        = 10
  description = "test 1"
  local_area  = "cn-north-beijing4"
  remote_area = "cn-south-guangzhou"

  tags = {
    foo = "bar"
  }
}

resource "huaweicloud_cc_global_connection_bandwidth" "gcb2" {
  name        = "%[1]s2"
  type        = "Area"  
  bordercross = false
  charge_mode = "bwd"
  size        = 12
  description = "test 2"
  local_area  = "cn-south-guangzhou"
  remote_area = "cn-southwest-guiyang1"

  tags = {
    sam = "lib"
  }
}
`, name)
}
