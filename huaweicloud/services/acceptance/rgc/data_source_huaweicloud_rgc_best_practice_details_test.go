package rgc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceBestPracticeDetails_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rgc_best_practice_details.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBestPracticeDetails_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "details.#"),
					resource.TestCheckResourceAttrSet(dataSource, "details.0.check_item"),
					resource.TestCheckResourceAttrSet(dataSource, "details.0.check_item_name"),
					resource.TestCheckResourceAttrSet(dataSource, "details.0.risk_description"),
					resource.TestCheckResourceAttrSet(dataSource, "details.0.result"),
					resource.TestCheckResourceAttrSet(dataSource, "details.0.scene"),
					resource.TestCheckResourceAttrSet(dataSource, "details.0.risk_level"),
				),
			},
		},
	})
}

const testAccDataSourceBestPracticeDetails_basic = `
data "huaweicloud_rgc_best_practice_details" "test" {
}
`
