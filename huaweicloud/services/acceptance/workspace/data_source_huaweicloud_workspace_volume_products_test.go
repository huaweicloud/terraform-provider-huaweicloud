package workspace

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccVolumeProductsDataSource_basic(t *testing.T) {
	var (
		dcName = "data.huaweicloud_workspace_volume_products.all"
		dc     = acceptance.InitDataSourceCheck(dcName)

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
				Config: testAccVolumeProductsDataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					// Query volume products without any filter parameter
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dcName, "volume_products.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dcName, "volume_products.0.resource_spec_code"),
					resource.TestCheckResourceAttrSet(dcName, "volume_products.0.volume_type"),
					resource.TestCheckResourceAttrSet(dcName, "volume_products.0.volume_product_type"),
					resource.TestCheckResourceAttrSet(dcName, "volume_products.0.resource_type"),
					resource.TestCheckResourceAttrSet(dcName, "volume_products.0.cloud_service_type"),
					resource.TestCheckResourceAttrSet(dcName, "volume_products.0.names.#"),
					resource.TestCheckResourceAttrSet(dcName, "volume_products.0.names.0.language"),
					resource.TestCheckResourceAttrSet(dcName, "volume_products.0.names.0.value"),
					resource.TestCheckResourceAttrSet(dcName, "volume_products.0.status"),
					// Filter By zone code
					dcFilterByZoneCode.CheckResourceExists(),
					resource.TestMatchResourceAttr(filterByZoneCode, "volume_products.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					// Filter By volume type
					dcFilterByVolumeType.CheckResourceExists(),
					resource.TestCheckOutput("is_volume_type_filter_useful", "true"),
				),
			},
		},
	})
}

const testAccVolumeProductsDataSource_basic = `
data "huaweicloud_workspace_volume_products" "all" {}

locals {
  volume_type = try(data.huaweicloud_workspace_volume_products.all.volume_products[0].volume_type, "NOT_FOUND")
}

# Filter by zone code
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_workspace_volume_products" "filter_by_zone_code" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
}

# Filter by volume type
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
