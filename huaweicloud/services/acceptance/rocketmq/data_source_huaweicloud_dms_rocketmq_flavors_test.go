package rocketmq

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceRocketmqFlavors_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_dms_rocketmq_flavors.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceDmsRocketmqFlavors_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "versions.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.#"),
					resource.TestCheckOutput("availability_zones_filter_is_useful", "true"),
					resource.TestCheckOutput("arch_type_filter_is_useful", "true"),
					resource.TestCheckOutput("charging_modes_is_useful", "true"),
					resource.TestCheckOutput("flavor_id_is_useful", "true"),
					resource.TestCheckOutput("storage_spec_code_filter_is_useful", "true"),
					resource.TestCheckOutput("type_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceDmsRocketmqFlavors_basic() string {
	return `
data "huaweicloud_dms_rocketmq_flavors" "test" {}

locals {
  flavor             = data.huaweicloud_dms_rocketmq_flavors.test.flavors[0]
  availability_zones = local.flavor.ios[0].availability_zones
  arch_type          = local.flavor.arch_types[0]
  charging_mode      = local.flavor.charging_modes[0]
  flavor_id          = local.flavor.id
  storage_spec_code  = local.flavor.ios[0].storage_spec_code
  type               = local.flavor.type
}

data "huaweicloud_dms_rocketmq_flavors" "availability_zones_filter" {
  availability_zones = local.availability_zones
}

output "availability_zones_filter_is_useful" {
  value = alltrue([for v in data.huaweicloud_dms_rocketmq_flavors.availability_zones_filter.flavors[*].ios :
  alltrue([for io in v : length(setintersection(io.availability_zones, local.availability_zones)) == length(local.availability_zones)])])
}

data "huaweicloud_dms_rocketmq_flavors" "arch_type_filter" {
  arch_type = local.arch_type
}

output "arch_type_filter_is_useful" {
  value = alltrue([for v in data.huaweicloud_dms_rocketmq_flavors.arch_type_filter.flavors[*].arch_types : contains(v, local.arch_type)])
}

data "huaweicloud_dms_rocketmq_flavors" "charging_mode_filter" {
  charging_mode = local.charging_mode
}

output "charging_modes_is_useful" {
  value = alltrue([for v in data.huaweicloud_dms_rocketmq_flavors.charging_mode_filter.flavors[*].charging_modes : contains(v, local.charging_mode)])
}

data "huaweicloud_dms_rocketmq_flavors" "flavor_id_filter" {
  flavor_id = local.flavor_id
}

output "flavor_id_is_useful" {
  value = alltrue([for v in data.huaweicloud_dms_rocketmq_flavors.flavor_id_filter.flavors[*].id : v == local.flavor_id])
}

data "huaweicloud_dms_rocketmq_flavors" "storage_spec_code_filter" {
  storage_spec_code = local.storage_spec_code
}

output "storage_spec_code_filter_is_useful" {
  value = alltrue([for v in data.huaweicloud_dms_rocketmq_flavors.storage_spec_code_filter.flavors[*].ios :
  alltrue([for io in v : io.storage_spec_code == local.storage_spec_code])])
}

data "huaweicloud_dms_rocketmq_flavors" "type_filter" { 
  type = local.type
}

output "type_filter_is_useful" {
  value = alltrue([for v in data.huaweicloud_dms_rocketmq_flavors.type_filter.flavors[*].type : v == local.type])
}
`
}
