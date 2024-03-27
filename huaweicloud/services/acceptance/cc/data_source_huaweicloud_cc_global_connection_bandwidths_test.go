package cc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceGlobalConnectionBandwidths_basic(t *testing.T) {
	defaultDataSourceName := "data.huaweicloud_cc_global_connection_bandwidths.filter_by_name"
	dc := acceptance.InitDataSourceCheck(defaultDataSourceName)
	name := acceptance.RandomAccResourceName()
	baseConfig := testAccDatasourceGlobalConnectionBandwidths_base(name)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceGlobalConnectionBandwidths_basic(baseConfig, name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(defaultDataSourceName, "globalconnection_bandwidths.0.name"),
					resource.TestCheckResourceAttrSet(defaultDataSourceName, "globalconnection_bandwidths.0.id"),
					resource.TestCheckResourceAttrSet(defaultDataSourceName, "globalconnection_bandwidths.0.size"),
					resource.TestCheckResourceAttrSet(defaultDataSourceName, "globalconnection_bandwidths.0.charge_mode"),
					resource.TestCheckOutput("is_name_useful", "true"),
					resource.TestCheckOutput("is_gcb_id_filter_useful", "true"),
					resource.TestCheckOutput("is_size_filter_useful", "true"),
					resource.TestCheckOutput("is_tags_filter_useful", "true"),
					resource.TestCheckOutput("is_tags_charge_mode_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceGlobalConnectionBandwidths_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cc_global_connection_bandwidth" "gcb1" {
  name        = "%[1]s_1"
  type        = "Region"  
  bordercross = false
  charge_mode = "bwd"
  size        = 5
  description = "test"
  sla_level   = "Ag"

  tags = {
    foo = "bar"
  }
}

resource "huaweicloud_cc_global_connection_bandwidth" "gcb2" {
  name        = "%[1]s_2"
  type        = "Region"  
  bordercross = false
  charge_mode = "95"
  size        = 5
  description = "test"
  sla_level   = "Ag"
  
  tags = {
    foo = "bar"
  }

  depends_on = [
    huaweicloud_cc_global_connection_bandwidth.gcb1
  ]
}

resource "huaweicloud_cc_global_connection_bandwidth" "gcb3" {
  name        = "%[1]s_3"
  type        = "Region"  
  bordercross = false
  charge_mode = "bwd"
  size        = 10
  description = "test"
  sla_level   = "Ag"
  
  tags = {
    sam = "lib"
  }

  depends_on = [
    huaweicloud_cc_global_connection_bandwidth.gcb1,
    huaweicloud_cc_global_connection_bandwidth.gcb2,
  ]
}
`, name)
}

func testAccDatasourceGlobalConnectionBandwidths_basic(config, name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_cc_global_connection_bandwidths" "filter_by_name" {
  name = "%[2]s_1"
  
  depends_on = [
    huaweicloud_cc_global_connection_bandwidth.gcb1,
    huaweicloud_cc_global_connection_bandwidth.gcb2,
    huaweicloud_cc_global_connection_bandwidth.gcb3
  ]
}

output "is_name_useful" {
  value = length(data.huaweicloud_cc_global_connection_bandwidths.filter_by_name.globalconnection_bandwidths) >= 1 && alltrue(
    [for v in data.huaweicloud_cc_global_connection_bandwidths.filter_by_name.globalconnection_bandwidths[*].name : v == "%[2]s_1"]
  )
}

locals {
  gcb_id = huaweicloud_cc_global_connection_bandwidth.gcb1.id
}
  
data "huaweicloud_cc_global_connection_bandwidths" "filter_by_gcb_id" {
  gcb_id = local.gcb_id

  depends_on = [
    huaweicloud_cc_global_connection_bandwidth.gcb1,
    huaweicloud_cc_global_connection_bandwidth.gcb2,
    huaweicloud_cc_global_connection_bandwidth.gcb3
  ]
}
	
output "is_gcb_id_filter_useful" {
  value = length(data.huaweicloud_cc_global_connection_bandwidths.filter_by_gcb_id.globalconnection_bandwidths) == 1 && alltrue(
    [for v in data.huaweicloud_cc_global_connection_bandwidths.filter_by_gcb_id.globalconnection_bandwidths[*].id : v == local.gcb_id]
  )
}

locals {
  size = huaweicloud_cc_global_connection_bandwidth.gcb1.size
}
  
data "huaweicloud_cc_global_connection_bandwidths" "filter_by_size" {
  size = local.size

  depends_on = [
    huaweicloud_cc_global_connection_bandwidth.gcb1,
    huaweicloud_cc_global_connection_bandwidth.gcb2,
    huaweicloud_cc_global_connection_bandwidth.gcb3
  ]
}
	
output "is_size_filter_useful" {
  value = length(data.huaweicloud_cc_global_connection_bandwidths.filter_by_size.globalconnection_bandwidths) >= 1 && alltrue(
    [for v in data.huaweicloud_cc_global_connection_bandwidths.filter_by_size.globalconnection_bandwidths[*].size : v == local.size]
  )
}

locals {
  tags = huaweicloud_cc_global_connection_bandwidth.gcb1.tags
}
  
data "huaweicloud_cc_global_connection_bandwidths" "filter_by_tags" {
  tags = local.tags

  depends_on = [
    huaweicloud_cc_global_connection_bandwidth.gcb1,
    huaweicloud_cc_global_connection_bandwidth.gcb2,
    huaweicloud_cc_global_connection_bandwidth.gcb3
  ]
}
	
output "is_tags_filter_useful" {
  value = length(data.huaweicloud_cc_global_connection_bandwidths.filter_by_tags.globalconnection_bandwidths) >= 1 && alltrue([
    for bw in data.huaweicloud_cc_global_connection_bandwidths.filter_by_tags.globalconnection_bandwidths : alltrue([
      for k, v in local.tags : bw.tags[k] == v
    ])
  ])
}

locals {
  charge_mode = huaweicloud_cc_global_connection_bandwidth.gcb2.charge_mode
  tags2       = huaweicloud_cc_global_connection_bandwidth.gcb2.tags
}

data "huaweicloud_cc_global_connection_bandwidths" "filter_by_tags_and_charge_mode" {
  tags        = local.tags2
  charge_mode = local.charge_mode

  depends_on = [
    huaweicloud_cc_global_connection_bandwidth.gcb1,
    huaweicloud_cc_global_connection_bandwidth.gcb2,
    huaweicloud_cc_global_connection_bandwidth.gcb3
  ]
}

output "is_tags_charge_mode_filter_useful" {
  value = length(data.huaweicloud_cc_global_connection_bandwidths.filter_by_tags_and_charge_mode.globalconnection_bandwidths) >= 1 && alltrue([
    for bw in data.huaweicloud_cc_global_connection_bandwidths.filter_by_tags_and_charge_mode.globalconnection_bandwidths : alltrue([
      for k, v in local.tags2 : bw.tags[k] == v
    ]) && bw.charge_mode == local.charge_mode
  ])
}
`, config, name)
}
