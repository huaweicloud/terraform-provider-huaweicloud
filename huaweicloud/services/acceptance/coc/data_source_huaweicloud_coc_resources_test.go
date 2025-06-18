package coc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCocResources_basic(t *testing.T) {
	dataSource := "data.huaweicloud_coc_resources.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCocResources_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.resource_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.ep_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.domain_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.cloud_service_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.region_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.tags.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.tags.0.key"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.tags.0.value"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.properties.%"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.is_delegated"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("ep_id_filter_is_useful", "true"),
					resource.TestCheckOutput("project_id_filter_is_useful", "true"),
					resource.TestCheckOutput("region_id_filter_is_useful", "true"),
					resource.TestCheckOutput("az_id_filter_is_useful", "true"),
					resource.TestCheckOutput("ip_filter_is_useful", "true"),
					resource.TestCheckOutput("resource_id_list_filter_is_useful", "true"),
					resource.TestCheckOutput("status_filter_is_useful", "true"),
					resource.TestCheckOutput("agent_state_filter_is_useful", "true"),
					resource.TestCheckOutput("image_name_filter_is_useful", "true"),
					resource.TestCheckOutput("os_type_filter_is_useful", "true"),
					resource.TestCheckOutput("tag_filter_is_useful", "true"),
					resource.TestCheckOutput("tag_key_filter_is_useful", "true"),
					resource.TestCheckOutput("vpc_id_filter_is_useful", "true"),
					resource.TestCheckOutput("flavor_name_filter_is_useful", "true"),
					resource.TestCheckOutput("charging_mode_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceCocResources_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_coc_resources" "test" {
  cloud_service_name = "ecs"
  type               = "cloudservers"
  depends_on         = [huaweicloud_compute_instance.test]
}

data "huaweicloud_coc_resources" "name_filter" {
  cloud_service_name = "ecs"
  type               = "cloudservers"
  name               = huaweicloud_compute_instance.test.name
}

output "name_filter_is_useful" {
  value = length(data.huaweicloud_coc_resources.name_filter.data) > 0 && alltrue(
    [for v in data.huaweicloud_coc_resources.name_filter.data[*].name : v == huaweicloud_compute_instance.test.name]
  )
}

data "huaweicloud_coc_resources" "ep_id_filter" {
  cloud_service_name = "ecs"
  type               = "cloudservers"
  ep_id              = "0"
  depends_on         = [huaweicloud_compute_instance.test]
}

output "ep_id_filter_is_useful" {
  value = length(data.huaweicloud_coc_resources.ep_id_filter.data) > 0 && alltrue(
    [for v in data.huaweicloud_coc_resources.ep_id_filter.data[*].ep_id : v == "0"]
  )
}

data "huaweicloud_coc_resources" "region_id_filter" {
  cloud_service_name = "ecs"
  type               = "cloudservers"
  region_id          = huaweicloud_compute_instance.test.region
}

output "region_id_filter_is_useful" {
  value = length(data.huaweicloud_coc_resources.region_id_filter.data) > 0 && alltrue(
    [for v in data.huaweicloud_coc_resources.region_id_filter.data[*].region_id : v == huaweicloud_compute_instance.test.region]
  )
}

data "huaweicloud_coc_resources" "az_id_filter" {
  cloud_service_name = "ecs"
  type               = "cloudservers"
  az_id              = "cn-north-9a"
  depends_on         = [huaweicloud_compute_instance.test]
}

output "az_id_filter_is_useful" {
  value = length(data.huaweicloud_coc_resources.az_id_filter.data) > 0
}

data "huaweicloud_coc_resources" "ip_type_filter" {
  cloud_service_name = "ecs"
  type               = "cloudservers"
  ip_type            = "fixed"
  ip                 = huaweicloud_compute_instance.test.access_ip_v4
  ip_list            = [huaweicloud_compute_instance.test.access_ip_v4]
}

output "ip_filter_is_useful" {
  value = length(data.huaweicloud_coc_resources.ip_type_filter.data) > 0
}

data "huaweicloud_coc_resources" "resource_id_list_filter" {
  cloud_service_name = "ecs"
  type               = "cloudservers"
  resource_id_list   = [huaweicloud_compute_instance.test.id]
}

output "resource_id_list_filter_is_useful" {
  value = length(data.huaweicloud_coc_resources.resource_id_list_filter.data) > 0 && alltrue(
    [for v in data.huaweicloud_coc_resources.resource_id_list_filter.data[*].resource_id : v == huaweicloud_compute_instance.test.id]
  )
}

locals {
  project_id = [for v in data.huaweicloud_coc_resources.resource_id_list_filter.data[*].project_id : v if v != ""][0]
}

data "huaweicloud_coc_resources" "project_id_filter" {
  cloud_service_name = "ecs"
  type               = "cloudservers"
  project_id         = local.project_id
}

output "project_id_filter_is_useful" {
  value = length(data.huaweicloud_coc_resources.project_id_filter.data) > 0 && alltrue(
    [for v in data.huaweicloud_coc_resources.project_id_filter.data[*].project_id : v == local.project_id]
  )
}

data "huaweicloud_coc_resources" "status_filter" {
  cloud_service_name = "ecs"
  type               = "cloudservers"
  status             = huaweicloud_compute_instance.test.status
}

output "status_filter_is_useful" {
  value = length(data.huaweicloud_coc_resources.status_filter.data) > 0
}

data "huaweicloud_coc_resources" "agent_state_filter" {
  cloud_service_name = "ecs"
  type               = "cloudservers"
  agent_state        = "ONLINE"
  depends_on         = [huaweicloud_compute_instance.test]
}

output "agent_state_filter_is_useful" {
  value = length(data.huaweicloud_coc_resources.agent_state_filter.data) > 0 && alltrue(
    [for v in data.huaweicloud_coc_resources.agent_state_filter.data[*].agent_state : v == "ONLINE"]
  )
}

data "huaweicloud_coc_resources" "image_name_filter" {
  cloud_service_name = "ecs"
  type               = "cloudservers"
  image_name         = huaweicloud_compute_instance.test.image_name
}

output "image_name_filter_is_useful" {
  value = length(data.huaweicloud_coc_resources.image_name_filter.data) > 0
}

data "huaweicloud_coc_resources" "os_type_filter" {
  cloud_service_name = "ecs"
  type               = "cloudservers"
  os_type            = "Linux"
  depends_on         = [huaweicloud_compute_instance.test]
}

output "os_type_filter_is_useful" {
  value = length(data.huaweicloud_coc_resources.os_type_filter.data) > 0
}

locals {
  tag = [for v in data.huaweicloud_coc_resources.resource_id_list_filter.data[*].tags[0].value : v if v != ""][0]
}

data "huaweicloud_coc_resources" "tag_filter" {
  cloud_service_name = "ecs"
  type               = "cloudservers"
  tag                = local.tag
}

output "tag_filter_is_useful" {
  value = length(data.huaweicloud_coc_resources.tag_filter.data) > 0 && alltrue(
    [for v in data.huaweicloud_coc_resources.tag_filter.data[*].tags[0].value : v == local.tag]
  )
}

locals {
  tag_key = [for v in data.huaweicloud_coc_resources.resource_id_list_filter.data[*].tags[0].key : v if v != ""][0]
}

data "huaweicloud_coc_resources" "tag_key_filter" {
  cloud_service_name = "ecs"
  type               = "cloudservers"
  tag_key            = local.tag_key
}

output "tag_key_filter_is_useful" {
  value = length(data.huaweicloud_coc_resources.tag_key_filter.data) > 0 && alltrue(
    [for v in data.huaweicloud_coc_resources.tag_key_filter.data[*].tags[0].key : v == local.tag_key]
  )
}

data "huaweicloud_coc_resources" "vpc_id_filter" {
  cloud_service_name = "ecs"
  type               = "cloudservers"
  vpc_id             = data.huaweicloud_vpc_subnet.test.vpc_id
  depends_on         = [huaweicloud_compute_instance.test]
}

output "vpc_id_filter_is_useful" {
  value = length(data.huaweicloud_coc_resources.vpc_id_filter.data) > 0
}

data "huaweicloud_coc_resources" "flavor_name_filter" {
  cloud_service_name = "ecs"
  type               = "cloudservers"
  flavor_name        = huaweicloud_compute_instance.test.flavor_name
}

output "flavor_name_filter_is_useful" {
  value = length(data.huaweicloud_coc_resources.flavor_name_filter.data) > 0
}

data "huaweicloud_coc_resources" "charging_mode_filter" {
  cloud_service_name = "ecs"
  type               = "cloudservers"
  charging_mode      = "0"
  depends_on         = [huaweicloud_compute_instance.test]
}

output "charging_mode_filter_is_useful" {
  value = length(data.huaweicloud_coc_resources.charging_mode_filter.data) > 0
}
`, testAccComputeInstance_basic(name))
}

const testAccCompute_data = `
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_images_image" "test" {
  name        = "Ubuntu 18.04 server 64bit"
  most_recent = true
}

data "huaweicloud_networking_secgroup" "test" {
  name = "default"
}
`

func testAccComputeInstance_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_compute_instance" "test" {
  name               = "%s"
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [data.huaweicloud_networking_secgroup.test.id]

  network {
    uuid = data.huaweicloud_vpc_subnet.test.id
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testAccCompute_data, rName)
}
