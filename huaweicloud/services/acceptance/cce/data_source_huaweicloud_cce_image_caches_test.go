package cce

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceHuaweiCloudCceImageCaches_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cce_image_caches.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testImageCaches_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("is_results_not_empty", "true"),
				),
			},
		},
	})
}

const testImageCaches_basic = `
data "huaweicloud_cce_image_caches" "test" {}

output "is_results_not_empty" {
  value = length(data.huaweicloud_cce_image_caches.test.image_caches) > 0
}
`
