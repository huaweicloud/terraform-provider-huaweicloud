package cce

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCceImageCaches_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cce_image_caches.test"
	rName := acceptance.RandomAccResourceNameWithDash()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCceClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCceImageCaches_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "image_caches.#"),
					resource.TestCheckResourceAttrSet(dataSource, "image_caches.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "image_caches.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "image_caches.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "image_caches.0.images.#"),
					resource.TestCheckResourceAttrSet(dataSource, "image_caches.0.image_cache_size"),
					resource.TestCheckResourceAttrSet(dataSource, "image_caches.0.retention_days"),
					resource.TestCheckResourceAttrSet(dataSource, "image_caches.0.building_config.#"),
					resource.TestCheckResourceAttrSet(dataSource, "image_caches.0.building_config.0.cluster"),
					resource.TestCheckResourceAttrSet(dataSource, "image_caches.0.building_config.0.image_pull_secrets.#"),
					resource.TestCheckResourceAttrSet(dataSource, "image_caches.0.status"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceCceImageCaches_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cce_image_cache" "test" {
  name   = "%[1]s"
  images = ["swr.%[2]s.myhuaweicloud.com/xpanse/kafka:latest"]

  building_config {
    cluster            = "%[3]s"
    image_pull_secrets = ["default:default-secret"]
  }

  image_cache_size = 20
  retention_days   = 7
}
`, name, acceptance.HW_REGION_NAME, acceptance.HW_CCE_CLUSTER_ID)
}

func testAccDataSourceCceImageCaches_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_cce_image_caches" "test" {
  depends_on = [huaweicloud_cce_image_cache.test]
}

data "huaweicloud_cce_image_caches" "name_filter" {
  depends_on = [huaweicloud_cce_image_cache.test]

  name = huaweicloud_cce_image_cache.test.name
}
locals {
  name = huaweicloud_cce_image_cache.test.name
}

output "name_filter_is_useful" {
  value = length(data.huaweicloud_cce_image_caches.name_filter.image_caches) > 0 && alltrue(
    [for v in data.huaweicloud_cce_image_caches.name_filter.image_caches[*].name : v == local.name]
  )
}
`, testAccDataSourceCceImageCaches_base(name))
}
