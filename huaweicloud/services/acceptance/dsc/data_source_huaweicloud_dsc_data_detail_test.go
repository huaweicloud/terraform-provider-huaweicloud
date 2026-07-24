package dsc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDscDataDetail_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_dsc_data_detail.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckDscInstance(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDscDataDetail_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "asset_type"),
					resource.TestCheckResourceAttrSet(dataSource, "ins_type"),
					resource.TestCheckResourceAttrSet(dataSource, "vpc_id"),
					resource.TestCheckResourceAttrSet(dataSource, "create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "scan_detail.#"),
					resource.TestCheckResourceAttrSet(dataSource, "security_strategy.#"),
					resource.TestCheckResourceAttrSet(dataSource, "threat_analysis.#"),
				),
			},
		},
	})
}

func testAccDataSourceDscDataDetail_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dsc_data_detail" "test" {
  ins_id = "%[1]s"
  type   = "DATABASE"
}
`, acceptance.HW_DSC_INSTANCE_ID)
}
