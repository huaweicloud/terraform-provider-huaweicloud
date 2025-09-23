package dbss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDbssFlavors_basic(t *testing.T) {
	rName := "data.huaweicloud_dbss_flavors.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceDbssFlavors_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "flavors.0.id"),
					resource.TestCheckResourceAttrSet(rName, "flavors.0.level"),
					resource.TestCheckResourceAttrSet(rName, "flavors.0.proxy"),
					resource.TestCheckResourceAttrSet(rName, "flavors.0.vcpus"),
					resource.TestCheckResourceAttrSet(rName, "flavors.0.memory"),
					resource.TestCheckResourceAttrSet(rName, "flavors.0.availability_zones.0"),

					resource.TestCheckOutput("flavor_id_filter_is_useful", "true"),

					resource.TestCheckOutput("availability_zone_filter_is_useful", "true"),

					resource.TestCheckOutput("level_filter_is_useful", "true"),

					resource.TestCheckOutput("memory_filter_is_useful", "true"),

					resource.TestCheckOutput("vcpus_filter_is_useful", "true"),

					resource.TestCheckOutput("proxy_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceDbssFlavors_basic() string {
	return `
data "huaweicloud_dbss_flavors" "test" {}

locals {
  flavor_id = data.huaweicloud_dbss_flavors.test.flavors[0].id
}
data "huaweicloud_dbss_flavors" "flavor_id_filter" {
  flavor_id = local.flavor_id
}
output "flavor_id_filter_is_useful" {
  value = length(data.huaweicloud_dbss_flavors.flavor_id_filter.flavors) > 0 && alltrue(
    [for v in data.huaweicloud_dbss_flavors.flavor_id_filter.flavors[*].id : v == local.flavor_id]
  )  
}

locals {
  availability_zone = data.huaweicloud_dbss_flavors.test.flavors[0].availability_zones[0]
}
data "huaweicloud_dbss_flavors" "availability_zone_filter" {
  availability_zone = local.availability_zone
}
output "availability_zone_filter_is_useful" {
  value = length(data.huaweicloud_dbss_flavors.availability_zone_filter.flavors) > 0 && alltrue(
    [for v in data.huaweicloud_dbss_flavors.availability_zone_filter.flavors[*].availability_zones :
    contains(v, local.availability_zone)]
  )  
}

locals {
  level = data.huaweicloud_dbss_flavors.test.flavors[0].level
}
data "huaweicloud_dbss_flavors" "level_filter" {
  level = local.level
}
output "level_filter_is_useful" {
  value = length(data.huaweicloud_dbss_flavors.level_filter.flavors) > 0 && alltrue(
    [for v in data.huaweicloud_dbss_flavors.level_filter.flavors[*].level : v == local.level]
  )  
}

locals {
  memory = data.huaweicloud_dbss_flavors.test.flavors[0].memory
}
data "huaweicloud_dbss_flavors" "memory_filter" {
  memory = local.memory
}
output "memory_filter_is_useful" {
  value = length(data.huaweicloud_dbss_flavors.memory_filter.flavors) > 0 && alltrue(
    [for v in data.huaweicloud_dbss_flavors.memory_filter.flavors[*].memory : v == local.memory]
  )  
}

locals {
  vcpus = data.huaweicloud_dbss_flavors.test.flavors[0].vcpus
}
data "huaweicloud_dbss_flavors" "vcpus_filter" {
  vcpus = local.vcpus
}
output "vcpus_filter_is_useful" {
  value = length(data.huaweicloud_dbss_flavors.vcpus_filter.flavors) > 0 && alltrue(
    [for v in data.huaweicloud_dbss_flavors.vcpus_filter.flavors[*].vcpus : v == local.vcpus]
  )  
}

locals {
  proxy = data.huaweicloud_dbss_flavors.test.flavors[0].proxy
}
data "huaweicloud_dbss_flavors" "proxy_filter" {
  proxy = local.proxy
}
output "proxy_filter_is_useful" {
  value = length(data.huaweicloud_dbss_flavors.proxy_filter.flavors) > 0 && alltrue(
    [for v in data.huaweicloud_dbss_flavors.proxy_filter.flavors[*].proxy : v == local.proxy]
  )  
}
`
}
