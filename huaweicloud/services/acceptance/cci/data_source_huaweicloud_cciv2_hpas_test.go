package cci

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceV2HPAs_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_cciv2_hpas.test"
	rName := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceV2HPAs_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "hpas.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "hpas.0.name", rName),
					resource.TestCheckResourceAttrSet(dataSourceName, "hpas.0.namespace"),
					resource.TestCheckResourceAttrSet(dataSourceName, "hpas.0.scale_target_ref.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "hpas.0.metrics"),
					resource.TestCheckResourceAttrSet(dataSourceName, "hpas.0.creation_timestamp"),
					resource.TestCheckResourceAttrSet(dataSourceName, "hpas.0.resource_version"),
					resource.TestCheckResourceAttrSet(dataSourceName, "hpas.0.uid"),
					resource.TestCheckResourceAttrSet(dataSourceName, "hpas.0.status"),
				),
			},
		},
	})
}

func testAccDataSourceV2HPAs_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_cciv2_namespaces" "test" {
  depends_on = [huaweicloud_cciv2_hpa.test]

  namespace = huaweicloud_cciv2_namespace.test.name
}
`, testAccV2HPA_basic(rName))
}
