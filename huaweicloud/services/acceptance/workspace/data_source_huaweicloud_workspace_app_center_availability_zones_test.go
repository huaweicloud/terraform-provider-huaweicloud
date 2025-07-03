package workspace

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccCenterAvailabilityZones_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_workspace_app_center_availability_zones.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCenterAvailabilityZones_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSourceName, "availability_zones.#", regexp.MustCompile(`^[0-9]+$`)),
					resource.TestCheckOutput("is_availability_zones_exist", "true"),
					resource.TestCheckOutput("is_availability_zone_set", "true"),
					resource.TestCheckOutput("is_display_name_set", "true"),
					resource.TestCheckOutput("is_default_az_set", "true"),
					resource.TestCheckOutput("is_visible_set", "true"),
					resource.TestCheckOutput("is_i18n_en_set", "true"),
					resource.TestCheckOutput("is_i18n_zh_set", "true"),
					resource.TestCheckOutput("is_product_ids_set", "true"),
					resource.TestCheckOutput("is_sold_out_set", "true"),
				),
			},
		},
	})
}

const testAccCenterAvailabilityZones_basic = `
data "huaweicloud_workspace_app_center_availability_zones" "test" {}

locals {
  azs = data.huaweicloud_workspace_app_center_availability_zones.test.availability_zones
  first_az = try(local.azs[0], null)
}

output "is_availability_zones_exist" {
  value = length(local.azs) >= 0
}

output "is_availability_zone_set" {
  value = length(local.azs) != 0 ? try(local.first_az.availability_zone != "", false) : true 
}

output "is_display_name_set" {
  value = length(local.azs) != 0 ? try(local.first_az.display_name != "", false) : true 
}

output "is_default_az_set" {
  value = length(local.azs) != 0 ? try(local.first_az.default_availability_zone != null, false) : true 
}

output "is_visible_set" {
  value = length(local.azs) != 0 ? try(local.first_az.visible != null, false) : true 
}

output "is_i18n_en_set" {
  value = length(local.azs) != 0 ? try(local.first_az.i18n.en_us != "", false) : true 
}

output "is_i18n_zh_set" {
  value = length(local.azs) != 0 ? try(local.first_az.i18n.zh_cn != "", false) : true 
}

output "is_product_ids_set" {
  value = length(local.azs) != 0 ? try(local.first_az.product_ids != null, false) : true 
}

output "is_sold_out_set" {
  value = length(local.azs) != 0 ? try(local.first_az.sold_out != null, false) : true
}
`
