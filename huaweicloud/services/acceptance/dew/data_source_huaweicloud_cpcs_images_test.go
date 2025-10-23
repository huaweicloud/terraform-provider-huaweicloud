package dew

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCpcsImages_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cpcs_images.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCpcsImages_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "images.0.image_id"),
					resource.TestCheckResourceAttrSet(dataSource, "images.0.image_name"),
					resource.TestCheckResourceAttrSet(dataSource, "images.0.service_type"),
					resource.TestCheckResourceAttrSet(dataSource, "images.0.arch_type"),
					resource.TestCheckResourceAttrSet(dataSource, "images.0.specification_id"),
					resource.TestCheckResourceAttrSet(dataSource, "images.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "images.0.version_type"),
					resource.TestCheckResourceAttrSet(dataSource, "images.0.vendor_name"),

					resource.TestCheckOutput("is_image_name_filter_useful", "true"),
					resource.TestCheckOutput("is_service_type_filter_useful", "true"),
				),
			},
		},
	})
}

const testDataSourceCpcsImages_basic = `
data "huaweicloud_cpcs_images" "test" {}

locals {
  image_name = data.huaweicloud_cpcs_images.test.images.0.image_name
  service_type = data.huaweicloud_cpcs_images.test.images.0.service_type
}

# Filter by image_name
data "huaweicloud_cpcs_images" "image_name_filter" {
  image_name = local.image_name
}

locals {
  image_name_filter_result = [
    for v in data.huaweicloud_cpcs_images.image_name_filter.images[*].image_name : v == local.image_name
  ]
}

output "is_image_name_filter_useful" {
  value = length(local.image_name_filter_result) > 0 && alltrue(local.image_name_filter_result)
}

# Filter by service_type
data "huaweicloud_cpcs_images" "service_type_filter" {
  service_type = local.service_type
}

locals {
  service_type_filter_result = [
    for v in data.huaweicloud_cpcs_images.service_type_filter.images[*].service_type : v == local.service_type
  ]
}

output "is_service_type_filter_useful" {
  value = length(local.service_type_filter_result) > 0 && alltrue(local.service_type_filter_result)
}
`
