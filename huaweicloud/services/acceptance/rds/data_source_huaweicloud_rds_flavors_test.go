package rds

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccRdsFlavorsDataSource_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rds_flavors.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFlavors_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.#"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.vcpus"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.memory"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.group_type"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.instance_mode"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.availability_zones.#"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.db_versions.#"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.az_status.%"),
					resource.TestCheckOutput("db_version_filter_is_useful", "true"),
					resource.TestCheckOutput("instance_mode_filter_is_useful", "true"),
					resource.TestCheckOutput("vcpus_filter_is_useful", "true"),
					resource.TestCheckOutput("memory_filter_is_useful", "true"),
					resource.TestCheckOutput("group_type_filter_is_useful", "true"),
					resource.TestCheckOutput("availability_zone_filter_is_useful", "true"),
					resource.TestCheckOutput("is_flexus_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccFlavors_basic() string {
	return `
data "huaweicloud_rds_flavors" "test" {
  db_type = "MySQL"
}

locals {
  db_version = "8.0"
}
data "huaweicloud_rds_flavors" "db_version_filter" {
  db_type    = "MySQL"
  db_version = "8.0"
}
output "db_version_filter_is_useful" {
  value = length(data.huaweicloud_rds_flavors.db_version_filter.flavors) > 0 && alltrue(
  [for v in data.huaweicloud_rds_flavors.db_version_filter.flavors[*].db_versions :
  alltrue([for vv in v : vv == local.db_version])]
  )
}

locals {
  instance_mode = "ha"
}
data "huaweicloud_rds_flavors" "instance_mode_filter" {
  db_type       = "MySQL"
  instance_mode = "ha"
}
output "instance_mode_filter_is_useful" {
  value = length(data.huaweicloud_rds_flavors.instance_mode_filter.flavors) > 0 && alltrue(
  [for v in data.huaweicloud_rds_flavors.instance_mode_filter.flavors[*].instance_mode : v == local.instance_mode]
  )
}

locals {
  vcpus = 4
}
data "huaweicloud_rds_flavors" "vcpus_filter" {
  db_type = "MySQL"
  vcpus   = 4
}
output "vcpus_filter_is_useful" {
  value = length(data.huaweicloud_rds_flavors.vcpus_filter.flavors) > 0 && alltrue(
  [for v in data.huaweicloud_rds_flavors.vcpus_filter.flavors[*].vcpus : v == local.vcpus]
  )
}

locals {
  memory = 4
}
data "huaweicloud_rds_flavors" "memory_filter" {
  db_type = "MySQL"
  memory  = 4
}
output "memory_filter_is_useful" {
  value = length(data.huaweicloud_rds_flavors.memory_filter.flavors) > 0 && alltrue(
  [for v in data.huaweicloud_rds_flavors.memory_filter.flavors[*].memory : v == local.memory]
  )
}

locals {
  group_type = "dedicated"
}
data "huaweicloud_rds_flavors" "group_type_filter" {
  db_type    = "MySQL"
  group_type = "dedicated"
}
output "group_type_filter_is_useful" {
  value = length(data.huaweicloud_rds_flavors.group_type_filter.flavors) > 0 && alltrue(
  [for v in data.huaweicloud_rds_flavors.group_type_filter.flavors[*].group_type : v == local.group_type]
  )
}

locals {
  availability_zone = data.huaweicloud_rds_flavors.test.flavors[0].availability_zones[0]
}
data "huaweicloud_rds_flavors" "availability_zone_filter" {
  db_type           = "MySQL"
  availability_zone = data.huaweicloud_rds_flavors.test.flavors[0].availability_zones[0]
}
output "availability_zone_filter_is_useful" {
  value = length(data.huaweicloud_rds_flavors.availability_zone_filter.flavors) > 0 && alltrue(
  [for v in data.huaweicloud_rds_flavors.availability_zone_filter.flavors[*].availability_zones :
  alltrue([for vv in v : vv == local.availability_zone])]
  )
}

locals {
  is_flexus = true
}
data "huaweicloud_rds_flavors" "is_flexus_filter" {
  db_type   = "MySQL"
  is_flexus = true
}
output "is_flexus_filter_is_useful" {
  value = length(data.huaweicloud_rds_flavors.is_flexus_filter.flavors) > 0
}
`
}
