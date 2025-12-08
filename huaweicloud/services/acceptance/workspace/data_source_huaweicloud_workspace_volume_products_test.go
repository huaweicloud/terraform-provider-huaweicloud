package workspace

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataVolumeProducts_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_workspace_volume_products.all"
		dc  = acceptance.InitDataSourceCheck(all)

		filterByZoneCode   = "data.huaweicloud_workspace_volume_products.filter_by_zone_code"
		dcFilterByZoneCode = acceptance.InitDataSourceCheck(filterByZoneCode)

		filterByVolumeType   = "data.huaweicloud_workspace_volume_products.filter_by_volume_type"
		dcFilterByVolumeType = acceptance.InitDataSourceCheck(filterByVolumeType)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataVolumeProducts_basic,
				Check: resource.ComposeTestCheckFunc(
					// Without any filter parameter.
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "volume_products.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "volume_products.0.resource_spec_code"),
					resource.TestCheckResourceAttrSet(all, "volume_products.0.volume_type"),
					resource.TestCheckResourceAttrSet(all, "volume_products.0.volume_product_type"),
					resource.TestCheckResourceAttrSet(all, "volume_products.0.resource_type"),
					resource.TestCheckResourceAttrSet(all, "volume_products.0.cloud_service_type"),
					resource.TestCheckResourceAttrSet(all, "volume_products.0.names.#"),
					resource.TestCheckResourceAttrSet(all, "volume_products.0.names.0.language"),
					resource.TestCheckResourceAttrSet(all, "volume_products.0.names.0.value"),
					resource.TestCheckResourceAttrSet(all, "volume_products.0.status"),
					// Filter by 'zone_code' parameter.
					dcFilterByZoneCode.CheckResourceExists(),
					resource.TestMatchResourceAttr(filterByZoneCode, "volume_products.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					// Filter by 'volume_type' parameter.
					dcFilterByVolumeType.CheckResourceExists(),
					resource.TestCheckOutput("is_volume_type_filter_useful", "true"),
				),
			},
		},
	})
}

const testAccDataVolumeProducts_basic = `
data "huaweicloud_workspace_volume_products" "all" {}

# Filter by zone code
data "huaweicloud_availability_zones" "test" {}

locals {
  zone_code = try(data.huaweicloud_availability_zones.test.names[0], "NOT_FOUND")
}

data "huaweicloud_workspace_volume_products" "filter_by_zone_code" {
  availability_zone = local.zone_code
}

locals {
  zone_code_filter_result = [
    for v in data.huaweicloud_workspace_volume_products.filter_by_zone_code.volume_products[*].availability_zone :
    v == local.zone_code
  ]
}

output "is_zone_code_filter_useful" {
  value = length(local.zone_code_filter_result) > 0 && alltrue(local.zone_code_filter_result)
}

# Filter by volume type
locals {
  volume_type = try(data.huaweicloud_workspace_volume_products.all.volume_products[0].volume_type, "NOT_FOUND")
}

data "huaweicloud_workspace_volume_products" "filter_by_volume_type" {
  volume_type = local.volume_type
}

locals {
  volume_type_filter_result = [
    for v in data.huaweicloud_workspace_volume_products.filter_by_volume_type.volume_products[*].volume_type :
    v == local.volume_type
  ]
}

output "is_volume_type_filter_useful" {
  value = length(local.volume_type_filter_result) > 0 && alltrue(local.volume_type_filter_result)
}
`
