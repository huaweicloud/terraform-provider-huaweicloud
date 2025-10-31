package rgc

import (
	"fmt"

	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceOperation_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rgc_operation.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRGCOperation(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOperationByOuId_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "operation_id"),
					resource.TestCheckResourceAttrSet(dataSource, "percentage_complete"),
					resource.TestCheckResourceAttrSet(dataSource, "status"),
					resource.TestCheckResourceAttrSet(dataSource, "percentage_details.#"),
					resource.TestCheckResourceAttrSet(dataSource, "percentage_details.0.percentage_name"),
					resource.TestCheckResourceAttrSet(dataSource, "percentage_details.0.percentage_status"),
				),
			},
		},
	})
}

func testAccDataSourceOperationByOuId_basic() string {
	return fmt.Sprintf(`
data  "huaweicloud_rgc_operation" "test" {
  organizational_unit_id = "%[1]s"
}
`, acceptance.HW_RGC_ORGANIZATIONAL_UNIT_ID)
}
