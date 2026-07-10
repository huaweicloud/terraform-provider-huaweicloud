package dsc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDscCatalogQuantityVariation_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dsc_catalog_quantity_variation.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDscTypeId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDscCatalogQuantityVariation_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "time_axis.#"),
					resource.TestCheckResourceAttrSet(dataSource, "sensitive_number_variation.#"),
					resource.TestCheckResourceAttrSet(dataSource, "total_number_variation.#"),
				),
			},
		},
	})
}

func testAccDataSourceDscCatalogQuantityVariation_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dsc_catalog_quantity_variation" "test" {
  type_id = "%[1]s"
}
`, acceptance.HW_DSC_TYPE_ID)
}
