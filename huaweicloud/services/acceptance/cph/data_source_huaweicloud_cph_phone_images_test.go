package cph

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourcePhoneImages_basic(t *testing.T) {
	rName := "data.huaweicloud_cph_phone_images.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourcePhoneImages_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "images.0.id"),
					resource.TestCheckResourceAttrSet(rName, "images.0.name"),
					resource.TestCheckResourceAttrSet(rName, "images.0.os_type"),
					resource.TestCheckResourceAttrSet(rName, "images.0.os_name"),
					resource.TestCheckResourceAttrSet(rName, "images.0.is_public"),
					resource.TestCheckResourceAttrSet(rName, "images.0.image_label"),

					resource.TestCheckOutput("is_public_filter_useful", "true"),

					resource.TestCheckOutput("is_image_label_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourcePhoneImages_basic() string {
	return `
data "huaweicloud_cph_phone_images" "test" {
}

data "huaweicloud_cph_phone_images" "is_public_filter" {
  is_public = 1
}

output "is_public_filter_useful" {
  value = alltrue([for v in data.huaweicloud_cph_phone_images.is_public_filter.images[*].is_public : v == 1])
}

data "huaweicloud_cph_phone_images" "image_label_filter" {
  image_label = "cloud_phone"
}

output "is_image_label_filter_useful" {
  value = alltrue([for v in data.huaweicloud_cph_phone_images.image_label_filter.images[*].image_label : v == "cloud_phone"])
}
`
}
