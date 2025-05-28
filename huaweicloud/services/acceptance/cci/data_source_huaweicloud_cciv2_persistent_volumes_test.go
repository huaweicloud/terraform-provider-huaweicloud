package cci

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceV2PersistentVolumes_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_cciv2_persistent_volumes.test"
	rName := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceV2PersistentVolumes_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "persistent_volumes.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "persistent_volumes.0.annotations.%"),
					resource.TestCheckResourceAttrSet(dataSourceName, "persistent_volumes.0.labels.%"),
					resource.TestCheckResourceAttrSet(dataSourceName, "persistent_volumes.0.creation_timestamp"),
					resource.TestCheckResourceAttrSet(dataSourceName, "persistent_volumes.0.resource_version"),
					resource.TestCheckResourceAttrSet(dataSourceName, "persistent_volumes.0.uid"),
					resource.TestCheckResourceAttrSet(dataSourceName, "persistent_volumes.0.finalizers.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "persistent_volumes.0.status.#"),
				),
			},
		},
	})
}

func testAccDataSourceV2PersistentVolumes_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_cciv2_persistent_volumes" "test" {
  depends_on = [huaweicloud_cciv2_persistent_volume.test]
}
`, testAccV2PersistentVolume_basic(rName), rName)
}
