package rgc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourcePreLaunchCheck_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rgc_pre_launch_check.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourcePreLaunchCheck_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("pre_launch_check", "true"),
					resource.TestCheckResourceAttrSet(dataSource, "pre_launch_check"),
				),
			},
		},
	})
}

const testDataSourcePreLaunchCheck_basic = `
data "huaweicloud_rgc_pre_launch_check" "test" {}

output "pre_launch_check" {
  value = data.huaweicloud_rgc_pre_launch_check.test.pre_launch_check
}
`
