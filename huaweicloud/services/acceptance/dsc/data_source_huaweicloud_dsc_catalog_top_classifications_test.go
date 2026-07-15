package dsc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Before running this test, ensure that sensitive data identification has been completed for specific assets
// under the HW_DSC_TYPE_ID.
func TestAccDataSourceCatalogTopClassifications_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_dsc_catalog_top_classifications.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDscTypeId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCatalogTopClassifications_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "classifications.#", regexp.MustCompile(`^[1-9]([0-9]+)?$`)),
					resource.TestCheckResourceAttrSet(dataSource, "classifications.0.classification_name"),
					resource.TestMatchResourceAttr(dataSource, "classifications.0.hit_number", regexp.MustCompile(`^[1-9]([0-9]+)?$`)),
					resource.TestMatchResourceAttr(dataSource, "classifications.0.column_details.#", regexp.MustCompile(`^[1-9]([0-9]+)?$`)),
					resource.TestCheckResourceAttrSet(dataSource, "classifications.0.column_details.0.asset_id"),
					resource.TestCheckResourceAttrSet(dataSource, "classifications.0.column_details.0.asset_name"),
					resource.TestCheckResourceAttrSet(dataSource, "classifications.0.column_details.0.column_fqn"),
					resource.TestCheckResourceAttrSet(dataSource, "classifications.0.column_details.0.db_type"),
				),
			},
		},
	})
}

func testAccDataSourceCatalogTopClassifications_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dsc_catalog_top_classifications" "test" {
  type_id = "%[1]s"
}
`, acceptance.HW_DSC_TYPE_ID)
}
