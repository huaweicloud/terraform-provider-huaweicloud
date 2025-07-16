package dcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDcsFlavors_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dcs_flavors.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDcsFlavors_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.#"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.cache_mode"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.engine"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.engine_versions"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.cpu_architecture"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.capacity"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.available_zones.#"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.charging_modes.#"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.ip_count"),
					resource.TestCheckOutput("instance_id_filter_is_useful", "true"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("capacity_filter_is_useful", "true"),
					resource.TestCheckOutput("engine_filter_is_useful", "true"),
					resource.TestCheckOutput("engine_version_filter_is_useful", "true"),
					resource.TestCheckOutput("cache_mode_filter_is_useful", "true"),
					resource.TestCheckOutput("ccpu_architecture_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDcsFlavors_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_dcs_flavors" "flavors" {
  cache_mode       = "ha"
  capacity         = 1
  cpu_architecture = "x86_64"
}

resource "huaweicloud_dcs_instance" "test" {
  name               = "%[1]s"
  engine_version     = "5.0"
  password           = "Huawei_test"
  engine             = "Redis"
  capacity           = 1
  vpc_id             = data.huaweicloud_vpc.test.id
  subnet_id          = data.huaweicloud_vpc_subnet.test.id
  availability_zones = [data.huaweicloud_availability_zones.test.names[0]]
  flavor             = data.huaweicloud_dcs_flavors.flavors.flavors[0].name
}
`, name)
}

func testAccDcsFlavors_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dcs_flavors" "test" {}

locals {
  instance_id = huaweicloud_dcs_instance.test.id
}
data "huaweicloud_dcs_flavors" "instance_id_filter" {
  instance_id = huaweicloud_dcs_instance.test.id
}
output "instance_id_filter_is_useful" {
  value = length(data.huaweicloud_dcs_flavors.instance_id_filter.flavors) > 0
}

locals {
  name = data.huaweicloud_dcs_flavors.test.flavors[0].name
}
data "huaweicloud_dcs_flavors" "name_filter" {
  name = data.huaweicloud_dcs_flavors.test.flavors[0].name
}
output "name_filter_is_useful" {
  value = length(data.huaweicloud_dcs_flavors.name_filter.flavors) > 0 && alltrue(
    [for v in data.huaweicloud_dcs_flavors.name_filter.flavors[*].name : v == local.name]
  )
}

locals {
  capacity = data.huaweicloud_dcs_flavors.test.flavors[0].capacity
}
data "huaweicloud_dcs_flavors" "capacity_filter" {
  capacity = data.huaweicloud_dcs_flavors.test.flavors[0].capacity
}
output "capacity_filter_is_useful" {
  value = length(data.huaweicloud_dcs_flavors.capacity_filter.flavors) > 0 && alltrue(
    [for v in data.huaweicloud_dcs_flavors.capacity_filter.flavors[*].capacity : v == local.capacity]
  )
}

locals {
  engine = data.huaweicloud_dcs_flavors.test.flavors[0].engine
}
data "huaweicloud_dcs_flavors" "engine_filter" {
  engine = data.huaweicloud_dcs_flavors.test.flavors[0].engine
}
output "engine_filter_is_useful" {
  value = length(data.huaweicloud_dcs_flavors.engine_filter.flavors) > 0 && alltrue(
    [for v in data.huaweicloud_dcs_flavors.engine_filter.flavors[*].engine : v == local.engine]
  )
}

locals {
  engine_version = "5.0"
}
data "huaweicloud_dcs_flavors" "engine_version_filter" {
  engine_version = "5.0"
}
output "engine_version_filter_is_useful" {
  value = length(data.huaweicloud_dcs_flavors.engine_version_filter.flavors) > 0
}

locals {
  cache_mode = data.huaweicloud_dcs_flavors.test.flavors[0].cache_mode
}
data "huaweicloud_dcs_flavors" "cache_mode_filter" {
  cache_mode = data.huaweicloud_dcs_flavors.test.flavors[0].cache_mode
}
output "cache_mode_filter_is_useful" {
  value = length(data.huaweicloud_dcs_flavors.cache_mode_filter.flavors) > 0 && alltrue(
    [for v in data.huaweicloud_dcs_flavors.cache_mode_filter.flavors[*].cache_mode : v == local.cache_mode]
  )
}

locals {
  cpu_architecture = data.huaweicloud_dcs_flavors.test.flavors[0].cpu_architecture
}
data "huaweicloud_dcs_flavors" "cpu_architecture_filter" {
  cpu_architecture = data.huaweicloud_dcs_flavors.test.flavors[0].cpu_architecture
}
output "ccpu_architecture_filter_is_useful" {
  value = length(data.huaweicloud_dcs_flavors.cpu_architecture_filter.flavors) > 0 && alltrue(
    [for v in data.huaweicloud_dcs_flavors.cpu_architecture_filter.flavors[*].cpu_architecture : v == local.cpu_architecture]
  )
}
`, testAccDcsFlavors_base(name))
}
