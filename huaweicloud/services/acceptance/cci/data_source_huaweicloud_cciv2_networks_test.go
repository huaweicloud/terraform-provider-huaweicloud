package cci

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceV2Networks_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_cciv2_networks.test"
	rName := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceV2Networks_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "networks.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "networks.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "networks.0.namespace"),
					resource.TestCheckResourceAttrSet(dataSourceName, "networks.0.annotations.%"),
					resource.TestCheckResourceAttrSet(dataSourceName, "networks.0.creation_timestamp"),
					resource.TestCheckResourceAttrSet(dataSourceName, "networks.0.resource_version"),
					resource.TestCheckResourceAttrSet(dataSourceName, "networks.0.uid"),
					resource.TestCheckResourceAttrSet(dataSourceName, "networks.0.security_group_ids.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "networks.0.subnets.#"),
				),
			},
		},
	})
}

func testAccDataSourceV2Networks_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_cciv2_networks" "test" {
  depends_on = [huaweicloud_cciv2_network.test]

  namespace = huaweicloud_cciv2_namespace.test.name
}
`, testAccV2Network_basic(rName))
}
