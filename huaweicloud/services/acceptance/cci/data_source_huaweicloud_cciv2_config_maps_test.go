package cci

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceV2ConfigMaps_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_cciv2_config_maps.test"
	rName := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceV2ConfigMaps_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "config_maps.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "config_maps.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "config_maps.0.namespace"),
					resource.TestCheckResourceAttrSet(dataSourceName, "config_maps.0.annotations.%"),
					resource.TestCheckResourceAttrSet(dataSourceName, "config_maps.0.labels.%"),
					resource.TestCheckResourceAttrSet(dataSourceName, "config_maps.0.creation_timestamp"),
					resource.TestCheckResourceAttrSet(dataSourceName, "config_maps.0.resource_version"),
					resource.TestCheckResourceAttrSet(dataSourceName, "config_maps.0.uid"),
					resource.TestCheckResourceAttrSet(dataSourceName, "config_maps.0.data.%"),
				),
			},
		},
	})
}

func testAccDataSourceV2ConfigMaps_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_cciv2_config_maps" "test" {
  depends_on = [huaweicloud_cciv2_config_map.test]

  namespace = huaweicloud_cciv2_namespace.test.name
}
`, testAccV2ConfigMap_basic(rName))
}
