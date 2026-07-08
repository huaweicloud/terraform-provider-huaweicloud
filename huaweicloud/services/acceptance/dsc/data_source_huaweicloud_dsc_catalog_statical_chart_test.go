package dsc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDscCatalogStaticalChart_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_dsc_catalog_statical_chart.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDscTypeId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDscCatalogStaticalChart_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "detection_rules.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "sensitive_col_infos.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "total_column_number"),
				),
			},
		},
	})
}

func testAccDataSourceDscCatalogStaticalChart_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dsc_catalog_statical_chart" "test" {
  type_id = "%[1]s"
}
`, acceptance.HW_DSC_TYPE_ID)
}
