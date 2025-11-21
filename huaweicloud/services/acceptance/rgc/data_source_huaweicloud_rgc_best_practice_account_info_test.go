package rgc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRgcBestPracticeAccountInfo_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rgc_best_practice_account_info.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRgcBestPracticeAccountInfo_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "account_type"),
					resource.TestCheckResourceAttrSet(dataSource, "effective_start_time"),
					resource.TestCheckResourceAttrSet(dataSource, "effective_expiration_time"),
					resource.TestCheckResourceAttrSet(dataSource, "current_time"),
				),
			},
		},
	})
}

var testAccDataSourceRgcBestPracticeAccountInfo_basic = `
data "huaweicloud_rgc_best_practice_account_info" "test" {}
`
