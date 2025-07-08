package workspace

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccAppAvailableVolumes_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_workspace_app_available_volumes.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAppAvailableVolumes_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSourceName, "volume_types.#", regexp.MustCompile(`^[0-9]+$`)),
					resource.TestCheckOutput("is_volume_types_exist", "true"),
					resource.TestCheckOutput("is_resource_spec_code_set", "true"),
					resource.TestCheckOutput("is_volume_type_set", "true"),
					resource.TestCheckOutput("is_volume_product_type_set", "true"),
					resource.TestCheckOutput("is_resource_type_set", "true"),
					resource.TestCheckOutput("is_cloud_service_type_set", "true"),
					resource.TestCheckOutput("is_name_set", "true"),
					resource.TestCheckOutput("is_volume_type_extra_specs_set", "true"),
					resource.TestCheckOutput("is_availability_zone_set", "true"),
					resource.TestCheckOutput("is_sold_out_availability_zone_set", "true"),
				),
			},
		},
	})
}

const testAccAppAvailableVolumes_basic = `
data "huaweicloud_workspace_app_available_volumes" "test" {}

locals {
  volume_types = data.huaweicloud_workspace_app_available_volumes.test.volume_types
  first_volume = try(local.volume_types[0], {})
  first_specs  = try(local.first_volume.volume_type_extra_specs[0], {})
}

output "is_volume_types_exist" {
  value = length(local.volume_types) >= 0
}

output "is_resource_spec_code_set" {
  value = length(local.volume_types) != 0 ? try(local.first_volume.resource_spec_code != "", false) : true
}

output "is_volume_type_set" {
  value = length(local.volume_types) != 0 ? try(local.first_volume.volume_type != "", false) : true
}

output "is_volume_product_type_set" {
  value = length(local.volume_types) != 0 ? try(local.first_volume.volume_product_type != "", false) : true
}

output "is_resource_type_set" {
  value = length(local.volume_types) != 0 ? try(local.first_volume.resource_type != "", false) : true
}

output "is_cloud_service_type_set" {
  value = length(local.volume_types) != 0 ? try(local.first_volume.cloud_service_type != "", false) : true
}

output "is_name_set" {
  value = length(local.volume_types) != 0 ? try(local.first_volume.name != null, false) : true
}

output "is_volume_type_extra_specs_set" {
  value = length(local.volume_types) != 0 ? try(local.first_volume.volume_type_extra_specs != null, false) : true
}

output "is_availability_zone_set" {
  value = length(local.volume_types) != 0 ? try(local.first_specs.availability_zone != "", false) : true
}

output "is_sold_out_availability_zone_set" {
  value = length(local.volume_types) != 0 ? try(local.first_specs.sold_out_availability_zone != null, false) : true
}
`
