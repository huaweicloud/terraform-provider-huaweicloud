package dsc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDscColumnDetailsByDatabase_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_dsc_column_details_by_database.test"
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
				Config: testAccDataSourceDscColumnDetailsByDatabase_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "results.#"),
				),
			},
		},
	})
}

func testAccDataSourceDscColumnDetailsByDatabase_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dsc_column_details_by_database" "test" {
  type_id = "%[1]s"
}
`, acceptance.HW_DSC_TYPE_ID)
}
