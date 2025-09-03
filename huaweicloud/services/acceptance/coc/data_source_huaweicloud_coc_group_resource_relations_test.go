package coc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCocGroupResourceRelations_basic(t *testing.T) {
	dataSource := "data.huaweicloud_coc_group_resource_relations.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: "0.12.1",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCocGroupResourceRelations_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.cmdb_resource_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.group_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.group_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.resource_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.cloud_service_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.region_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.ep_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.domain_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.tags.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.tags.0.key"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.tags.0.value"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.properties"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.create_time"),
					resource.TestCheckOutput("application_id_filter_is_useful", "true"),
					resource.TestCheckOutput("component_id_filter_is_useful", "true"),
					resource.TestCheckOutput("group_id_filter_is_useful", "true"),
					resource.TestCheckOutput("resource_id_list_filter_is_useful", "true"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("region_id_filter_is_useful", "true"),
					resource.TestCheckOutput("az_id_filter_is_useful", "true"),
					resource.TestCheckOutput("ip_filter_is_useful", "true"),
					resource.TestCheckOutput("status_filter_is_useful", "true"),
					resource.TestCheckOutput("image_name_filter_is_useful", "true"),
					resource.TestCheckOutput("os_type_filter_is_useful", "true"),
					resource.TestCheckOutput("charging_mode_filter_is_useful", "true"),
					resource.TestCheckOutput("flavor_name_filter_is_useful", "true"),
					resource.TestCheckOutput("is_collected_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceCocGroupResourceRelations_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_coc_group_resource_relations" "test" {
  cloud_service_name = "ecs"
  type               = "cloudservers"
  vendor             = "RMS"
  application_id     = huaweicloud_coc_application.test.id
  depends_on         = [huaweicloud_coc_group_resource_relation.test]
}

output "application_id_filter_is_useful" {
  value = length(data.huaweicloud_coc_group_resource_relations.test.data) > 0
}

data "huaweicloud_coc_group_resource_relations" "component_id_filter" {
  cloud_service_name = "ecs"
  type               = "cloudservers"
  vendor             = "RMS"
  component_id       = huaweicloud_coc_component.test.id
  depends_on         = [huaweicloud_coc_group_resource_relation.test]
}

output "component_id_filter_is_useful" {
  value = length(data.huaweicloud_coc_group_resource_relations.component_id_filter.data) > 0
}

data "huaweicloud_coc_group_resource_relations" "group_id_filter" {
  cloud_service_name = "ecs"
  type               = "cloudservers"
  vendor             = "RMS"
  group_id           = huaweicloud_coc_group.test.id
  depends_on         = [huaweicloud_coc_group_resource_relation.test]
}

output "group_id_filter_is_useful" {
  value = length(data.huaweicloud_coc_group_resource_relations.group_id_filter.data) > 0 && alltrue(
   [for v in data.huaweicloud_coc_group_resource_relations.group_id_filter.data[*].group_id :
     v == huaweicloud_coc_group.test.id]
 )
}

data "huaweicloud_coc_group_resource_relations" "resource_id_list_filter" {
  cloud_service_name = "ecs"
  type               = "cloudservers"
  vendor             = "RMS"
  application_id     = huaweicloud_coc_application.test.id
  resource_id_list   = [huaweicloud_compute_instance.test.id]
  depends_on         = [huaweicloud_coc_group_resource_relation.test]
}

output "resource_id_list_filter_is_useful" {
  value = length(data.huaweicloud_coc_group_resource_relations.resource_id_list_filter.data) > 0 && alltrue(
    [for v in data.huaweicloud_coc_group_resource_relations.resource_id_list_filter.data[*].resource_id :
      v == huaweicloud_compute_instance.test.id]
  )
}

data "huaweicloud_coc_group_resource_relations" "name_filter" {
  cloud_service_name = "ecs"
  type               = "cloudservers"
  vendor             = "RMS"
  application_id     = huaweicloud_coc_application.test.id
  name               = huaweicloud_compute_instance.test.name
  depends_on         = [huaweicloud_coc_group_resource_relation.test]
}

output "name_filter_is_useful" {
  value = length(data.huaweicloud_coc_group_resource_relations.name_filter.data) > 0 && alltrue(
    [for v in data.huaweicloud_coc_group_resource_relations.name_filter.data[*].name :
      v == huaweicloud_compute_instance.test.name]
  )
}

data "huaweicloud_coc_group_resource_relations" "region_id_filter" {
  cloud_service_name = "ecs"
  type               = "cloudservers"
  vendor             = "RMS"
  application_id     = huaweicloud_coc_application.test.id
  region_id          = huaweicloud_compute_instance.test.region
  depends_on         = [huaweicloud_coc_group_resource_relation.test]
}

output "region_id_filter_is_useful" {
  value = length(data.huaweicloud_coc_group_resource_relations.region_id_filter.data) > 0 && alltrue(
    [for v in data.huaweicloud_coc_group_resource_relations.region_id_filter.data[*].region_id :
      v == huaweicloud_compute_instance.test.region]
  )
}

data "huaweicloud_coc_group_resource_relations" "az_id_filter" {
  cloud_service_name = "ecs"
  type               = "cloudservers"
  vendor             = "RMS"
  application_id     = huaweicloud_coc_application.test.id
  az_id              = huaweicloud_compute_instance.test.availability_zone
  depends_on         = [huaweicloud_coc_group_resource_relation.test]
}

output "az_id_filter_is_useful" {
  value = length(data.huaweicloud_coc_group_resource_relations.az_id_filter.data) > 0
}

data "huaweicloud_coc_group_resource_relations" "ip_type_filter" {
  cloud_service_name = "ecs"
  type               = "cloudservers"
  vendor             = "RMS"
  application_id     = huaweicloud_coc_application.test.id
  ip_type            = "fixed"
  ip                 = huaweicloud_compute_instance.test.access_ip_v4
  ip_list            = [huaweicloud_compute_instance.test.access_ip_v4]
  depends_on         = [huaweicloud_coc_group_resource_relation.test]
}

output "ip_filter_is_useful" {
  value = length(data.huaweicloud_coc_group_resource_relations.ip_type_filter.data) > 0
}

data "huaweicloud_coc_group_resource_relations" "status_filter" {
  cloud_service_name = "ecs"
  type               = "cloudservers"
  vendor             = "RMS"
  application_id     = huaweicloud_coc_application.test.id
  status             = huaweicloud_compute_instance.test.status
  depends_on         = [huaweicloud_coc_group_resource_relation.test]
}

output "status_filter_is_useful" {
  value = length(data.huaweicloud_coc_group_resource_relations.status_filter.data) > 0
}

data "huaweicloud_coc_group_resource_relations" "image_name_filter" {
  cloud_service_name = "ecs"
  type               = "cloudservers"
  vendor             = "RMS"
  application_id     = huaweicloud_coc_application.test.id
  image_name         = huaweicloud_compute_instance.test.image_name
  depends_on         = [huaweicloud_coc_group_resource_relation.test]
}

output "image_name_filter_is_useful" {
  value = length(data.huaweicloud_coc_group_resource_relations.image_name_filter.data) > 0
}

data "huaweicloud_coc_group_resource_relations" "os_type_filter" {
  cloud_service_name = "ecs"
  type               = "cloudservers"
  vendor             = "RMS"
  application_id     = huaweicloud_coc_application.test.id
  os_type            = "Linux"
  depends_on         = [huaweicloud_coc_group_resource_relation.test]
}

output "os_type_filter_is_useful" {
  value = length(data.huaweicloud_coc_group_resource_relations.os_type_filter.data) > 0
}

data "huaweicloud_coc_group_resource_relations" "charging_mode_filter" {
  cloud_service_name = "ecs"
  type               = "cloudservers"
  vendor             = "RMS"
  application_id     = huaweicloud_coc_application.test.id
  charging_mode      = "0"
  depends_on         = [huaweicloud_coc_group_resource_relation.test]
}

output "charging_mode_filter_is_useful" {
  value = length(data.huaweicloud_coc_group_resource_relations.charging_mode_filter.data) > 0
}

data "huaweicloud_coc_group_resource_relations" "flavor_name_filter" {
  cloud_service_name = "ecs"
  type               = "cloudservers"
  vendor             = "RMS"
  application_id     = huaweicloud_coc_application.test.id
  flavor_name        = huaweicloud_compute_instance.test.flavor_name
  depends_on         = [huaweicloud_coc_group_resource_relation.test]
}

output "flavor_name_filter_is_useful" {
  value = length(data.huaweicloud_coc_group_resource_relations.flavor_name_filter.data) > 0
}

data "huaweicloud_coc_group_resource_relations" "is_collected_filter" {
  cloud_service_name = "ecs"
  type               = "cloudservers"
  vendor             = "RMS"
  application_id     = huaweicloud_coc_application.test.id
  is_collected       = false
  depends_on         = [huaweicloud_coc_group_resource_relation.test]
}

output "is_collected_filter_is_useful" {
  value = length(data.huaweicloud_coc_group_resource_relations.is_collected_filter.data) > 0
}
`, testAccGroupResourceRelation_basic(name))
}
