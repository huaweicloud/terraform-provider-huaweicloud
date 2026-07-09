package dsc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCatalogStatistics_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_dsc_catalog_statistics.test"
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
				Config: testDataSourceCatalogStatistics_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "bucket.#"),
					resource.TestCheckResourceAttrSet(dataSource, "column.#"),
					resource.TestCheckResourceAttrSet(dataSource, "database.#"),
					resource.TestCheckResourceAttrSet(dataSource, "file.#"),
					resource.TestCheckResourceAttrSet(dataSource, "table.#"),
				),
			},
		},
	})
}

func testDataSourceCatalogStatistics_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dsc_catalog_statistics" "test" {
  type_id = "%[1]s"
}
`, acceptance.HW_DSC_TYPE_ID)
}
