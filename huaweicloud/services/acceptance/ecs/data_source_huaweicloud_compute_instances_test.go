package ecs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccComputeInstancesDataSource_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceNameWithDash()
	dataSource1 := "data.huaweicloud_compute_instances.basic"
	dataSource2 := "data.huaweicloud_compute_instances.filter_by_name"
	dataSource3 := "data.huaweicloud_compute_instances.filter_by_id"
	dataSource4 := "data.huaweicloud_compute_instances.filter_by_ip"
	dataSource5 := "data.huaweicloud_compute_instances.filter_by_tags"

	dc1 := acceptance.InitDataSourceCheck(dataSource1)
	dc2 := acceptance.InitDataSourceCheck(dataSource2)
	dc3 := acceptance.InitDataSourceCheck(dataSource3)
	dc4 := acceptance.InitDataSourceCheck(dataSource4)
	dc5 := acceptance.InitDataSourceCheck(dataSource5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstancesDataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc1.CheckResourceExists(),
					dc2.CheckResourceExists(),
					dc3.CheckResourceExists(),
					dc4.CheckResourceExists(),
					dc5.CheckResourceExists(),
					resource.TestCheckOutput("is_results_not_empty", "true"),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("is_id_filter_useful", "true"),
					resource.TestCheckOutput("is_ip_filter_useful", "true"),
					resource.TestCheckOutput("is_tags_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccComputeInstancesDataSource_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_compute_instance" "test" {
  name               = "%[2]s"
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [data.huaweicloud_networking_secgroup.test.id]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]

  network {
    uuid = data.huaweicloud_vpc_subnet.test.id
  }

  tags = {
    server_name = "%[2]s"
  }
}

data "huaweicloud_compute_instances" "basic" {
  depends_on = [huaweicloud_compute_instance.test]
}

data "huaweicloud_compute_instances" "filter_by_name" {
  name = huaweicloud_compute_instance.test.name
}

data "huaweicloud_compute_instances" "filter_by_id" {
  instance_id = huaweicloud_compute_instance.test.id
}

data "huaweicloud_compute_instances" "filter_by_ip" {
  fixed_ip_v4 = huaweicloud_compute_instance.test.network[0].fixed_ip_v4
}

data "huaweicloud_compute_instances" "filter_by_tags" {
  tags = {
    server_name = "%[2]s"
  }

  depends_on = [huaweicloud_compute_instance.test]
}

locals {
  name_filter_result = [for v in data.huaweicloud_compute_instances.filter_by_name.instances[*].name : v == "%[2]s"]
  id_filter_result = [
    for v in data.huaweicloud_compute_instances.filter_by_id.instances[*].id : v == huaweicloud_compute_instance.test.id
  ]
  ip_filter_result = [
    for v in data.huaweicloud_compute_instances.filter_by_id.instances[*].network[0].fixed_ip_v4 :
    v == huaweicloud_compute_instance.test.network[0].fixed_ip_v4
  ]
  tags_filter_result = [for v in data.huaweicloud_compute_instances.filter_by_tags.instances[*].tags.server_name : v == "%[2]s"]
}

output "is_results_not_empty" {
  value = length(data.huaweicloud_compute_instances.basic.instances) > 0
}

output "is_name_filter_useful" {
  value = alltrue(local.name_filter_result) && length(local.name_filter_result) > 0
}

output "is_id_filter_useful" {
  value = alltrue(local.id_filter_result) && length(local.id_filter_result) > 0
}

output "is_ip_filter_useful" {
  value = alltrue(local.ip_filter_result) && length(local.ip_filter_result) > 0
}

output "is_tags_filter_useful" {
  value = alltrue(local.tags_filter_result) && length(local.tags_filter_result) > 0
}
`, testAccCompute_data, rName)
}
