package ecs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceEcsComputeResizeFlavors_basic(t *testing.T) {
	dataSource := "data.huaweicloud_compute_resize_flavors.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceEcsComputeResizeFlavors_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.#"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.vcpus"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.ram"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.disk"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.os_flavor_access_is_public"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.links.#"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.links.0.rel"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.links.0.href"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.extra_specs.#"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.extra_specs.0.resource_type"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.extra_specs.0.cond_operation_az"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.extra_specs.0.quota_max_rate"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.extra_specs.0.info_cpu_name"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.extra_specs.0.quota_sub_network_interface_max_num"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.extra_specs.0.cond_network"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.extra_specs.0.hw_numa_nodes"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.extra_specs.0.quota_vif_max_num"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.extra_specs.0.ecs_performancetype"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.extra_specs.0.quota_min_rate"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.extra_specs.0.ecs_virtualization_env_types"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.extra_specs.0.quota_max_pps"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.extra_specs.0.cond_compute"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.extra_specs.0.ecs_generation"),
					resource.TestCheckOutput("instance_uuid_filter_is_useful", "true"),
					resource.TestCheckOutput("source_flavor_id_filter_is_useful", "true"),
					resource.TestCheckOutput("source_flavor_name_filter_is_useful", "true"),
					resource.TestCheckOutput("sort_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceEcsComputeResizeFlavors_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  vpc_id     = huaweicloud_vpc.test.id
  name       = "%[1]s"
  cidr       = cidrsubnet(huaweicloud_vpc.test.cidr, 4, 0)
  gateway_ip = cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 4, 0), 1)
}

resource "huaweicloud_networking_secgroup" "test" {
  name = "%[1]s"
}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

data "huaweicloud_images_images" "test" {
  flavor_id = data.huaweicloud_compute_flavors.test.ids[0]

  os         = "Windows"
  visibility = "public"
}

resource "huaweicloud_compute_instance" "test" {
  name               = "%[1]s"
  image_id           = data.huaweicloud_images_images.test.images[0].id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [huaweicloud_networking_secgroup.test.id]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  system_disk_type   = "SSD"

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }
}
`, name)
}

func testDataSourceEcsComputeResizeFlavors_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_compute_resize_flavors" "test" {
  source_flavor_id = data.huaweicloud_compute_flavors.test.flavors[0].id
}

data "huaweicloud_compute_resize_flavors" "instance_uuid_filter" {
  instance_uuid = huaweicloud_compute_instance.test.id
}
output "instance_uuid_filter_is_useful" {
  value = length(data.huaweicloud_compute_resize_flavors.instance_uuid_filter.flavors) > 0
}

data "huaweicloud_compute_resize_flavors" "source_flavor_id_filter" {
  source_flavor_id = data.huaweicloud_compute_flavors.test.flavors[0].id
}
output "source_flavor_id_filter_is_useful" {
  value = length(data.huaweicloud_compute_resize_flavors.source_flavor_id_filter.flavors) > 0
}

data "huaweicloud_compute_resize_flavors" "source_flavor_name_filter" {
  source_flavor_name = data.huaweicloud_compute_resize_flavors.test.flavors[0].name
}
output "source_flavor_name_filter_is_useful" {
  value = length(data.huaweicloud_compute_resize_flavors.source_flavor_name_filter.flavors) > 0
}

data "huaweicloud_compute_resize_flavors" "sort_filter" {
  source_flavor_id = data.huaweicloud_compute_flavors.test.flavors[0].id
  sort_dir         = "desc"
  sort_key         = "vcpus"
}
output "sort_filter_is_useful" {
  value = length(data.huaweicloud_compute_resize_flavors.sort_filter.flavors) > 0
}
`, testDataSourceEcsComputeResizeFlavors_base(name))
}
