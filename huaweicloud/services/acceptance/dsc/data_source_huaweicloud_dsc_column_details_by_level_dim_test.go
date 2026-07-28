package dsc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDscColumnDetailsByLevelDim_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_dsc_column_details_by_level_dim.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDscLabelId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDscColumnDetailsByLevelDim_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "results.#"),
					resource.TestCheckResourceAttrSet(dataSource, "results.0.level_name"),
				),
			},
		},
	})
}

func testAccDataSourceDscColumnDetailsByLevelDim_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dsc_column_details_by_level_dim" "test" {
  label_id = "%[1]s"
}
`, acceptance.HW_DSC_LABEL_ID)
}
