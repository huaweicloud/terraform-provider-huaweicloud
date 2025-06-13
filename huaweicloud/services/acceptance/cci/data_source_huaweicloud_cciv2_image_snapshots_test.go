package cci

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceV2ImageSnapshots_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_cciv2_image_snapshots.test"
	rName := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceV2ImageSnapshots_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "image_snapshots.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "image_snapshots.0.annotations.%"),
					resource.TestCheckResourceAttrSet(dataSourceName, "image_snapshots.0.labels.%"),
					resource.TestCheckResourceAttrSet(dataSourceName, "image_snapshots.0.image_snapshots_size"),
					resource.TestCheckResourceAttrSet(dataSourceName, "image_snapshots.0.resource_version"),
					resource.TestCheckResourceAttrSet(dataSourceName, "image_snapshots.0.uid"),
					resource.TestCheckResourceAttrSet(dataSourceName, "image_snapshots.0.images.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "image_snapshots.0.status.#"),
				),
			},
		},
	})
}

func testAccDataSourceV2ImageSnapshots_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cciv2_image_snapshot" "test" {
  name = "%[2]s"

  building_config {
    namespace = huaweicloud_cciv2_namespace.test.name
  }

  image_snapshots_size   = 20
  ttl_days_after_created = 100

  images {
    image = "nginx:latest"
  }
}

data "huaweicloud_cciv2_image_snapshots" test {
  name = "%[2]s"
}
`, testAccV2Namespace_basic(rName), rName)
}
