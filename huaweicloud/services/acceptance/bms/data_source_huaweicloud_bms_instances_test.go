package bms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccBmsInstancesDataSource_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	resourceName := "data.huaweicloud_bms_instances.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckUserId(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccBmsInstancesDataSource_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBmsFlavorDataSourceID(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "servers.#"),
					resource.TestCheckResourceAttrSet(resourceName, "servers.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "servers.0.name"),
					resource.TestCheckResourceAttrSet(resourceName, "servers.0.nics.#"),
					resource.TestCheckResourceAttrSet(resourceName, "servers.0.flavor_id"),
					resource.TestCheckResourceAttrSet(resourceName, "servers.0.flavor_name"),
					resource.TestCheckResourceAttrSet(resourceName, "servers.0.disk"),
					resource.TestCheckResourceAttrSet(resourceName, "servers.0.vcpus"),
					resource.TestCheckResourceAttrSet(resourceName, "servers.0.memory"),
					resource.TestCheckResourceAttrSet(resourceName, "servers.0.status"),
					resource.TestCheckResourceAttrSet(resourceName, "servers.0.created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "servers.0.updated_at"),
					resource.TestCheckResourceAttrSet(resourceName, "servers.0.vpc_id"),
					resource.TestCheckResourceAttrSet(resourceName, "servers.0.agency_name"),
					resource.TestCheckResourceAttrSet(resourceName, "servers.0.image_name"),
					resource.TestCheckResourceAttrSet(resourceName, "servers.0.image_type"),
					resource.TestCheckResourceAttrSet(resourceName, "servers.0.description"),
					resource.TestCheckResourceAttrSet(resourceName, "servers.0.locked"),
					resource.TestCheckResourceAttrSet(resourceName, "servers.0.user_id"),
					resource.TestCheckResourceAttrSet(resourceName, "servers.0.key_pair"),
					resource.TestCheckResourceAttrSet(resourceName, "servers.0.volumes_attached.#"),
					resource.TestCheckResourceAttrSet(resourceName, "servers.0.vm_state"),
					resource.TestCheckResourceAttrSet(resourceName, "servers.0.disk_config"),
					resource.TestCheckResourceAttrSet(resourceName, "servers.0.availability_zone"),
					resource.TestCheckResourceAttrSet(resourceName, "servers.0.root_device_name"),
					resource.TestCheckResourceAttrSet(resourceName, "servers.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "servers.0.user_data"),
					resource.TestCheckResourceAttrSet(resourceName, "servers.0.security_groups.#"),
					resource.TestCheckResourceAttrSet(resourceName, "servers.0.image_id"),

					resource.TestCheckOutput("flavor_id_filter_is_useful", "true"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("status_filter_is_useful", "true"),
					resource.TestCheckOutput("enterprise_project_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccBmsInstancesDataSource_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_bms_instances" "test" {
  depends_on = [huaweicloud_bms_instance.test]
}

data "huaweicloud_bms_instances" "flavor_id_filter" {
  flavor_id = data.huaweicloud_bms_instances.test.servers.0.flavor_id
}
locals {
  flavor_id = data.huaweicloud_bms_instances.test.servers.0.flavor_id
}
output "flavor_id_filter_is_useful" {
  value = length(data.huaweicloud_bms_instances.flavor_id_filter.servers) > 0 && alltrue(
    [for v in data.huaweicloud_bms_instances.flavor_id_filter.servers[*].flavor_id : v == local.flavor_id]
  )  
}

data "huaweicloud_bms_instances" "name_filter" {
  name = huaweicloud_bms_instance.test.name
}
locals {
  name = data.huaweicloud_bms_instances.test.servers.0.name
}
output "name_filter_is_useful" {
  value = length(data.huaweicloud_bms_instances.name_filter.servers) > 0 && alltrue(
    [for v in data.huaweicloud_bms_instances.name_filter.servers[*].name : v == local.name]
  )  
}

data "huaweicloud_bms_instances" "status_filter" {
  status = data.huaweicloud_bms_instances.test.servers.0.status
}
locals {
  status = data.huaweicloud_bms_instances.test.servers.0.status
}
output "status_filter_is_useful" {
  value = length(data.huaweicloud_bms_instances.status_filter.servers) > 0 && alltrue(
    [for v in data.huaweicloud_bms_instances.status_filter.servers[*].status : v == local.status]
  )  
}

data "huaweicloud_bms_instances" "enterprise_project_id_filter" {
  enterprise_project_id = data.huaweicloud_bms_instances.test.servers.0.enterprise_project_id
} 
locals {
  enterprise_project_id = data.huaweicloud_bms_instances.test.servers.0.enterprise_project_id
}
output "enterprise_project_id_filter_is_useful" {
  value = length(data.huaweicloud_bms_instances.enterprise_project_id_filter.servers) > 0 && alltrue(
    [for v in data.huaweicloud_bms_instances.enterprise_project_id_filter.servers[*].enterprise_project_id : v == local.enterprise_project_id]
  )  
}
`, testAccBmsInstance_basic(name))
}
