package kafka

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccKafkaFlavorsDataSource_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_dms_kafka_flavors.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKafkaFlavorsDataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSourceName, "versions.#", regexp.MustCompile(`[1-9]\d*`)),
					resource.TestMatchResourceAttr(dataSourceName, "flavors.#", regexp.MustCompile(`[1-9]\d*`)),
					resource.TestCheckOutput("type_validation", "true"),
					resource.TestCheckOutput("arch_types_validation", "true"),
					resource.TestCheckOutput("charging_modes_validation", "true"),
					resource.TestCheckOutput("storage_spec_code_validation", "true"),
					resource.TestCheckOutput("availability_zones_validation", "true"),
				),
			},
		},
	})
}

const testAccKafkaFlavorsDataSource_basic = `
data "huaweicloud_dms_kafka_flavors" "baisc" {
  type = "cluster"
}

data "huaweicloud_dms_kafka_flavors" "test" {
  type               = local.test_refer.type
  arch_type          = local.test_refer.arch_types[0]
  charging_mode      = local.test_refer.charging_modes[0]
  storage_spec_code  = local.test_refer.ios[0].storage_spec_code
  availability_zones = local.test_refer.ios[0].availability_zones
}

locals {
  test_refer   = data.huaweicloud_dms_kafka_flavors.baisc.flavors[0]
  test_results = data.huaweicloud_dms_kafka_flavors.test
}

output "type_validation" {
  value = contains(local.test_results.flavors[*].type, local.test_refer.type)
}

output "arch_types_validation" {
  value = !contains([for a in local.test_results.flavors[*].arch_types : contains(a, local.test_refer.arch_types[0])], false)
}

output "charging_modes_validation" {
  value = !contains([for c in local.test_results.flavors[*].charging_modes : contains(c, local.test_refer.charging_modes[0])], false)
}

output "storage_spec_code_validation" {
  value = !contains([for ios in local.test_results.flavors[*].ios : !contains([for io in ios : io.storage_spec_code == local.test_refer.ios[0].storage_spec_code], false)], false)
}

output "availability_zones_validation" {
  value = !contains([for ios in local.test_results.flavors[*].ios : !contains([for io in ios : length(setintersection(io.availability_zones, local.test_refer.ios[0].availability_zones)) == length(local.test_refer.ios[0].availability_zones)], false)], false)
}
`
